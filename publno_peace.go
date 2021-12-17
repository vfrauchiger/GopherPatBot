package main

//little module to modify publication numbers

import (
	"fmt"
	"regexp"
	"strings"
)

func numberIngestion(publno string) []string {
	publno = strings.ReplaceAll(publno, "-", "")
	publno = strings.ReplaceAll(publno, "/", "")
	publno = strings.ReplaceAll(publno, " ", "")
	re := regexp.MustCompile(`(?P<cc>\w{2})(?P<nu>\d+)(?P<kd>[A-Za-z][0-9A-Za-z]?)`)
	match := re.FindStringSubmatch(publno)
	publnoList := []string{match[1], match[2], match[3]}
	fmt.Println(publnoList)

	if publnoList[0] == "US" {
		publnoList = treatUS(publnoList)
	} else if publnoList[0] == "EP" {
		publnoList = treatEP(publnoList)
	} else if publnoList[0] == "WO" {
		publnoList = treatWO(publnoList)
	}
	return publnoList
}

func treatUS(publnoList []string) []string {
	kdCodeDict := map[string]string{
		"AA": "A1",
		"AB": "A2",
		"BA": "B1",
		"BB": "B2",
		"E":  "E1",
	}
	keys := make([]string, 0, len(kdCodeDict))
	for key := range kdCodeDict {
		keys = append(keys, key)
	}
	// modify kind code
	for _, key := range keys {
		if publnoList[2] == key {
			publnoList[2] = kdCodeDict[key]
		}
	}

	// modify number
	if len(publnoList[1]) == 11 {
		publnoList[1] = publnoList[1][:4] + publnoList[1][5:]
	}

	return publnoList

}

/*
Von 1978 bis zum 30. Juni 2002 lautete es WOJJNNNNN (Ländercode, 2 Ziffern für das Jahr und 5 Ziffern für die laufende Nummer).
Vom 1. Juli 2002 bis zum 31. Dezember 2003 lautete es WOJJNNNNNN (Ländercode, 2 Ziffern für das Jahr und 6 Ziffern für die laufende Nummer).
Seit dem 1. Januar 2004 lautet es WOJJJJNNNNNN (Ländercode, 4 Ziffern für das Jahr und 6 Ziffern für die laufende Nummer).
*/

func treatWO(publnoList []string) []string {
	//modify number

	if len(publnoList[1]) == 7 {
		return publnoList
	} else if len(publnoList[1]) == 8 {
		if publnoList[1][:2] == "02" || publnoList[1][:2] == "03" {
			return publnoList
		} else {
			publnoList[1] = "20" + publnoList[1]
			return publnoList
		}
	} else if len(publnoList[1]) == 10 {
		if publnoList[1][2:4] == "02" || publnoList[1][2:4] == "03" {
			publnoList[1] = publnoList[1][2:]
			return publnoList
		} else {
			return publnoList
		}
	}

	return publnoList
}

func treatEP(publnoList []string) []string {
	return publnoList
}
