package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"encoding/csv"
	"strings"
	"io"
	"net/url"
)

type result struct {
	url string
	score string
}

func main() {
	file := flag.String("f", "", "CSV file with the url in the first column")
	apiKey := flag.String("api", "", "API Key given by SEMRush")
	lang := flag.String("lang", "us", "SEMRush language database")
	flag.Parse()
	if *file == "" || *apiKey == "" {
		flag.Usage()
	} else {
		urlList := parseFile(file)
		processUrlList(urlList, apiKey, lang)
	}
}

func parseFile(file *string) []string {
	dat, err := ioutil.ReadFile(*file)
	CheckError(err)

	r := csv.NewReader(strings.NewReader(string(dat)))
	var urlList []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		CheckError(err)
		urlList = append(urlList, record[0])
	}
	return urlList
}

func processUrlList(urlList []string, apiKey *string, lang *string) {
	for i := 0; i < len(urlList); i++ {
		domain := getDomainFromUrl(urlList[i])
		score := GetDomainScore(domain, apiKey, lang)
		fmt.Printf("%d\t%s\t%s\t%d\n", i, urlList[i], domain, score)
	}
}

func getDomainFromUrl(address string) string {
	u, err := url.Parse(address)
	CheckError(err)
	return u.Host
}

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
