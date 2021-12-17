package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func loadExcTable(filename string) []string {
	// function takes a filename (.xlsx) and return the first col of the first sheet as a slice of strings
	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// get name of first worksheet
	table := f.GetSheetName(0)

	//read all cols
	cols, err := f.GetCols(table)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(cols[0])

	//return the first column without the header cell
	return cols[0][1:]
}
