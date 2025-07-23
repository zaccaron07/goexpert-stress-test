package usecase

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteSuccess(t *testing.T) {
	var counter int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqNum := atomic.AddInt64(&counter, 1)

		if reqNum%2 == 0 {
			w.WriteHeader(http.StatusTooManyRequests)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()

	totalRequests := 100
	cfg := Config{
		URL:         srv.URL,
		Requests:    totalRequests,
		Concurrency: 3,
	}

	tester := NewStressTester(srv.Client())

	report, err := tester.Execute(context.Background(), cfg)
	assert.NoError(t, err)
	assert.Equal(t, cfg.Requests, report.TotalRequests)

	expectedOK := totalRequests / 2
	expectedTooManyRequests := totalRequests - expectedOK

	assert.Equal(t, expectedOK, report.StatusCodes[http.StatusOK])
	assert.Equal(t, expectedTooManyRequests, report.StatusCodes[http.StatusTooManyRequests])
	assert.Equal(t, cfg.URL, report.URL)
	assert.NotZero(t, report.TotalTime)
}

func TestExecuteNetworkError(t *testing.T) {
	cfg := Config{
		URL:         "",
		Requests:    5,
		Concurrency: 2,
	}

	tester := NewStressTester(&http.Client{})

	report, err := tester.Execute(context.Background(), cfg)
	assert.NoError(t, err)
	assert.Equal(t, cfg.Requests, report.TotalRequests)

	assert.Equal(t, cfg.Requests, report.StatusCodes[0])
	assert.Zero(t, report.StatusCodes[http.StatusOK])
	assert.Equal(t, cfg.URL, report.URL)
	assert.NotZero(t, report.TotalTime)
}
