package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Timeout excedido ao acessar API local")
		}
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Servidor retornou status %d\n", resp.StatusCode)
		return
	}

	var c Cotacao
	err = json.Unmarshal(body, &c)
	if err != nil {
		panic(err)
	}

	salvaCotacaoArquivo(c.Bid)
	fmt.Println(string(body))
}

func salvaCotacaoArquivo(bid string) {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	line := "Dólar: " + bid

	_, err = f.Write([]byte(line))
	if err != nil {
		panic(err)
	}
}
