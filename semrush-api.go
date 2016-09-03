// Copyright 2016 Go SEMRush Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"net/http"
	"io/ioutil"
	"encoding/csv"
	"strings"
	"io"
	"strconv"
)

const (
	SEMRUSH_API = "http://api.semrush.com/"
)

// This will make the API call to SemRUSH to get the score of the domain
func GetDomainScore(domain string, apiKey *string, lang *string) int {
	address := SEMRUSH_API + "?key=" + *apiKey + "&database=" + *lang + "&type=domain_ranks&export_columns=Ot" + "&domain=" + domain
	res, err := http.Get(address)
	CheckError(err)
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	CheckError(err)

	r := csv.NewReader(strings.NewReader(string(body)))
	skipFirstLine := true
	for {
		record, err := r.Read()
		if skipFirstLine {
			skipFirstLine = false
			continue
		}
		if err == io.EOF {
			break
		}
		CheckError(err)
		i, err := strconv.Atoi(record[0])
		CheckError(err)
		return i
	}
	return 0
}
