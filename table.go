package main

import "fmt"

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
