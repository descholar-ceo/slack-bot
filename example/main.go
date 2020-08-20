package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Res is from internet
type Res map[string]interface{}

func main() {

	var result Res
	resp, err := http.Get("")
	if err != nil {
		fmt.Printf("Ooops! Something went wrong %v\n", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	// isValid := json.Valid(body)
	fmt.Println(result)
}
