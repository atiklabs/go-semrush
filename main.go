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

var scoreMap = map[string]int{}

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

// This parses a CSV file and returns an array with every url
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

// This will use SEMRush API for every url in urlList to get a domain score
func processUrlList(urlList []string, apiKey *string, lang *string) {
	for i := 0; i < len(urlList); i++ {
		domain := getDomainFromUrl(urlList[i])
		score, ok := scoreMap[domain];
		if !ok {
			score = GetDomainScore(domain, apiKey, lang)
			scoreMap[domain] = score
		}
		fmt.Printf("%s,%d,%s\n", domain, score, urlList[i])
	}
}

// Extract the domain name from the url address
func getDomainFromUrl(address string) string {
	u, err := url.Parse(address)
	CheckError(err)
	return u.Host
}

// A common way to treat errors
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
