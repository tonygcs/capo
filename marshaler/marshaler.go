package marshaler

var mars Marshaler = &JSONMarshaler{}

// GetMarshaler returns the current marshaler instance.
func GetMarshaler() Marshaler {
	return mars
}

// SetMarshaler sets the entity that will transform the requests and responses
// in the entities that the application can handle.
func SetMarshaler(m Marshaler) {
	mars = m
}

// Marshaler is the entity that will transform the HTTP requests and responses.
type Marshaler interface {
	ContentTypeHeader() string
	Marshal(entity interface{}) ([]byte, error)
	Unmarshal(data []byte, entity interface{}) error
}
