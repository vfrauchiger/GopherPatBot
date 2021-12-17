package main

import (
	"fmt"
	"testing"
)

func TestGetNumberOfPages(t *testing.T) {
	publication := "US.9492605.B2"

	numberOfPages := getNumberOfPages(publication)

	fmt.Println(numberOfPages)

}
