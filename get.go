package main

import (
	"encoding/json"
	"net/http"

	"google.golang.org/appengine"

	"gitlab.com/canya-com/canwork-database-client/model"
)

// GetRequest : GET
type GetRequest struct {
	*http.Request
}

// MonitorTransaction : tracks all pending or on_timeout transactions in the database
func (r *GetRequest) MonitorTransaction() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		var badRequest = BadRequest{
			Context: appengine.NewContext(request),
			Writer:  writer,
		}

		writer.Header().Set("Content-Type", "application/json")

		if request.Method != GET {
			message := NewInvalidMethodMessage(GET, request)
			badRequest.Message = message
			badRequest.OnBadRequest(http.StatusForbidden)
			return
		}

		var tx model.Transaction
		txs := []model.Transaction{}

		databaseClient := tx.Table()
		tx.IsWebhookCalled = 0
		databaseClient.Find(&txs, tx)

		rows, err := databaseClient.Rows()
		if err != nil {
			message := err.Error()
			badRequest.Message = message
			badRequest.OnBadRequest(http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			err := rows.Scan(&tx.Hash, &tx.From, &tx.Status, &tx.Network, &tx.CreatedAt, &tx.CompletedAt, &tx.Timeout, &tx.IsWebhookCalled, &tx.WebhookOnSuccess, &tx.WebhookOnTimeout)
			if err != nil {
				message := err.Error()
				badRequest.Message = message
				badRequest.OnBadRequest(http.StatusInternalServerError)
				return
			}
			txs = append(txs, tx)
		}

		output, err := json.Marshal(txs)
		if err != nil {
			message := err.Error()
			badRequest.Message = message
			badRequest.OnBadRequest(http.StatusInternalServerError)
			return
		}

		writer.Write(output)
	}
}

// TransactionDetails : GET tx by id
func (r *GetRequest) TransactionDetails() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		var badRequest = BadRequest{
			Context: appengine.NewContext(request),
			Writer:  writer,
		}

		writer.Header().Set("Content-Type", "application/json")

		if request.Method != GET {
			message := NewInvalidMethodMessage(GET, request)
			badRequest.Message = message
			badRequest.OnBadRequest(http.StatusForbidden)
			return
		}

		hash := request.URL.Query().Get("hash")

		query := model.Transaction{
			Hash: hash,
		}

		var tx model.Transaction

		databaseClient := query.GetRecordByHash(&tx)
		if databaseClient.RecordNotFound() {
			message := "No records on transaction table"
			badRequest.Message = message
			badRequest.OnBadRequest(http.StatusInternalServerError)
			return
		}

		row := databaseClient.Row()
		err := row.Scan(&tx.Hash, &tx.From, &tx.Status, &tx.Network, &tx.CreatedAt, &tx.CompletedAt, &tx.Timeout, &tx.IsWebhookCalled, &tx.WebhookOnSuccess, &tx.WebhookOnTimeout)
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
}
