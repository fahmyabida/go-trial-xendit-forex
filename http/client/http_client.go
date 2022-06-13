package client

import (
	"encoding/json"
	"fmt"
	"go-forex/model"
	"go-forex/utils"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type HttpClient struct {
	svcConfig model.SvcConfig
}

func NewHttpClient(svcConfig model.SvcConfig) HttpClient {
	return HttpClient{svcConfig}
}

func (c HttpClient) GetForex(traceId, originCurrency, destinationCurrency string) (response model.RateData) {
	url := c.svcConfig.URLGetRate
	method := c.svcConfig.URLGetRateMethod

	payload := strings.NewReader(`{
		"origin_currency": "` + originCurrency + `",
		"destination_currency": "` + destinationCurrency + `"
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Err(fmt.Errorf("'%v', %v", traceId, err))
		return
	}
	utils.LogOUT(traceId, method, url)
	var bodyResponse []byte = []byte("empty body response")
	defer func() {
		utils.LogIN(traceId, method, url, string(bodyResponse))
	}()
	res, err := client.Do(req)
	if err != nil {
		log.Err(fmt.Errorf("'%v', %v", traceId, err))
		return
	}
	defer res.Body.Close()

	bodyResponse, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Err(fmt.Errorf("'%v', %v", traceId, err))
		return
	}
	json.Unmarshal(bodyResponse, &response)
	return response
}
