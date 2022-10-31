package marshaler

import "encoding/json"

type JSONMarshaler struct {
}

func (m *JSONMarshaler) ContentTypeHeader() string {
	return "application/json"
}

func (m *JSONMarshaler) Unmarshal(data []byte, entity interface{}) error {
	return json.Unmarshal(data, entity)
}

func (m *JSONMarshaler) Marshal(entity interface{}) ([]byte, error) {
	return json.Marshal(entity)
}
