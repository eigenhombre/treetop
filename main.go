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
	maxTopOfPath := 0
	for file := range fileChan {
		if file.Path == "." {
			continue
		}
		topOfPath := topOfPath(file.Path)
		if len(topOfPath) > maxTopOfPath {
			maxTopOfPath = len(topOfPath)
		}
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
	// sort by bytes
	var topLevels []*TopLevel
	for _, tl := range fileMap {
		topLevels = append(topLevels, tl)
	}
	sortByBytes(topLevels)

	// print results
	for _, tl := range topLevels {
		fmtString := fmt.Sprintf("%%14d B %%%ds: %%6d files\n", maxTopOfPath)
		fmt.Printf(fmtString, tl.Bytes, tl.Name, tl.NumFiles)
	}
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

func topOfPath(path string) string {
	terms := strings.Split(path, "/")
	for _, term := range terms {
		if term == "." || term == ".." {
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
