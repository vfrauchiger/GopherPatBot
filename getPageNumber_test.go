package main

import (
	"fmt"
	"testing"
)

func TestGetNumberOfPages(t *testing.T) {
	publication := "US.9519121.B2"

	numberOfPages := getNumberOfPages(publication)

	fmt.Println(numberOfPages)

}
