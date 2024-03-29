package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Define una estructura para mapear la respuesta JSON.
type ApiResponse struct {
	Result string   `json:"Result"`
	Data   []string `json:"result"`
	Date   string   `json:"date"`
}

func main() {
	urlStr := "https://www.corporacionbi.com/gt/bancoindustrial/wp-content/plugins/jevelin_showcase_exchange_rate/service/mod_moneda.php"
	method := "POST"

	payload := strings.NewReader("action=getMoneda")

	client := &http.Client{}
	req, err := http.NewRequest(method, urlStr, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Language", "es-419,es;q=0.6")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Origin", "https://www.corporacionbi.com")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-GPC", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("sec-ch-ua", `"Chromium";v="116", "Not)A;Brand";v="24", "Brave";v="116"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", `"macOS"`)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Deserializa la respuesta JSON.
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	location, err := time.LoadLocation("America/Guatemala")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	// Agregar la fecha actual en la zona horaria de Guatemala en la raíz del objeto
	currentDate := time.Now().In(location).Format("2006-01-02 15:04:05")
	// Agrega la nueva propiedad con un string vacío.
	apiResponse.Date = currentDate

	// Serializa el objeto modificado a JSON.
	modifiedJSON, err := json.Marshal(apiResponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile("./currency_all/tipo_de_cambio_bi.json", modifiedJSON, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
