// Copyright 2016 Go SEMRush Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"encoding/csv"
	"strings"
	"io"
	"net/url"
	"strconv"
	"os"
)

const (
	GOROUTINES = 50
)

type csvEntry struct {
	url string
	csv []string
}

var scoreMap = map[string]int{}

func main() {
	file := flag.String("f", "", "CSV file with the url in the first column")
	header := flag.Bool("header", false, "First row of file is header")
	apiKey := flag.String("api", "", "API Key given by SEMRush")
	lang := flag.String("lang", "us", "SEMRush language database")
	domain := flag.Bool("domain", false, "The first column of the CSV file is the domain name instead of the url ")
	async := flag.Bool("async", false, "50x faster, result returned is shuffled, spends more API points because connects once for each URL instead for each domain")
	flag.Parse()
	if *file == "" || *apiKey == "" {
		flag.Usage()
	} else {
		csvEntryList := parseFile(file, *header)
		if *async {
			processUrlListAsync(csvEntryList, apiKey, lang, *domain)
		} else {
			processUrlList(csvEntryList, apiKey, lang, *domain)
		}
	}
}

// This parses a CSV file and returns an array with every url
func parseFile(file *string, header bool) []csvEntry {
	dat, err := ioutil.ReadFile(*file)
	CheckError(err)
	r := csv.NewReader(strings.NewReader(string(dat)))
	var csvEntryList []csvEntry
	for {
		record, err := r.Read()
		if header {
			fmt.Println(strings.Join(record, ",") + ",Traffic By Semrush")
			header = false
			continue
		}
		if err == io.EOF {
			break
		}
		CheckError(err)
		csvEntryList = append(csvEntryList, csvEntry{record[0], record})
	}
	return csvEntryList
}

// This will use SEMRush API for every url in urlList to get a domain score
func processUrlList(csvEntryList []csvEntry, apiKey *string, lang *string, isDomain bool) {
	w := csv.NewWriter(os.Stdout)
	for i := 0; i < len(csvEntryList); i++ {
		domain := csvEntryList[i].url
		if !isDomain {
			domain = getDomainFromUrl(csvEntryList[i].url)
		}
		score, ok := scoreMap[domain];
		if !ok {
			score = GetDomainScore(domain, apiKey, lang)
			scoreMap[domain] = score
		}
		csvEntryList[i].csv = append(csvEntryList[i].csv, strconv.Itoa(score))
		err := w.Write(csvEntryList[i].csv);
		CheckError(err)
	}
	w.Flush()
	CheckError(w.Error())
}

// Same as processUrlList by threaded with goroutines
func processUrlListAsync(csvEntryList []csvEntry, apiKey *string, lang *string, isDomain bool) {
	for i := 0; i < len(csvEntryList); i += GOROUTINES {
		done := make(chan bool, GOROUTINES)
		for j := i; j < len(csvEntryList) && j < i + GOROUTINES; j++ {
			domain := csvEntryList[j].url
			if !isDomain {
				domain = getDomainFromUrl(csvEntryList[j].url)
			}
			go asyncGetDomainScore(domain, csvEntryList[j].csv, apiKey, lang, done)
		}
		for j := i; j < len(csvEntryList) && j < i + GOROUTINES; j++ {
			<-done
		}
	}
}

// This will be executed as a goroutine
func asyncGetDomainScore(domain string, csvArray []string, apiKey *string, lang *string, done chan bool) {
	score := GetDomainScore(domain, apiKey, lang)
	csvArray = append(csvArray, strconv.Itoa(score))
	w := csv.NewWriter(os.Stdout)
	err := w.Write(csvArray);
	CheckError(err)
	w.Flush()
	CheckError(w.Error())
	done <- true
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
