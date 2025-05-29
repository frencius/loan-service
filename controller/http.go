package controller

import (
	"encoding/json"
	"net/http"

	"github.com/frencius/loan-service/model"
)

func WriteHTTPResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		respCode := http.StatusInternalServerError
		result := model.ComposeErrorResponse(respCode, err.Error(), "JSON marshal failed")
		WriteHTTPResponse(w, respCode, result)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(bodyBytes)
}

func WriteResponseFile(w http.ResponseWriter, fileName, fileType string, response []byte, status int) {
	w.Header().Set("Content-Type", "text/"+fileType)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)

	w.WriteHeader(status)

	_, _ = w.Write(response)
}
