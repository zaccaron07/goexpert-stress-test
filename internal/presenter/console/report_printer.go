package console

import (
	"fmt"
	"sort"

	"goexpert-stress-test/internal/entity"
)

type ReportPrinter struct{}

func NewReportPrinter() *ReportPrinter {
	return &ReportPrinter{}
}

func (p *ReportPrinter) Print(report entity.Report) error {
	fmt.Println("===== Relatório =====")
	fmt.Printf("URL: %s\n", report.URL)
	fmt.Printf("Tempo total: %v\n", report.TotalTime)
	fmt.Printf("Total de requests: %d\n", report.TotalRequests)
	fmt.Printf("Requests com HTTP 200: %d\n", report.StatusCodes[200])
	fmt.Println("Distribuição de status codes:")

	codes := make([]int, 0, len(report.StatusCodes))
	for code := range report.StatusCodes {
		codes = append(codes, code)
	}
	sort.Ints(codes)

	for _, code := range codes {
		fmt.Printf("  %d -> %d\n", code, report.StatusCodes[code])
	}
	return nil
}
