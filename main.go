package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var (
	user = os.Getenv("JMAP_USERNAME")
	pass = os.Getenv("JMAP_PASSWORD")
)

func main() {

	accountId, apiUrl, err := getAccountStuff()
	if err != nil {
		fmt.Printf("error getting account stuff: %v", err)
	}
	log.Print(accountId)
	log.Print(apiUrl)
	j := fmt.Sprintf(`{
        "using": ["urn:ietf:params:jmap:core", "urn:ietf:params:jmap:mail"],
        "methodCalls": [
            [
                "Mailbox/query",
                {
                    "accountId": "%s",
                    "filter": {"role": "inbox", "hasAnyRole": true}
                },
                "a"
            ]
        ]
}`, accountId)
	inbox_res, err := makeJmapCall(apiUrl, j)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(inbox_res), &respMap)
	if err != nil {
		fmt.Printf(": %v", err)
	}

	methodResponses := respMap["methodResponses"].([]interface{})
	mailboxQuery := methodResponses[0].([]interface{})
	queryResponse := mailboxQuery[1].(map[string]interface{})
	inboxId := queryResponse["ids"].([]interface{})[0]

	j2 := fmt.Sprintf(`{
        "using": ["urn:ietf:params:jmap:core", "urn:ietf:params:jmap:mail"],
        "methodCalls": [
            [
                "Email/query",
                {
                    "accountId": "%s",
                    "filter": {"inMailbox": "%s"},
                    "sort": [{"property": "receivedAt", "isAscending": false}],
                    "limit": 10
                },
                "a"
            ],
            [
                "Email/get",
                {
                    "accountId": "%s",
                    "properties": ["id", "subject", "receivedAt"],
                    "#ids": {"resultOf": "a", "name": "Email/query", "path": "/ids/*"}
                },
                "b"
            ]
        ]
}`, accountId, inboxId, accountId)

	jmapResp2, err := makeJmapCall(apiUrl, j2)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	resp2Map := make(map[string]interface{})
	err = json.Unmarshal([]byte(jmapResp2), &resp2Map)
	if err != nil {
		fmt.Printf(": %v", err)
	}

	methodResponses2 := resp2Map["methodResponses"].([]interface{})
	queryResponse2 := methodResponses2[1].([]interface{})
	something := queryResponse2[1].(map[string]interface{})
	emailList := something["list"].([]interface{})
	for _, email := range emailList {
		e := email.(map[string]interface{})
		fmt.Printf("%s - %s\n", e["subject"], e["receivedAt"])
	}

}
