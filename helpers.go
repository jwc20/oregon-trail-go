package oregontrail

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const RandomIntURL = "https://www.random.org/integers/?num=1&min=1&max=10&col=1&base=10&format=plain&rnd=new"

func GetRandomInt() int {
	request := NewGetRandomIntRequest()
	response := NewGetRandomIntResponseFromClient(request)
	defer response.Body.Close()

	result := ExtractRandomInteger(response)
	return result
}

func NewGetRandomIntRequest() *http.Request {
	req, err := http.NewRequest("GET", RandomIntURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	return req
}

func NewGetRandomIntResponseFromClient(req *http.Request) *http.Response {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func ExtractRandomInteger(resp *http.Response) int {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	result, err := strconv.Atoi(strings.TrimSpace(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	return result
}
