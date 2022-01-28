package main

import (
	"fmt"
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

func main() {
	fmt.Println()
	targetPath := "."
	if len(os.Args) > 1 {
		targetPath = os.Args[1]
	}
	fileMap := make(map[string]*TopLevel)
	// Start iterating over files...:
	go iterate(targetPath)

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
		if !ok || t1.Sub(t0) > 30*time.Millisecond {
			cursor.Up(prevHeight)
			prevHeight = showTable(fileMap)
			t0 = t1
		}
		if !ok {
			break
		}
	}
	fmt.Println()
}

func showTable(fileMap map[string]*TopLevel) int {
	// sort by total number of bytes in top level directory:
	var topLevels []*TopLevel
	for _, tl := range fileMap {
		topLevels = append(topLevels, tl)
	}
	sortByBytes(topLevels)

	table := makeTable(topLevels)

	max := 20
	widths := []int{0, 0, 0}
	for i, row := range table {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
		if i == max {
			break
		}
	}
	ret := 0
	for _, row := range table {
		for i, cell := range row {
			fmt.Printf("%*s", widths[i], cell)
		}
		fmt.Println()
		ret++
		if ret > max {
			break
		}
	}
	return ret
}

func makeTable(topLevels []*TopLevel) [][]string {
	table := make([][]string, len(topLevels))
	for i, tl := range topLevels {
		table[i] = make([]string, 3)
		table[i][0] = fmt.Sprintf("%s B ", commafiedInt(int(tl.Bytes)))
		table[i][1] = fmt.Sprintf("%s: ", tl.Name)
		table[i][2] = fmt.Sprintf("%d files     ", tl.NumFiles)
	}
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

func iterate(path string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip read errors for now; typically, these are permissions issues;
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
	close(fileChan)
}

func commafiedInt(i int) string {
	s := fmt.Sprintf("%d", i)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}
