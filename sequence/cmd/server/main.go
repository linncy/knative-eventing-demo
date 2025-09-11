package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	step := getenv("STEP_NAME", "step")
	addr := ":8080"

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	handler := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		b, _ := io.ReadAll(r.Body)

		var body any
		if len(b) > 0 && (b[0] == '{' || b[0] == '[') {
			_ = json.Unmarshal(b, &body)
		} else if len(b) > 0 {
			body = string(b)
		}

		// Collect CloudEvents binary mode headers (optional, for debug echo)
		ce := map[string]string{}
		for k, v := range r.Header {
			if len(v) > 0 && strings.HasPrefix(strings.ToLower(k), "ce-") {
				ce[k] = v[0]
			}
		}

		pod := getenv("HOSTNAME", "")
		log.Printf("[%-5s] %s %s len(body)=%d ce=%d", step, r.Method, r.URL.Path, len(b), len(ce))

		// Respond in CloudEvents structured mode
		resp := map[string]any{
			"specversion": "1.0",
			"id":          fmt.Sprintf("%s-%d", step, time.Now().UnixNano()),
			"source":      fmt.Sprintf("kn-seq/%s", step),
			"type":        "dev.demo.processed",
			"time":        time.Now().Format(time.RFC3339Nano),
			"data": map[string]any{
				"step":       step,
				"path":       r.URL.Path,
				"pod":        pod,
				"ce_headers": ce,
				"body_echo":  body,
				"message":    fmt.Sprintf("handled by %s", step),
			},
		}
		w.Header().Set("Content-Type", "application/cloudevents+json")
		_ = json.NewEncoder(w).Encode(&resp)
	}

	mux.HandleFunc("/", handler)
	mux.HandleFunc("/process", handler)

	log.Printf("starting %s on %s ...", step, addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
