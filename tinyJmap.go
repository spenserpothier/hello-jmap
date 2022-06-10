package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func makeJmapCall(apiUrl string, jmapCall string) (jsonString string, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBufferString(jmapCall))

	req.SetBasicAuth(user, pass)

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("error: %v", err)
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return string(body), nil
}

func getAccountStuff() (accountId string, apiUrl string, err error) {

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, "https://jmap.fastmail.com/.well-known/jmap", nil)

	req.SetBasicAuth(user, pass)
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("error: %v", err)
		return "", "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	sessionMap := make(map[string]interface{})

	err = json.Unmarshal(body, &sessionMap)
	if err != nil {
		fmt.Printf("error unmarshling: %v", err)
		return "", "", err
	}
	primaryAccountsMap := sessionMap["primaryAccounts"].(map[string]interface{})
	accountId = fmt.Sprint(primaryAccountsMap["urn:ietf:params:jmap:mail"])
	apiUrl = fmt.Sprint(sessionMap["apiUrl"])
	return accountId, apiUrl, nil
}
