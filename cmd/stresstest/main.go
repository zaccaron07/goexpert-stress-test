package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"goexpert-stress-test/internal/presenter/console"
	"goexpert-stress-test/internal/usecase"
)

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 1, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")

	flag.Parse()

	if *url == "" {
		log.Fatal("o parâmetro --url é obrigatório")
	}
	if *requests <= 0 {
		log.Fatal("--requests deve ser maior que 0")
	}
	if *concurrency <= 0 {
		log.Fatal("--concurrency deve ser maior que 0")
	}

	if *concurrency > *requests {
		*concurrency = *requests
	}

	client := &http.Client{}

	stressTester := usecase.NewStressTester(client)

	report, err := stressTester.Execute(context.Background(), usecase.Config{
		URL:         *url,
		Requests:    *requests,
		Concurrency: *concurrency,
	})
	if err != nil {
		log.Fatal(err)
	}

	printer := console.NewReportPrinter()
	if err := printer.Print(report); err != nil {
		log.Fatal(err)
	}
}
