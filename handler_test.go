package health

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewHandler(t *testing.T) {
	// How can I test a function that returns a Struct ?
	// A better Idea? Please tell me!
	h := NewHandler()
	handler := &h

	if handler == nil {
		t.Error("&NewHandler() == nil, wants !nil")
	}
}

func Test_Handler_ServeHTTP_Down(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.AddChecker("DownChecker", downTestChecker{})

	h.ServeHTTP(w, r)

	jsonbytes, _ := ioutil.ReadAll(w.Body)
	jsonstring := strings.TrimSpace(string(jsonbytes))

	wants := `{"DownChecker":{"status":"DOWN"},"status":"DOWN"}`

	if jsonstring != wants {
		t.Errorf("jsonReturned == %s, wants %s", jsonstring, wants)
	}

	contentType := w.Header().Get("Content-Type")
	wants = "application/json"

	if contentType != wants {
		t.Errorf("type == %s, wants %s", contentType, wants)
	}

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusServiceUnavailable)
	}
}

func Test_Handler_ServeHTTP_Up(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.AddChecker("UpChecker", upTestChecker{})
	h.AddInfo("custom", "info")

	h.ServeHTTP(w, r)
	jsonbytes, _ := ioutil.ReadAll(w.Body)
	jsonstring := strings.TrimSpace(string(jsonbytes))

	wants := `{"UpChecker":{"status":"UP"},"custom":"info","status":"UP"}`

	if jsonstring != wants {
		t.Errorf("jsonstring == %s, wants %s", jsonstring, wants)
	}

	contentType := w.Header().Get("Content-Type")
	wants = "application/json"

	if contentType != wants {
		t.Errorf("type == %s, wants %s", contentType, wants)
	}

	if w.Code != http.StatusOK {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusOK)
	}
}

type postgresComplexChecker struct{}

func (c postgresComplexChecker) Check(ctx context.Context) Health {
	pg := NewHealth()
	pg.Up()

	stats := map[string]interface{}{
		"idle":                 0,               // The number of idle connections.
		"in_use":               2,               // The number of connections currently in use.
		"max_idle_closed":      0,               // The total number of connections closed due to SetMaxIdleConns.
		"max_idle_time_closed": 3 * time.Second, // The total number of connections closed due to SetConnMaxIdleTime.
		"max_life_time_closed": 2 * time.Second, // The total number of connections closed due to SetConnMaxLifetime.
		"max_open_connections": 10,              // Maximum number of open connections to the database.
		"open_connections":     2,               // Pool Status
		"wait_count":           5,               // Counters
		"wait_duration":        time.Second,     // The total time blocked waiting for a new connection.
	}

	sub := NewHealth()
	for k, v := range stats {
		sub.AddInfo(k, v)
	}

	pg.AddInfo("stats", sub)

	return pg
}

func Test_Handler_ServeHTTP_Up_Plain(t *testing.T) {
	r, _ := http.NewRequest("GET", "/?format=plain", nil)
	w := httptest.NewRecorder()

	h := Handler{}
	h.AddChecker("postgres", postgresComplexChecker{})

	h.ServeHTTP(w, r)

	jsonbytes, _ := ioutil.ReadAll(w.Body)
	jsonstring := strings.TrimSpace(string(jsonbytes))
	substings := strings.Split(jsonstring, "\n")

	ss := []string{
		"status 1",
		"postgres_stats_wait_count 5",
		"postgres_status 1",
		"postgres_stats_idle 0",
		"postgres_stats_max_open_connections 10",
		"postgres_stats_in_use 2",
		"postgres_stats_status 0",
		"postgres_stats_open_connections 2",
		"postgres_stats_max_idle_closed 0",
	}

	assert.ElementsMatch(t, ss, substings)

	contentType := w.Header().Get("Content-Type")
	wants := "text/plain"

	if contentType != wants {
		t.Errorf("type == %s, wants %s", contentType, wants)
	}

	if w.Code != http.StatusOK {
		t.Errorf("w.Code == %d, wants %d", w.Code, http.StatusServiceUnavailable)
	}
}
