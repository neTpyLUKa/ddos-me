package server

import (
	"app/payload"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type HttpHandler struct {
	secret *payload.Secret
}

type IncrementBody struct {
	Payload string `json:"Payload"`
}

func (h *HttpHandler) SendPayload(res http.ResponseWriter, p *payload.Payload) error {
	bytes, err := payload.Encode(p, h.secret)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return err
	}
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(bytes)
	return err
}

func (h *HttpHandler) CreateHandler(res http.ResponseWriter, req *http.Request) {
	// TODO call close ?
	log.Printf("%+v\n", req)
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	p, err := payload.CreatePayload(h.secret)
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
	log.Printf("%+v\n", req)
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(req.Body)
	t := IncrementBody{}
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Error decoding HTTP body:", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println(t.Payload)
	p, err := payload.Decode([]byte(t.Payload))
	if err != nil {
		log.Println("Error decoding payload:", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	p.Counter++

	err = h.SendPayload(res, p)
	if err != nil {
		log.Println("Error sending payload:", err.Error())
		return
	}
	return
}

type MergeBody struct {
	Payloads []string `json:"Payloads"`
}

func (h *HttpHandler) MergeHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("%+v\n", req)
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(req.Body)
	t := MergeBody{}
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Error decoding HTTP body:", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println(t.Payloads)
	var payloads []*payload.Payload

	for _, rawP := range t.Payloads {
		p, err := payload.Decode([]byte(rawP))
		if err != nil {
			log.Println("Error decoding payload:", err.Error())
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		payloads = append(payloads, p)
	}

	result, err := payload.MergePayload(payloads)
	if err != nil {
		log.Println("Error merging payloads:", err.Error())
		res.WriteHeader(http.StatusForbidden)
	}

	err = h.SendPayload(res, result)
	if err != nil {
		log.Println("Error sending payload:", err.Error())
		return
	}
	return
}

func StartServer() {
	secret := payload.Secret{}
	handler := HttpHandler{secret: &secret}
	r := mux.NewRouter()
	r.HandleFunc("/create", handler.CreateHandler)
	r.HandleFunc("/increment", handler.IncrementHandler)
	r.HandleFunc("/merge", handler.MergeHandler)
	log.Println("Starting HTTP Server...")
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:9000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
