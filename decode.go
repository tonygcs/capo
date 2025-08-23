package capo

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// jsonDecoder is a JSON decoder.
func jsonDecoder(body []byte, data any) error {
	err := json.NewDecoder(bytes.NewBuffer(body)).Decode(data)
	if err != nil {
		return Error(err, InvalidRequestBodyCode, InvalidRequestBodyMsg, http.StatusBadRequest)
	}
	return nil
}
