package post

import (
	"encoding/json"
	"net/http"

	"google.golang.org/appengine"

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

	hash := request.URL.Query().Get("hash")

	query := model.Transaction{
		Hash: hash,
	}

	var tx model.Transaction

	client := query.GetRecordByHash(&tx)
	if client.RecordNotFound() {
		message := "No records for transaction table"
		badRequest.Message = message
		badRequest.OnBadRequest(http.StatusInternalServerError)
		return
	}

	row := client.Row()
	err := row.Scan(&tx.Hash, &tx.From)
	if err != nil {
		message := err.Error()
		badRequest.Message = message
		badRequest.OnBadRequest(http.StatusInternalServerError)
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
