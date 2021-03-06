package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// In orde to get the code working you will have to obtain application keys from the European Patent Office
var consKey string = ""    // add OPS Key here
var consSecKey string = "" // add OPS Secret here

var token string

type AuthResponse struct {
	RefreshTokenExpiresIn string   `json:"refresh_token_expires_in"`
	APIProductList        string   `json:"api_product_list"`
	APIProductListJSON    []string `json:"api_product_list_json"`
	OrganizationName      string   `json:"organization_name"`
	DeveloperEmail        string   `json:"developer.email"`
	TokenType             string   `json:"token_type"`
	IssuedAt              string   `json:"issued_at"`
	ClientID              string   `json:"client_id"`
	AccessToken           string   `json:"access_token"`
	ApplicationName       string   `json:"application_name"`
	Scope                 string   `json:"scope"`
	ExpiresIn             string   `json:"expires_in"`
	RefreshCount          string   `json:"refresh_count"`
	Status                string   `json:"status"`
}

type WorldPatentData struct {
	XMLName         xml.Name `xml:"world-patent-data"`
	Text            string   `xml:",chardata"`
	Xmlns           string   `xml:"xmlns,attr"`
	Ops             string   `xml:"ops,attr"`
	Xlink           string   `xml:"xlink,attr"`
	DocumentInquiry struct {
		Text                 string `xml:",chardata"`
		PublicationReference struct {
			Text       string `xml:",chardata"`
			DocumentID struct {
				Text           string `xml:",chardata"`
				DocumentIDType string `xml:"document-id-type,attr"`
				Country        string `xml:"country"`
				DocNumber      string `xml:"doc-number"`
				Kind           string `xml:"kind"`
			} `xml:"document-id"`
		} `xml:"publication-reference"`
		InquiryResult struct {
			Text                 string `xml:",chardata"`
			PublicationReference struct {
				Text       string `xml:",chardata"`
				DocumentID struct {
					Text           string `xml:",chardata"`
					DocumentIDType string `xml:"document-id-type,attr"`
					Country        string `xml:"country"`
					DocNumber      string `xml:"doc-number"`
					Kind           string `xml:"kind"`
				} `xml:"document-id"`
			} `xml:"publication-reference"`
			DocumentInstance []struct {
				Text                  string `xml:",chardata"`
				System                string `xml:"system,attr"`
				NumberOfPages         string `xml:"number-of-pages,attr"`
				Desc                  string `xml:"desc,attr"`
				Link                  string `xml:"link,attr"`
				DocumentFormatOptions struct {
					Text           string   `xml:",chardata"`
					DocumentFormat []string `xml:"document-format"`
				} `xml:"document-format-options"`
				DocumentSection []struct {
					Text      string `xml:",chardata"`
					Name      string `xml:"name,attr"`
					StartPage string `xml:"start-page,attr"`
				} `xml:"document-section"`
			} `xml:"document-instance"`
		} `xml:"inquiry-result"`
	} `xml:"document-inquiry"`
}

func authenticatePB() string {
	authURL := "https://ops.epo.org/3.2/auth/accesstoken"
	encKey := consKey + ":" + consSecKey
	encKey = base64.StdEncoding.EncodeToString([]byte(encKey))
	//Body of R
	reqBody := strings.NewReader("grant_type=client_credentials")

	req, err := http.NewRequest(
		"POST",
		authURL,
		reqBody,
	)
	if err != nil {
		fmt.Println("Request-Error: ", err)
	}
	// Headers
	req.Header.Add("Authorization", string(encKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Response Error: ", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read Error: ", err)
	}
	resp.Body.Close()

	var jsonData AuthResponse

	json.Unmarshal(data, &jsonData)

	fmt.Println("Status: ", resp.StatusCode)
	fmt.Println(jsonData.AccessToken, jsonData.DeveloperEmail)

	return jsonData.AccessToken
}

func getNumberOfPages(publno string) int {

	pageURL := "https://ops.epo.org/rest-services/published-data/publication/docdb/" + publno + "/images"

	token = authenticatePB()

	//reqBody := strings.NewReader(publno)

	req, err := http.NewRequest(
		"GET",
		pageURL,
		nil,
	)
	if err != nil {
		fmt.Println("Page Number Req error: ", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	fmt.Println(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Response Error Pages: ", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read Error Pages: ", err)
	}
	resp.Body.Close()

	var xmlResp WorldPatentData
	xml.Unmarshal(data, &xmlResp)

	//fmt.Println(string(data))
	if strings.Contains(string(data), "No results found") == true {
		fmt.Printf("The Publication %s was not found!\n", publno)
		return 0
	} else {
		fmt.Println(len(xmlResp.DocumentInquiry.InquiryResult.DocumentInstance))
		for index, _ := range xmlResp.DocumentInquiry.InquiryResult.DocumentInstance {
			if xmlResp.DocumentInquiry.InquiryResult.DocumentInstance[index].Desc == "FullDocument" {
				numbOfPages, err := strconv.Atoi(xmlResp.DocumentInquiry.InquiryResult.DocumentInstance[index].NumberOfPages)
				if err != nil {
					fmt.Println("Conversion Error: ", err)
				}
				return numbOfPages
			} else {
				continue
			}
		}

	}
	return 0
}

func getOnePublication(publnoSlice []string) {
	var fileNames []string

	publno := publnoSlice[0] + "." + publnoSlice[1] + "." + publnoSlice[2]
	numbOfPages := getNumberOfPages(publno)
	if numbOfPages == 0 {
		return
	}
	fmt.Printf("The Publication has %d pages\n", numbOfPages)

	for i := 1; i < (numbOfPages + 1); i++ {
		var reader io.Reader
		urlpdf := "http://ops.epo.org/rest-services/published-data/images/" + publnoSlice[0] + "/" + publnoSlice[1] + "/" + publnoSlice[2] + "/fullimage.pdf?Range=" + strconv.Itoa(i)
		//savePath := publno + "_" + strconv.Itoa(i) + ".pdf"
		client := &http.Client{}

		req, err := http.NewRequest(
			"GET",
			urlpdf,
			reader,
		)
		if err != nil {
			fmt.Println(err)
		}

		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()

		savePath := "./" + publnoSlice[0] + publnoSlice[1] + publnoSlice[2] + "_" + strconv.Itoa(i) + ".pdf"
		fileNames = append(fileNames, savePath)
		file, err := os.Create(savePath)
		if err != nil {
			fmt.Println(err)
		}
		//defer file.Close()

		pages, err := io.Copy(file, resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(pages)
		file.Close()
		fmt.Println(savePath)

	}

	err := api.MergeCreateFile(fileNames, "./"+publnoSlice[0]+publnoSlice[1]+publnoSlice[2]+".pdf", nil)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, file := range fileNames {
			os.Remove(file)
		}
	}

}
