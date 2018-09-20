package post

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"gitlab.com/canya-com/canwork-database-client/model"
	HTTP "gitlab.com/canya-com/canya-ethereum-tx-api/internal"
)

// MonitorTransaction : POST, get tx by id
func MonitorTransaction(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var badRequest = HTTP.BadRequest{
		Context: appengine.NewContext(request),
		Writer:  writer,
	}

	writer.Header().Set("Content-Type", "application/json")

	if request.Method != HTTP.POST {
		message := HTTP.NewInvalidMethodMessage(HTTP.POST, request)
		badRequest.Message = message
		badRequest.OnBadRequest(http.StatusForbidden)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		message := err.Error()
		badRequest.Message = message
		badRequest.OnBadRequest(http.StatusInternalServerError)
		return
	}

	now := time.Now()
	tx := model.Transaction{
		Timeout: now.AddDate(0, 0, 1).Unix(),
	}

	err = json.Unmarshal(body, &tx)
	if err != nil {
		message := err.Error()
		badRequest.Message = message
		badRequest.OnBadRequest(http.StatusInternalServerError)
		return
	}

	if !tx.IsValid() {
		message := fmt.Sprintf("Invalid transaction hash. Format must have 0x prefix and be of length %d", tx.Length()*2)
		badRequest.Message = message
		badRequest.OnBadRequest(http.StatusInternalServerError)
		return
	}

	if tx.RecordExists() {
		message := "Transaction record exists"
		tx.GetRecordByHash(&tx).Row().Scan(&tx.Hash, &tx.From)
		badRequest.Transaction = tx
		badRequest.Message = message
		badRequest.OnBadRequest(http.StatusInternalServerError)
		return
	}

	tx.CreatedAt = now.Unix()

	databaseClient := tx.New()

	err = databaseClient.Error
	if err != nil {
		errors := databaseClient.GetErrors()
		group := []HTTP.BadRequest{}
		for _, err := range errors {
			log.Errorf(badRequest.Context, err.Error())
			group = append(group, HTTP.BadRequest{
				Message: err.Error(),
			})
		}
		response, _ := json.Marshal(group)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(response)
		return
	}

	output, err := json.Marshal(tx)
	if err != nil {
		message := "GET GetTransaction failed to Marshal transaction"
		badRequest.Message = message
		badRequest.OnBadRequest(http.StatusInternalServerError)
		return
	}

	writer.Write(output)
}
