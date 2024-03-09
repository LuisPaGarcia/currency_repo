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
	"https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/usd.json",
	"https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/gtq.json",
	"https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/crc.json",
	"https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/mxn.json",
	"https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/clp.json",
	"https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/cop.json",
	"https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/pen.json",
}

func main() {
	var wg sync.WaitGroup
	resultsChan := make(chan map[string]interface{}, len(urls)) // Create a channel to collect results

	for _, url := range urls {
		wg.Add(1)
		go fetchCurrency(url, &wg, resultsChan) // Pass the channel as an argument
	}

	go func() {
		wg.Wait()
		close(resultsChan) // Close the channel once all goroutines are done
	}()

	finalResult := make(map[string]interface{}) // Map to combine all results
	for result := range resultsChan {
		for key, value := range result {
			finalResult[key] = value // Merge result into finalResult map
		}
	}

	location, err := time.LoadLocation("America/Guatemala")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}
	// Agregar la fecha actual en la zona horaria de Guatemala en la raíz del objeto
	currentDate := time.Now().In(location).Format("2006-01-02 15:04:05")
	finalResult["timestamp"] = currentDate

	// Convert the combined results to JSON
	finalJSON, err := json.MarshalIndent(finalResult, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling final JSON:", err)
		return
	}

	err = ioutil.WriteFile("./currency_all/currency_rates.json", finalJSON, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}

	fmt.Println("Fetch completed.")
}

// Adjusted fetchCurrency function to send results to a channel
func fetchCurrency(url string, wg *sync.WaitGroup, results chan<- map[string]interface{}) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(err)
		return
	}

	key := strings.Split(url, "/")[8]
	key = strings.Split(key, ".")[0]

	if currencyData, ok := data[key].(map[string]interface{}); ok {
		result := make(map[string]interface{})
		// Especifica las keys deseadas
		desiredKeys := []string{"usd", "gtq", "crc", "mxn", "clp", "cop", "pen"}
		for _, desiredKey := range desiredKeys {
			if value, exists := currencyData[desiredKey]; exists {
				result[key+"_"+desiredKey] = value
			}
		}
		results <- result // Send the result to the channel
	} else {
		fmt.Println("Invalid currency data for URL:", url)
	}
}
