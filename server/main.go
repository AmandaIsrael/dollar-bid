package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", DollarBidHandler)

	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func DollarBidHandler(w http.ResponseWriter, r *http.Request) {
	usdbrl := &Usdbrl{}
	result := usdbrl.GetDollarBid(r.Context())

	if result == "" {
		log.Println("[SERVER] Erro ao obter cotação")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erro interno do servidor"})
		return
	}

	log.Println("[SERVER] Cotação retornada com sucesso")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"bid": result})
}
