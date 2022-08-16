package health

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Handler is a HTTP Server Handler implementation
type Handler struct {
	CompositeChecker
}

// NewHandler returns a new Handler
func NewHandler() Handler {
	return Handler{}
}

// ServeHTTP returns a json encoded Health
// set the status to http.StatusServiceUnavailable if the check is down
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	const plain = "plain"

	switch format {
	case plain:
		w.Header().Add("Content-Type", "text/plain")
	default:
		w.Header().Add("Content-Type", "application/json")
	}

	health := h.CompositeChecker.Check(r.Context())

	if health.IsDown() {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	switch format {
	case plain:
		_ = plainText(w, health)
	default:
		_ = json.NewEncoder(w).Encode(health)
	}
}

func plainText(w io.Writer, health Health) error {
	var buf bytes.Buffer

	data := walkStatus("", health)
	if len(data) > 0 {
		for k, v := range data {
			buf.WriteString(fmt.Sprintf("%s %v\n", k, v))
		}
	}

	_, err := w.Write(buf.Bytes())

	return err
}

func walkStatus(prefix string, h Health) map[string]interface{} {
	const status = "status"

	r := make(map[string]interface{})

	statusKey := prefix + status
	if h.status == Up {
		r[statusKey] = 1
	} else {
		r[statusKey] = 0
	}

	for k, v := range h.info {
		switch value := v.(type) {
		case Health:
			res := walkStatus(k+"_", value)
			for k, v := range res {
				r[prefix+k] = v
			}
		case int:
			r[prefix+k] = value
		}
	}

	return r
}
