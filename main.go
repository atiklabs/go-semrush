package main

import (
	"fmt"
	"flag"
)

type result struct {
	url string
	score string
}

func main() {
	file := flag.String("f", "", "CSV file with the url in the first column")
	flag.Parse()
	if (*file == "") {
		flag.Usage()
	} else {
		//API be963c9550758935863b1583b41a6ef6
		urlList := parseFile(file)
		processUrlList(urlList)
	}
}

func parseFile(file *string) []string {
	urlList := []string {"a", "b"}
	return urlList
}

func processUrlList(urlList []string) {
	for i := 0; i < len(urlList); i++ {
		fmt.Printf("URL: %s\n", urlList[i])
	}
}
