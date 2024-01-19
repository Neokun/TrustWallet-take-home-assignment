package api

import (
	"encoding/json"
	"net/http"

	"github.com/TrustWallet/tx-parser/internal/parser"
)

type register struct {
	parserSvc parser.Parser
}

func NewRegister(parserSvc parser.Parser) *register {
	reg := &register{
		parserSvc: parserSvc,
	}

	return reg
}

func (reg *register) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Address string `json:"address"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	if !reg.parserSvc.Subscribe(data.Address) {
		http.Error(w, "Address already subscribed", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Subscribed successfully"))
}

func (reg *register) GetCurrentBlockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is accepted", http.StatusMethodNotAllowed)
		return
	}

	blockNum := reg.parserSvc.GetCurrentBlock()
	response, err := json.Marshal(map[string]interface{}{"block": blockNum})
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (reg *register) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is accepted", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is missing", http.StatusBadRequest)
		return
	}

	txns := reg.parserSvc.GetTransactions(address)
	response, err := json.Marshal(txns)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
