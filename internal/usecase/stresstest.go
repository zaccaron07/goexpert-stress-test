package usecase

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"goexpert-stress-test/internal/entity"
)

type Config struct {
	URL         string
	Requests    int
	Concurrency int
}

type StressTester interface {
	Execute(ctx context.Context, cfg Config) (entity.Report, error)
}

type stressTester struct {
	client *http.Client
}

func NewStressTester(client *http.Client) StressTester {
	return &stressTester{client: client}
}

func (s *stressTester) Execute(ctx context.Context, cfg Config) (entity.Report, error) {
	start := time.Now()

	tasks := make(chan struct{})
	statusCh := make(chan int)

	statusCounts, aggWg := startAggregator(statusCh)
	workersWg := s.startWorkers(ctx, cfg.URL, cfg.Concurrency, tasks, statusCh)

	go func() {
		for i := 0; i < cfg.Requests; i++ {
			tasks <- struct{}{}
		}
		close(tasks)
	}()

	workersWg.Wait()
	close(statusCh)
	aggWg.Wait()

	report := entity.Report{
		URL:           cfg.URL,
		TotalTime:     time.Since(start),
		TotalRequests: cfg.Requests,
		StatusCodes:   statusCounts,
	}

	return report, nil
}

func (s *stressTester) startWorkers(ctx context.Context, url string, concurrency int, tasks <-chan struct{}, statusCh chan<- int) *sync.WaitGroup {
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.worker(ctx, url, tasks, statusCh)
		}()
	}
	return &wg
}

func (s *stressTester) worker(ctx context.Context, url string, tasks <-chan struct{}, statusCh chan<- int) {
	for range tasks {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			statusCh <- 0
			continue
		}

		resp, err := s.client.Do(req)
		if err != nil {
			log.Printf("Erro ao fazer request: %v", err)
			statusCh <- 0
			continue
		}

		statusCh <- resp.StatusCode
		resp.Body.Close()
	}
}

func startAggregator(statusCh <-chan int) (map[int]int, *sync.WaitGroup) {
	counts := make(map[int]int)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for code := range statusCh {
			counts[code]++
		}
	}()
	return counts, &wg
}
