package payload

import (
	"testing"
)

func TestEncode(t *testing.T) {
	secret := Secret{}
	payload := Payload{Counter: 0, Id: []int64{1}}
	bytes, err := Encode(&payload, &secret)
	if err != nil {
		t.Error("Error encoding payload:", err.Error())
	}
	if string(bytes) != "{\"Counter\":0,\"Id\":[1]}" {
		t.Error("Bad encoding:", string(bytes))
	}
}

func TestDecode(t *testing.T) {
	bytes := []byte("{\"Counter\":5, \"Id\":[0]}")
	payload, err := Decode(bytes)
	if err != nil {
		t.Error("Error decoding payload:", err.Error())
	}
	if payload == nil || payload.Counter != 5 || len(payload.Id) != 1 || payload.Id[0] != 0 {
		t.Error("Unexpected value decoded")
	}
}