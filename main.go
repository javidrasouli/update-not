package main

import (
	"context"
	"fmt"
	"github.com/tonkeeper/tonapi-go"
	"github.com/xuri/excelize/v2"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var Contract = "EQCvxJy4eG8hyHBFsZ7eePxrRsUQSFE_jpptRAYBmcG_DOGS"

// Fetches the HTML content from a URL and returns it as a string
func fetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error is :", err.Error())
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
	Address         string
	Payload         string
	Value           string
	Coin            string
	ContractAddress bool
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
			if n.FirstChild != nil && strings.Contains(n.FirstChild.Data, "DOGS") {
				content.Value = n.FirstChild.Data
				content.Coin = "NOT"
			}
			if n.FirstChild != nil && strings.Contains(n.FirstChild.Data, "TON") {
				content.Value = n.FirstChild.Data
				content.Coin = "TON"
			}
			if n.FirstChild != nil && strings.Contains(n.FirstChild.Data, "USD₮") {
				content.Value = n.FirstChild.Data
				content.Coin = "USDT"
			}
		}

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, i := range n.Attr {
				if i.Key == "href" && strings.Contains(i.Val, Contract) {
					content.Value = n.FirstChild.Data
					content.ContractAddress = true
					content.Coin = "DOGS"
				}
			}
			if strings.Contains(n.FirstChild.Data, "NOT") {
				content.Value = n.FirstChild.Data
				content.Coin = "NOT"
			}
			if strings.Contains(n.FirstChild.Data, "TON") {
				content.Value = n.FirstChild.Data
				content.Coin = "TON"
			}
			if strings.Contains(n.FirstChild.Data, "USD₮") {
				content.Value = n.FirstChild.Data
				content.Coin = "USDT"
			}
			for _, attribute := range n.Attr {
				if attribute.Key == "href" && attribute.Val == "/EQD5X3jciHiG4dA8fI3Y6oiXMkibk3RCJ0U2gFmeTsee2sgC" {
					content.Address = "UQD5X3jciHiG4dA8fI3Y6oiXMkibk3RCJ0U2gFmeTsee2pXH"
				}
				if attribute.Key == "href" && attribute.Val == "/UQD5X3jciHiG4dA8fI3Y6oiXMkibk3RCJ0U2gFmeTsee2pXH" {
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

func extractNumbers(input string) string {
	// Define the regular expression to match numbers
	re := regexp.MustCompile(`[0-9]+`)
	// Find all occurrences of numbers in the string
	numbers := re.FindAllString(input, -1)

	// Join all the numbers to form a single string
	result := ""
	for _, num := range numbers {
		result += num
	}

	return result
}

func getMineWalletAddress(address string, minAccountAddress string) string {
	if address == minAccountAddress {
		return "UQD5X3jciHiG4dA8fI3Y6oiXMkibk3RCJ0U2gFmeTsee2pXH"
	}

	return ""
}

func main() {
	client, err := tonapi.New()
	if err != nil {
		log.Fatal(err)
	}

	jet := tonapi.GetAccountParams{
		AccountID: "EQCvxJy4eG8hyHBFsZ7eePxrRsUQSFE_jpptRAYBmcG_DOGS",
	}

	jetton, err := client.GetAccount(context.Background(), jet)
	if err != nil {
		panic(err)
	}

	accParam := tonapi.GetAccountParams{
		AccountID: "UQD5X3jciHiG4dA8fI3Y6oiXMkibk3RCJ0U2gFmeTsee2pXH",
	}

	account, err := client.GetAccount(context.Background(), accParam)
	if err != nil {
		panic(err)
	}

	mineAddress := account.GetAddress()
	dogsAddr := jetton.GetAddress()

	// Open the Excel file
	f, err := excelize.OpenFile("NOT_TON.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	ownWallet := "UQD5X3jciHiG4dA8fI3Y6oiXMkibk3RCJ0U2gFmeTsee2pXH"

	finishd := false
	sellCont := 1
	for !finishd {
		fmt.Println("get count :", sellCont)

		cellValue, err := f.GetCellValue("Sheet1", fmt.Sprintf("D%d", sellCont))
		if err != nil {
			sellCont = sellCont + 1
			continue
		}
		Done, _ := f.GetCellValue("Sheet1", fmt.Sprintf("M%d", sellCont))
		if Done == "DONE" {
			sellCont = sellCont + 1
			continue
		}

		time.Sleep(2 * time.Second)

		if cellValue == "" {
			sellCont = sellCont + 1
			finishd = true
			continue
		}

		//url := fmt.Sprintf("https://tonviewer.com/transaction/%s", cellValue)
		//
		//htmlContent, err := fetchHTML(url)
		//if err != nil {
		//	sellCont = sellCont + 1
		//	continue
		//}
		//content := parseHTML(htmlContent)

		param := tonapi.GetJettonsEventsParams{
			EventID: cellValue,
		}

		tr, err := client.GetJettonsEvents(context.Background(), param)
		if err != nil {
			fmt.Println("error to get tarnsaction : ", err.Error())
			sellCont = sellCont + 1
			finishd = true
			continue
		}

		content := new(Content)
		actions := tr.GetActions()

		if len(actions) > 0 {
			privew := actions[0].GetSimplePreview()
			transfer, ok := actions[0].GetJettonTransfer().Get()
			if ok {
				content.Payload = transfer.Comment.Value
				content.Value = privew.Value.Value
				content.Coin = transfer.Jetton.Name
				if content.Coin == "Dogs" {
					content.ContractAddress = dogsAddr == transfer.Jetton.Address
				}
				content.Address = getMineWalletAddress(transfer.Recipient.Value.GetAddress(), mineAddress)
			}
		}

		err = f.SetCellValue("Sheet1", fmt.Sprintf("I%d", sellCont), content.Address)
		if err != nil {
			sellCont = sellCont + 1
			continue
		}
		err = f.SetCellValue("Sheet1", fmt.Sprintf("H%d", sellCont), content.Payload)
		if err != nil {
			sellCont = sellCont + 1
			continue
		}
		err = f.SetCellValue("Sheet1", fmt.Sprintf("J%d", sellCont), content.Value)
		if err != nil {
			sellCont = sellCont + 1
			continue
		}

		checkCount := content.Address == ownWallet

		err = f.SetCellValue("Sheet1", fmt.Sprintf("L%d", sellCont), checkCount)
		if err != nil {
			sellCont = sellCont + 1
			continue
		}

		if !checkCount {
			if content.Address != ownWallet {
				err = f.SetCellValue("Sheet1", fmt.Sprintf("M%d", sellCont), "address not match")
				if err != nil {
					sellCont = sellCont + 1
					continue
				}
			}
		}

		if content.Coin == "Dogs" {
			err = f.SetCellValue("Sheet1", fmt.Sprintf("K%d", sellCont), content.ContractAddress)
			if err != nil {
				sellCont = sellCont + 1
				continue
			}
		}
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("M%d", sellCont), "DONE")

		sellCont = sellCont + 1

	}

	err = f.Save()
	if err != nil {
		fmt.Println("error to SAVE :", err.Error())
	}

}
