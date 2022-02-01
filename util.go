package main

import (
	"fmt"
	"strings"
)

func commafiedInt(i int) string {
	s := fmt.Sprintf("%d", i)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

func topOfPath(targetPath, path string) string {
	var newPath string
	if strings.HasPrefix(path, targetPath) {
		newPath = path[len(targetPath):]
	} else {
		newPath = path
	}
	terms := strings.Split(newPath, "/")
	for _, term := range terms {
		if term == "." || term == ".." || term == "" {
			continue
		}
		return term
	}
	return ""
}
