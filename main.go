package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/anaskhan96/soup"
)

type Subtitle struct {
	Id        string   `json:"id"`
	Desc      string   `json:"desc"`
	Languages []string `json:"languages"`
}

func main() {
	soup.Header("User-Agent", "curl/7.64.1")
	resp, err := soup.Get("https://www.a4k.net/search?term=%E7%A1%85%E8%B0%B7")
	if err != nil {
		log.Printf("there's an error: %v", err)
		return
	}
	doc := soup.HTMLParse(resp)
	//fmt.Println(doc.FullText())
	items := doc.FindStrict("ul", "class", "ui relaxed divided list").FindAll("li")

	var funcGetLanguages func(nodes []soup.Root) []string = func(nodes []soup.Root) []string {
		var result []string
		for _, item := range nodes {
			language := item.Attrs()["data-content"]
			result = append(result, language)
		}
		return result
	}

	for _, item := range items[:5] {
		i := item.FindStrict("div", "class", "content").Find("h3").Find("a")

		subtitle := &Subtitle{
			Id:        i.Attrs()["href"][strings.LastIndex(i.Attrs()["href"], "/")+1:],
			Desc:      i.Text(),
			Languages: funcGetLanguages(item.FindStrict("div", "class", "language").Find("span", "class", "h4").FindAll("i")),
		}

		jsonText, _ := json.Marshal(subtitle)
		fmt.Println(string(jsonText))
	}
}
