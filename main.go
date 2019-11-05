package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func downloadPic(page string, url string) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	//open file for editing
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	filename := dir + fmt.Sprintf("/files/%v.png", page)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Success! Page %s downloaded!\n", page)
}

func getUrl(page string, url string) {
	loc := url + page

	response, e := http.Get(loc)
	if e != nil {
		log.Fatal(e)
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	z := html.NewTokenizer(response.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.SelfClosingTagToken:
			t := z.Token()

			isAnchor := t.Data == "img"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "src" {
						value := a.Val
						runes := []rune(value)
						substring := string(runes[16:22])
						if substring == "comics" {
							newString := "http://www." + string(runes[2:])
							downloadPic(page, newString)
							return
						}
					}
				}
			}
		}
	}
}

func main() {
	website := "https://xkcd.com/"

	start := time.Now()
	for i := 1; i <= 1960; i++ {
		go getUrl(strconv.Itoa(i), website)
		time.Sleep(time.Millisecond * 20)
	}
	end := time.Since(start)
	fmt.Println(end)

	var input string
	if _, err := fmt.Scanln(&input); err != nil {
		log.Fatal(err)
	}
}
