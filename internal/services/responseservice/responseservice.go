package responseservice

// Payload represents a standardised JSON payload returned by API
type Payload struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// NewPayload returns a new ResponsePayload
func NewPayload(s string, d interface{}) Payload {
	return Payload{
		Status: s,
		Data:   d,
	}
}
