package capo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// jsonEncoder is a JSON encoder.
func jsonEncoder(data any) ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)

	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal the data :: %w", err)
	}

	return buf.Bytes(), nil
}
