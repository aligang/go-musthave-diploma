package accural

import (
	"bytes"
	"context"
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

func FetchOrderInfo(ctx context.Context, orderID string, config *config.Config) (*message.AccuralMessage, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	buf := &bytes.Buffer{}
	URI := fmt.Sprintf("%s/api/orders/%s", config.AccuralSystemAddress, orderID)
	request, err := http.NewRequest("GET", URI, buf)
	request = request.WithContext(ctx)
	if err != nil {
		logging.Warn("Error During Request preparation: %s", err.Error())
		return nil, err
	}

	requestDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		logging.Warn("Error During Request Dump: %s", err.Error())
		return nil, err
	}
	logging.Debug("Sending request to: URI: %s", URI)
	logging.Debug("request content: %s", string(requestDump))
	fmt.Println(string(requestDump))

	select {
	default:
	case <-ctx.Done():
		return nil, errors.New("context was stopped")
	}
	//request.Header.Add("Accept-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		logging.Warn("Error During communication with: %s", URI)
		return nil, err
	}
	defer response.Body.Close()
	responseDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		logging.Warn("Error During Response Dump: %s", err.Error())
		return nil, err
	}
	logging.Debug("Response content: %s", string(responseDump))
	fmt.Println(string(responseDump))

	accuralRecord := message.New()
	switch {
	case response.StatusCode == http.StatusNoContent:

		accuralRecord.Order = orderID
		accuralRecord.Status = "NEW"
	case response.StatusCode != http.StatusOK:
		logging.Warn("Got response from %s with code: %d, Could not fetch order info",
			config.AccuralSystemAddress, response.StatusCode)
		return nil, errors.New("problem during fetching order info")
	default:
		responsePayload, err := io.ReadAll(response.Body)
		if err != nil {
			logging.Warn("Could not read data from wire")
			return nil, err
		}
		logging.Warn("got accural record %s", string(responsePayload))
		err = json.Unmarshal(responsePayload, accuralRecord)
		if err != nil {
			logging.Warn("Could not decode json")
			return nil, err
		}
	}
	return accuralRecord, nil
}
