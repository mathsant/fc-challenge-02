package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type CDNViaCepResponse struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func makeRequest(url string, result chan<- *http.Response) {
	response, err := http.Get(url)
	if err != nil {
		result <- nil
		return
	}
	result <- response
}

func main() {
	channel1 := make(chan *http.Response)
	channel2 := make(chan *http.Response)

	url1 := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", "36904-278")
	url2 := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", "36904278")

	go makeRequest(url1, channel1) // CDN
	go makeRequest(url2, channel2) // Via Cep

	select {
	case msg := <-channel1: // cdn
		var data CDNViaCepResponse
		json.NewDecoder(msg.Body).Decode(&data)
		fmt.Printf("Received from CDN: ID: %s | Dados Resgatados: %v\n", msg.Request.URL, data)
	case msg := <-channel2: // via cep
		var data ViaCepResponse
		json.NewDecoder(msg.Body).Decode(&data)
		fmt.Printf("Received from Via CEP: URL: %s | Dados Resgatados: %v\n", msg.Request.URL, data)
	case <-time.After(time.Second * 1):
		println("Timeout!")
	}
}
