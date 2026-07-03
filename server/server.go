package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type Cotacao struct {
	Usdbrl struct {
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
	} `json:"USDBRL"`
}

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

var db *sql.DB

func main() {
	InicializaBanco()

	http.HandleFunc("/cotacao", BuscaCotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func InicializaBanco() {
	var err error

	db, err = sql.Open("sqlite", "../cotacao.db")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cotacao (id INTEGER PRIMARY KEY, bid TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		panic(err)
	}
}

func BuscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cotacao, err := BuscaCotacao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = insertCotacao(db, cotacao)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := CotacaoResponse{
		Bid: cotacao.Usdbrl.Bid,
	}

	json.NewEncoder(w).Encode(response)
}

func BuscaCotacao() (*Cotacao, error) {
	ctxApi, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxApi, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctxApi.Err() == context.DeadlineExceeded {
			log.Println("Timeout excedido ao acessar API externa")
		}

		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var c Cotacao
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func insertCotacao(db *sql.DB, cotacao *Cotacao) error {
	ctxDb, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	stmt, err := db.Prepare("insert into cotacao (bid) values ($1)")

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctxDb, cotacao.Usdbrl.Bid)
	if err != nil {
		if ctxDb.Err() == context.DeadlineExceeded {
			log.Println("Timeout excedido ao inserir no banco de dados")
		}
		return err
	}
	return nil
}
