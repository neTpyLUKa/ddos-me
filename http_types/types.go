package http_types

import (
	"app/payload"
	"encoding/json"
	"io"
)

type Response struct {
	*payload.Payload
	Signature string `json:"Signature"`
}

type CreateBody struct {
	PublicKey string `json:"PublicKey"`
	Index int `json:"Index"`
}

type MergeBody struct {
	Payloads []*payload.Payload `json:"Payloads"`
	Signatures []string `json:"Signatures"`
}

func CreatePayloadBody(resp *Response) ([]byte, error) {
	result, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DecodePayloadBody(body io.ReadCloser) (*Response, error) {
	decoder := json.NewDecoder(body)
	t := Response{}
	err := decoder.Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func CreateCreateBody(p string, index int) ([]byte, error) {
	result, err := json.Marshal(CreateBody{PublicKey: p, Index: index})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DecodeCreateBody(body io.ReadCloser) (*CreateBody, error) {
	decoder := json.NewDecoder(body)
	t := CreateBody{}
	err := decoder.Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func CreateMergeBody(payloads []*payload.Payload, signatures []string) ([]byte, error) {
	result, err := json.Marshal(MergeBody{Payloads: payloads, Signatures: signatures})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DecodeMergeBody(body io.ReadCloser) (*MergeBody, error) {
	decoder := json.NewDecoder(body)
	t := MergeBody{}
	err := decoder.Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

