package accural

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/accural/message"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

func FetchOrderInfo(orderId string, config *config.Config) (*message.AccuralMessage, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	buf := &bytes.Buffer{}
	URI := fmt.Sprintf("%s/api/orders/%s", config.AccuralSystemAddress, orderId)
	request, err := http.NewRequest("GET", URI, buf)
	if err != nil {
		logging.Warn("Error During Request preparation: %s", err.Error())
		return nil, err
	}
	requestDump, err := httputil.DumpRequestOut(request, true)
	logging.Debug("Sending request to: URI: %s", URI)
	logging.Debug("request content: %s", string(requestDump))
	fmt.Println(string(requestDump))
	//request.Header.Add("Accept-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		logging.Warn("Error During communication with: %s", URI)
		return nil, err
	}
	responseDump, err := httputil.DumpResponse(response, true)
	logging.Debug("Response content: %s", string(responseDump))
	fmt.Println(string(responseDump))
	if response.StatusCode != http.StatusOK {
		logging.Warn("Got response from %s with code: %d, Could not fetch order info",
			config.AccuralSystemAddress, response.StatusCode)
		return nil, errors.New("problem during fetching order info")
	}

	responsePayload, err := io.ReadAll(response.Body)
	if err != nil {
		logging.Warn("Could not read data from wire")
		return nil, err
	}
	logging.Warn("got accural record %s", string(responsePayload))
	accuralRecord := &message.AccuralMessage{}
	err = json.Unmarshal(responsePayload, accuralRecord)
	if err != nil {
		logging.Warn("Could not decode json")
		return nil, err
	}
	return accuralRecord, nil
}
