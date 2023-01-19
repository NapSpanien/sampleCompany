package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// MyJSONResponse respuesta JSON standard
type MyJSONResponse struct {
	MyResponse
	status      int
	contentType string
	writer      http.ResponseWriter
}

// MyResponse Estructura interna de la respuesta
type MyResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"body"`
	Msg  string      `json:"message"`
}

// CreateJSONResponse crea una respuesta a partir de r MyResponse en w
func createJSONResponse(w http.ResponseWriter, r MyResponse, statusCode int) MyJSONResponse {
	jsonResponse := MyJSONResponse{status: statusCode, contentType: "application/json", writer: w}
	jsonResponse.MyResponse.Data = r.Data
	jsonResponse.MyResponse.Code = r.Code
	jsonResponse.MyResponse.Msg = r.Msg
	return jsonResponse
}

// SendJSONResponse escribe la respuesta JSON en el writer
func (my *MyJSONResponse) SendJSONResponse() {
	my.writer.Header().Set("Content-Type", my.contentType)
	my.writer.WriteHeader(my.status)

	output, _ := json.Marshal(&my)
	fmt.Fprintf(my.writer, string(output))
}

// GenerateResponse GenerateResponse
func GenerateResponse(w http.ResponseWriter, r MyResponse, status int) {
	MyJSONResponse := createJSONResponse(w, r, status)
	MyJSONResponse.SendJSONResponse()
}
