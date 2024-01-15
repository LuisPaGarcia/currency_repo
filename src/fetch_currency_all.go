// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var urls = []string{
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/gtq/clp.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/gtq/clp.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/gtq/cop.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/gtq/pen.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/gtq/mxn.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/gtq/usd.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/gtq/crc.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/cop/clp.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/cop/gtq.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/cop/pen.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/cop/mxn.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/cop/usd.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/cop/crc.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/pen/clp.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/pen/gtq.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/pen/cop.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/pen/mxn.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/pen/usd.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/pen/crc.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/mxn/clp.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/mxn/gtq.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/mxn/cop.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/mxn/pen.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/mxn/usd.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/mxn/crc.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/usd/clp.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/usd/gtq.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/usd/cop.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/usd/pen.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/usd/mxn.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/usd/crc.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/crc/clp.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/crc/gtq.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/crc/cop.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/crc/pen.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/crc/mxn.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/crc/usd.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/clp/crc.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/clp/gtq.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/clp/cop.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/clp/pen.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/clp/mxn.json",
	"https://raw.githubusercontent.com/fawazahmed0/currency-api/1/latest/currencies/clp/usd.json",
}

func fetchAPI(url string, wg *sync.WaitGroup, results chan<- map[string]interface{}) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", url, err)
		return
	}

	var data interface {
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", url, err)
		return
	}

	urlParts := strings.Split(url, "/")
	propName := strings.TrimSuffix(urlParts[len(urlParts)-1], ".json")
	key := urlParts[len(urlParts)-2] + "_" + strings.TrimSuffix(propName, ".json")
	results <- map[string]interface{}{key: data.(map[string]interface{})[propName]}
}

func main() {
	var wg sync.WaitGroup
	results := make(chan map[string]interface{}, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go fetchAPI(url, &wg, results)
	}

	wg.Wait()
	close(results)

	finalResults := make(map[string]interface{})
	for result := range results {
		for k, v := range result {
			finalResults[k] = v
		}
	}

	location, err := time.LoadLocation("America/Guatemala")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	// Agregar la fecha actual en la zona horaria de Guatemala en la raíz del objeto
	currentDate := time.Now().In(location).Format("2006-01-02 15:04:05")
	finalResults["date"] = currentDate

	file, err := json.MarshalIndent(finalResults, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling results to JSON:", err)
		return
	}

	err = ioutil.WriteFile("./currency_all/currency_rates.json", file, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}
