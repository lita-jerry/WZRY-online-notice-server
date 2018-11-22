package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	url := "http://127.0.0.1:9090/login?username=zhao&password=123"

	req, _ := http.NewRequest("POST", url, nil)

	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "a880ca77-3976-bd28-1329-9b2ecb7e049d")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
