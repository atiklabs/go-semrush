package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"encoding/csv"
	"strings"
	"io"
)

type result struct {
	url string
	score string
}

func main() {
	file := flag.String("f", "", "CSV file with the url in the first column")
	apiKey := flag.String("api", "", "API Key given by SEMRush")
	flag.Parse()
	if *file == "" || *apiKey == "" {
		flag.Usage()
	} else {
		urlList := parseFile(file)
		processUrlList(urlList, apiKey)
	}
}

func parseFile(file *string) []string {
	dat, err := ioutil.ReadFile(*file)
	check(err)

	r := csv.NewReader(strings.NewReader(string(dat)))
	var urlList []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		check(err)
		urlList = append(urlList, record[0])
	}
	return urlList
}

func processUrlList(urlList []string, apiKey *string) {
	for i := 0; i < len(urlList); i++ {
		fmt.Printf("URL: %s\n", urlList[i])
		LetsPrint()
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
