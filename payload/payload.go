package payload

import (
	"encoding/json"
	"math/rand"
	"time"
)


func init() {
	rand.Seed(time.Now().UnixNano())
}

type Payload struct {
	Counter int64 `json:"Counter"`
	Id []int64 `json:"Id"`
}

func CreatePayload(secret *Secret) (*Payload, error) {
	return &Payload{Counter: 0, Id: []int64{rand.Int63()}}, nil
}

func Encode(payload *Payload, secret *Secret) ([]byte, error) {
	// TODO use secret
	b, err := json.Marshal(payload)
	return b, err
}

func Decode(bytes []byte) (*Payload, error) {
	payload := &Payload{}
	err := json.Unmarshal(bytes, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func MergePayload(payloads []*Payload) (*Payload, error) {
	// TODO add check for duplicates ?
	result := Payload{}
	for _, payload := range payloads {
		result.Id = append(result.Id, payload.Id...)
		result.Counter += payload.Counter
	}
	return &result, nil
}