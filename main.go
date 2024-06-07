package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
)

// Fetches the HTML content from a URL and returns it as a string
func fetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

type Content struct {
	Address string
	Payload string
	Value   string
}

// Parses the HTML content and extracts specific elements
func parseHTML(htmlContent string) *Content {
	content := new(Content)
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal(err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {

		if n.Type == html.TextNode && strings.Contains(n.Data, "Encrypted") {
			fmt.Println("Element content:", n.Data)
			content.Payload = n.Data
		}

		if n.Type == html.ElementNode && n.Data == "div" {
			for _, attribute := range n.Attr {
				if attribute.Key == "class" && strings.Contains(attribute.Val, "payload") {
					if n.FirstChild.LastChild != nil {
						content.Payload = n.FirstChild.LastChild.Data
					}
				}
			}
			if n.FirstChild != nil && strings.Contains(n.FirstChild.Data, "NOT") {
				content.Value = n.FirstChild.Data
			}
			if n.FirstChild != nil && strings.Contains(n.FirstChild.Data, "TON") {
				content.Value = n.FirstChild.Data
			}
		}

		if n.Type == html.ElementNode && n.Data == "a" {
			if strings.Contains(n.FirstChild.Data, "NOT") {
				content.Value = n.FirstChild.Data
			}
			for _, attribute := range n.Attr {
				if attribute.Key == "href" && attribute.Val == "/EQD5X3jciHiG4dA8fI3Y6oiXMkibk3RCJ0U2gFmeTsee2sgC" {
					content.Address = "UQD5X3jciHiG4dA8fI3Y6oiXMkibk3RCJ0U2gFmeTsee2pXH"
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return content
}

func main() {
	// Open the Excel file
	f, err := excelize.OpenFile("NOT_TON.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	finishd := false
	sellCont := 1
	for !finishd {
		fmt.Println("get count :", sellCont)
		cellValue, err := f.GetCellValue("بررسی فوری", fmt.Sprintf("D%d", sellCont))
		if err != nil {
			sellCont = sellCont + 1
			continue
		}
		if cellValue == "" {
			sellCont = sellCont + 1
			finishd = true
			continue
		}

		htmlContent, err := fetchHTML(cellValue)
		if err != nil {
			sellCont = sellCont + 1
			continue
		}

		content := parseHTML(htmlContent)
		err = f.SetCellValue("بررسی فوری", fmt.Sprintf("G%d", sellCont), content.Address)
		if err != nil {
			sellCont = sellCont + 1
			continue
		}
		err = f.SetCellValue("بررسی فوری", fmt.Sprintf("H%d", sellCont), content.Payload)
		if err != nil {
			sellCont = sellCont + 1
			continue
		}
		err = f.SetCellValue("بررسی فوری", fmt.Sprintf("I%d", sellCont), content.Value)
		if err != nil {
			sellCont = sellCont + 1
			continue
		}

		sellCont = sellCont + 1
	}

	err = f.Save()
	if err != nil {
		fmt.Println("error to SAVE :", err.Error())
	}

}
