package oregontrail

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const url = "https://www.random.org/integers/?num=1&min=1&max=10&col=1&base=10&format=plain&rnd=new"

func GetRandomInt() int {
	request := newGetRandomIntRequest()
	response := newGetRandomIntResponseFromClient(request)
	defer response.Body.Close()

	result := extractRandomInteger(response)
	return result
}

func newGetRandomIntRequest() *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	return req
}

func newGetRandomIntResponseFromClient(req *http.Request) *http.Response {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func extractRandomInteger(resp *http.Response) int {
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
