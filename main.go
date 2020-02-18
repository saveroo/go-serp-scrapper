package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/gocolly/colly"
)

var google = "https://www.google.com/search?q="

// Void
func withSoup(query string) {
	resp, err := soup.Get(google + query)
	if err != nil {
		os.Exit(1)
	}
	var re = regexp.MustCompile(`(?m)/url\?q=(.*)`)
	doc := soup.HTMLParse(resp)
	links := doc.Find("div", "id", "main").FindAll("a")
	for _, link := range links {
		if re.MatchString(link.Attrs()["href"]) {
			fmt.Println("| Link:", re.FindStringSubmatch(link.Attrs()["href"])[1])
		}
	}
}

// Void
func withColly(query string) {
	c := colly.NewCollector()
	re := regexp.MustCompile(`(?m)/url\?q=(.*)`)
	c.OnHTML("#main", func(e *colly.HTMLElement) {
		dom := e
		dom.ForEach("a[href]", func(i int, e *colly.HTMLElement) {
			links, _ := e.DOM.Attr("href")
			if re.MatchString(links) {
				if e.DOM.ChildrenFiltered("div").Prev().Is("div") {
					title := e.DOM.ChildrenFiltered("div").Prev().Text()
					fmt.Println(i, title, "||", re.FindStringSubmatch(links)[1])
				}
			}
		})
	})
	c.Visit(google + query)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("(Unfinished Experimental)")
		fmt.Println("(type 1/ctrl+c to exit)")
		fmt.Println("=============")
		fmt.Print("Search Query: ")
		cmdStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		cmdStr = strings.TrimSuffix(cmdStr, "\n")
		if cmdStr == "1" {
			os.Exit(1)
		}
		withColly(cmdStr)
	}
}
