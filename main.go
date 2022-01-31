package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atomicgo/cursor"
)

type FileInfo struct {
	Name string
	Size int64
	Path string
}

type TopLevel struct {
	Name     string
	NumFiles int64
	Bytes    int64
}

var fileChan = make(chan FileInfo, 10000)

func dirsInDir(dir string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("can't read dir: " + dir)
	}
	ret := []fs.FileInfo{}
	for _, file := range files {
		if file.IsDir() {
			ret = append(ret, file)
		}
	}
	return ret, nil
}

func main() {
	fmt.Println()
	targetPath := "."
	if len(os.Args) > 1 {
		targetPath = os.Args[1]
	}
	fileMap := make(map[string]*TopLevel)

	// Files in top level of directory:
	files, err := dirsInDir(targetPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	type done struct{}
	dirsOpenChan := make(chan done, len(files))

	for _, file := range files {
		go func(file string) {
			collectDirStats(targetPath + "/" + file)
			dirsOpenChan <- done{}
		}(file.Name())
	}

	// goroutine that waits for all dir listings to be finished:
	go func() {
		for range files {
			<-dirsOpenChan
		}
		close(fileChan)
	}()

	// ... consume and report results:
	prevHeight := 0
	t0 := time.Now()
	for {
		file, ok := <-fileChan
		if ok {
			topOfPath := topOfPath(targetPath, file.Path)
			tl := fileMap[topOfPath]
			if tl == nil {
				tl = &TopLevel{
					Name:     topOfPath,
					Bytes:    file.Size,
					NumFiles: 1,
				}
			} else {
				tl.Bytes += file.Size
				tl.NumFiles++
			}
			fileMap[topOfPath] = tl
		}
		t1 := time.Now()
		if !ok || t1.Sub(t0) > 100*time.Millisecond {
			cursor.Up(prevHeight)
			prevHeight = showTable(makeTable(sortedTopLevels(fileMap), 20))
			t0 = t1
		}
		if !ok {
			break
		}
	}
	fmt.Println()
}

// sortedTopLevels returns a slice of topLevels sorted by name.
func sortedTopLevels(fileMap map[string]*TopLevel) []*TopLevel {
	// sort by total number of bytes in top level directory:
	var topLevels []*TopLevel
	for _, tl := range fileMap {
		topLevels = append(topLevels, tl)
	}
	sortByBytes(topLevels)
	return topLevels
}

func showTable(table [][]string) int {
	widths := []int{0, 0, 0}
	for _, row := range table {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	ret := 0
	for _, row := range table {
		for i, cell := range row {
			fmt.Printf("%*s", widths[i], cell)
		}
		fmt.Println()
		ret++
	}
	return ret
}

func makeTable(topLevels []*TopLevel, maxTopLevels int) [][]string {
	table := [][]string{}
	i, bytes, files := 0, 0, 0
	for _, tl := range topLevels {
		if i < maxTopLevels {
			table = append(table, []string{
				fmt.Sprintf("%s B ", commafiedInt(int(tl.Bytes))),
				fmt.Sprintf("%s: ", tl.Name),
				fmt.Sprintf("%d files%30s", tl.NumFiles, ""),
			})
		}
		i++
		bytes += int(tl.Bytes)
		files += int(tl.NumFiles)
	}
	table = append(table, []string{"", "", ""})
	table = append(table, []string{
		fmt.Sprintf("%s B ", commafiedInt(bytes)),
		"TOTAL  ",
		fmt.Sprintf("%s files%30s", commafiedInt(files), ""),
	})
	return table
}

func sortByBytes(topLevels []*TopLevel) {
	for i := 0; i < len(topLevels); i++ {
		for j := i + 1; j < len(topLevels); j++ {
			if topLevels[i].Bytes < topLevels[j].Bytes ||
				topLevels[i].Bytes == topLevels[j].Bytes && topLevels[i].Name > topLevels[j].Name {
				topLevels[i], topLevels[j] = topLevels[j], topLevels[i]
			}
		}
	}
}

func topOfPath(targetPath, path string) string {
	newPath := strings.Replace(path, targetPath, "", 1)
	terms := strings.Split(newPath, "/")
	for _, term := range terms {
		if term == "." || term == ".." || term == "" {
			continue
		}
		return term
	}
	return ""
}

func collectDirStats(path string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip read errors (generally, permissions) for now; typically, these are permissions issues;
			// handle them later?  Present statistics?
			// log.Fatalf(err.Error())
			return nil
		}
		fileChan <- FileInfo{
			Name: info.Name(),
			Size: info.Size(),
			Path: path,
		}
		return nil
	})
	// close(fileChan)
}

func commafiedInt(i int) string {
	s := fmt.Sprintf("%d", i)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}
