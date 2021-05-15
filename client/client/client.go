package client

import (
	"app/client/client_config"
	"app/http_types"
	"app/payload"
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"net/http"
	"time"
)

type Client struct {
	Decoded    []*payload.Payload
	Signatures []string
	httpClient *http.Client // TODO use http3
	address *common.Address
	publicKey *common.Address
}

func CreateClient(parsedConfig *client_config.Config) (*Client, error) {
	// size int, address common.Address, publicKey *ecdsa.PublicKey
	size := parsedConfig.PayloadsPerConnection + 1
	address := common.HexToAddress(parsedConfig.Address)
	publicKey := common.HexToAddress(parsedConfig.PublicKey)

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	c := http.Client{Transport: tr}
	return &Client{httpClient: &c,
				   Decoded: make([]*payload.Payload, size),
				   Signatures: make([]string, size),
				   address: &address,
	               publicKey: &publicKey}, nil
}

func (c *Client) SavePayload(i int, resp *http.Response) error {

	t, err := http_types.DecodePayloadBody(resp.Body)
	if err != nil {
		return err
	}
	// TODO check signature
	log.Println(t.Payload)
	c.Decoded[i] = t.Payload
	c.Signatures[i] = t.Signature
	return nil
}
// TODO add check for correct server response

func (c *Client) Create(i int) error {
	body, err := http_types.CreateCreateBody(c.address.Hex(), i)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Post("http://127.0.0.1:9000/create", "", bytes.NewReader(body))
	if err != nil {
		return err
	}
	err = c.SavePayload(i, resp)
	return err
}

func (c *Client) Increment(i int) error {
	body, err := http_types.CreatePayloadBody(&http_types.Response{Payload: c.Decoded[i], Signature: c.Signatures[i]})
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Post("http://127.0.0.1:9000/increment", "", bytes.NewReader(body))
	if err != nil {
		return err
	}

	err = c.SavePayload(i, resp)
	return err
}

// writes result to 0 index
func (c *Client) Merge() error {
	body, err := http_types.CreateMergeBody(c.Decoded[1:], c.Signatures[1:])
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Post("http://127.0.0.1:9000/merge", "", bytes.NewReader(body))
	//log.Println("Body = ", resp.Body)
	if err != nil {
		return err
	}

	err = c.SavePayload(0, resp)
	return err
}

func (c *Client) Loop() error { // TODO use this function ?
	return nil
}
