package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ApiResponse struct {
	USDBRL Usdbrl `json:"USDBRL"`
}

type Usdbrl struct {
	gorm.Model `json:"-"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

const (
	UsdbrlUrl       = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	ApiTimeout      = 200 * time.Millisecond
	DatabaseTimeout = 10 * time.Millisecond
)

func (u *Usdbrl) GetDollarBid(ctx context.Context) string {
	ctx, cancel := context.WithTimeout(ctx, ApiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", UsdbrlUrl, nil)
	if err != nil {
		log.Printf("[SERVER] Erro ao criar requisição: %v\n", err)
		return ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("[SERVER] Timeout de 200ms excedido ao buscar cotação da API externa")
		} else {
			log.Printf("[SERVER] Erro na requisição HTTP: %v\n", err)
		}
		return ""
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[SERVER] Erro ao ler resposta da API externa: %v\n", err)
		return ""
	}

	var apiResp ApiResponse
	err = json.Unmarshal(result, &apiResp)
	if err != nil {
		log.Printf("[SERVER] Erro ao decodificar JSON da API externa: %v\n", err)
		return ""
	}

	var db, dbErr = openConnectionWithDataBase()
	if dbErr != nil {
		log.Println("[SERVER] Erro ao conectar com o banco de dados")
		return ""
	}

	err = saveToDatabase(db, &apiResp.USDBRL)
	if err != nil {
		log.Println("[SERVER] Erro ao salvar cotação no banco de dados")
		return ""
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("[SERVER] Erro ao obter instância do banco de dados: %v\n", err)
		return ""
	}
	defer sqlDB.Close()

	log.Println("[SERVER] Cotação obtida da API externa com sucesso")
	return apiResp.USDBRL.Bid
}

func openConnectionWithDataBase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("cotacoes.db"), &gorm.Config{})
	if err != nil {
		log.Printf("[DATABASE] Erro ao conectar com o banco: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&Usdbrl{})
	if err != nil {
		log.Printf("[DATABASE] Erro na migração: %v", err)
		return nil, err
	}

	log.Println("[DATABASE] Conexão com SQLite estabelecida com sucesso")
	return db, nil
}

func saveToDatabase(db *gorm.DB, usdbrl *Usdbrl) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), DatabaseTimeout)
	defer cancel()

	result := db.WithContext(dbCtx).Create(usdbrl)
	if result.Error != nil {
		log.Printf("[DATABASE] Erro ao salvar no banco: %v", result.Error)
		return result.Error
	}
	log.Printf("[DATABASE] Cotação salva no banco com sucesso. ID: %d", usdbrl.ID)
	return nil
}
