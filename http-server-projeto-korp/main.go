package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

const serviceName = "http-server-projeto-korp"

var (
	requestsTotal      uint64
	serviceAvailable   int64 = 1
)

type projetoKorpResponse struct {
	Nome    string `json:"nome"`
	Horario string `json:"horario"`
}

func projetoKorpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		atomic.AddUint64(&requestsTotal, 1)
		return
	}

	resp := projetoKorpResponse{
		Nome:    "Projeto Korp",
		Horario: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("erro ao serializar resposta: %v", err)
		atomic.AddUint64(&requestsTotal, 1)
		return
	}
	atomic.AddUint64(&requestsTotal, 1)
}

func metricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	total := atomic.LoadUint64(&requestsTotal)
	available := atomic.LoadInt64(&serviceAvailable)

	fmt.Fprintf(w, "# HELP projeto_korp_service_available Disponibilidade do serviço (1=disponível, 0=indisponível)\n")
	fmt.Fprintf(w, "# TYPE projeto_korp_service_available gauge\n")
	fmt.Fprintf(w, "projeto_korp_service_available %d\n", available)

	fmt.Fprintf(w, "# HELP projeto_korp_requests_total Total de requisições ao endpoint /projeto-korp\n")
	fmt.Fprintf(w, "# TYPE projeto_korp_requests_total counter\n")
	fmt.Fprintf(w, "projeto_korp_requests_total %d\n", total)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/projeto-korp", projetoKorpHandler)
	mux.HandleFunc("/metrics", metricsHandler)

	addr := ":8080"
	log.Printf("%s escutando em %s", serviceName, addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		atomic.StoreInt64(&serviceAvailable, 0)
		log.Fatalf("servidor encerrado: %v", err)
	}
}
