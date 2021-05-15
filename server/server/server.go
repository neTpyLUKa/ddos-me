package server

import (
	"app/http_types"
	"app/payload"
	"app/server/server_config"
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type HttpHandler struct {
	privateKey *ecdsa.PrivateKey
}

func (h *HttpHandler) SendPayload(res http.ResponseWriter, p *payload.Payload) error {
	signature, err := payload.CreateSignature(p, h.privateKey)
	if err != nil {
		return err
	}
	resp := http_types.Response{Payload: p, Signature: string(signature)}
	result, err := json.Marshal(resp)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return err
	}
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(result)
	return err
}

func (h *HttpHandler) CreateHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("%+v\n", req.Body)
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	t, err := http_types.DecodeCreateBody(req.Body)
	if err != nil {
		log.Println("Error decoding HTTP body:", err.Error())
		res.WriteHeader(http.StatusBadRequest)
	}

	p, err := payload.CreatePayload(common.HexToAddress(t.PublicKey), int64(t.Index), h.privateKey)
	if err != nil {
		log.Println("Error creating payload:", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.SendPayload(res, p)
	if err != nil {
		log.Println("Error sending payload:", err.Error())
		return
	}
	return
}

func (h *HttpHandler) IncrementHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resp, err := http_types.DecodePayloadBody(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO check signature
	resp.Counter++

	log.Printf("%+v\n", resp.Payload)

	err = h.SendPayload(res, resp.Payload)
	if err != nil {
		log.Println("Error sending payload:", err.Error())
		return
	}
	return
}

func (h *HttpHandler) MergeHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("%+v\n", req)
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	t, err := http_types.DecodeMergeBody(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := payload.MergePayload(t.Payloads)
	if err != nil {
		log.Println("Error merging payloads:", err.Error())
		res.WriteHeader(http.StatusForbidden)
		return
	}

	err = h.SendPayload(res, result)
	if err != nil {
		log.Println("Error sending payload:", err.Error())
		return
	}
	return
}

func StartServer(parsedConfig *server_config.Config) {
	privateKey, err := crypto.HexToECDSA(parsedConfig.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	handler := HttpHandler{privateKey: privateKey}
	r := mux.NewRouter()
	r.HandleFunc("/create", handler.CreateHandler)
	r.HandleFunc("/increment", handler.IncrementHandler)
	r.HandleFunc("/merge", handler.MergeHandler)
	log.Println("Starting HTTP Server...")
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:9000",
		WriteTimeout: 15 * time.Second, // TODO add to server_config
		ReadTimeout:  15 * time.Second, // TODO add to server_config
	} // TODO use connection limiter, use parameter from server_config
	log.Fatal(srv.ListenAndServe())
}
