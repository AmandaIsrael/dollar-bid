package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type BidResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("[CLIENT] Erro ao criar requisição: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("[CLIENT] Timeout de 300ms excedido ao buscar cotação")
		} else {
			log.Printf("[CLIENT] Erro na requisição HTTP: %v", err)
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[CLIENT] Servidor retornou erro: status %d", resp.StatusCode)
		return
	}

	var bidResp BidResponse
	err = json.NewDecoder(resp.Body).Decode(&bidResp)
	if err != nil {
		log.Printf("[CLIENT] Erro ao decodificar resposta: %v", err)
		return
	}

	writeFile(fmt.Sprintf("Dólar: %s", bidResp.Bid))
	log.Printf("[CLIENT] Cotação salva no arquivo cotacao.txt: %s", bidResp.Bid)
}

func writeFile(content string) {
	err := os.WriteFile("cotacao.txt", []byte(content), 0644)
	if err != nil {
		log.Fatalf("Erro ao escrever no arquivo: %v", err)
	}
}
