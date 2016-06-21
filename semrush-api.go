package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
)

const (
	SEMRUSH_API = "http://api.semrush.com/"
)

func GetDomainScore(domain string, apiKey *string, lang *string) int {
	address := SEMRUSH_API + "?key=" + *apiKey + "&database=" + *lang + "&type=domain_ranks&export_columns=Dn,Ot" + "&domain=" + domain
	res, err := http.Get(address)
	CheckError(err)
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	CheckError(err)
	fmt.Println(address)
	fmt.Printf("%s", body)
	return 1
}