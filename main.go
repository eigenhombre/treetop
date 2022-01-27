package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	targetPath := "."
	if len(os.Args) > 1 {
		targetPath = os.Args[1]
	}
	fileMap := make(map[string]*TopLevel)
	// Start goroutine to iterate over files
	go iterate(targetPath)

	// Report results
	for file := range fileChan {
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
	// sort by total number of bytes in top level directory:
	var topLevels []*TopLevel
	for _, tl := range fileMap {
		topLevels = append(topLevels, tl)
	}
	sortByBytes(topLevels)

	table, widths := makeTable(topLevels)

	for _, row := range table {
		for i, cell := range row {
			fmt.Printf("%*s", widths[i], cell)
		}
		fmt.Println()
	}
}

func makeTable(topLevels []*TopLevel) ([][]string, []int) {
	table := make([][]string, len(topLevels))
	for i, tl := range topLevels {
		table[i] = make([]string, 3)
		table[i][0] = fmt.Sprintf("%s B ", commafiedInt(int(tl.Bytes)))
		table[i][1] = fmt.Sprintf("%s:", tl.Name)
		table[i][2] = fmt.Sprintf("%d files", tl.NumFiles)
	}
	widths := []int{0, 0, 0}
	for _, row := range table {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	return table, widths
}

func sortByBytes(topLevels []*TopLevel) {
	for i := 0; i < len(topLevels); i++ {
		for j := i + 1; j < len(topLevels); j++ {
			if topLevels[i].Bytes < topLevels[j].Bytes {
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
			log.Fatalf(err.Error())
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
