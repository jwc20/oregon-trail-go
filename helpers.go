package oregontrail

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const RandomIntURL = "https://www.random.org/integers/?num=1&min=1&max=10&col=1&base=10&format=plain&rnd=new"

func GetRandomInt() int {
	request := NewGetRandomIntRequest()
	response, err := NewGetRandomIntResponseFromClient(request)
	if err != nil {
		// backup
		result := rand.Intn(10)
		return result
	}
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

func NewGetRandomIntResponseFromClient(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp, err
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
