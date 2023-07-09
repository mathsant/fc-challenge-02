package main

import (
	"fmt"
	"net/http"
	"time"
)

func makeRequest(url string, result chan<- *http.Response) {
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	response, err := client.Get(url)
	if err != nil {
		result <- nil
		return
	}

	result <- response
}

func fetchData(cep string) *http.Response {
	url1 := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep)
	url2 := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	result := make(chan *http.Response, 2)

	go makeRequest(url1, result)
	go makeRequest(url2, result)

	for i := 0; i < 2; i++ {
		response := <-result
		if response != nil {
			return response
		}
	}

	return nil
}

func main() {
	cep := "36904-278"

	response := fetchData(cep)

	if response != nil {
		fmt.Printf("API: %s\n", response.Request.URL)
	} else {
		fmt.Println("Nenhuma resposta recebida dentro do tempo limite.")
	}
}
