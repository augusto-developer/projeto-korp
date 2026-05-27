package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Response representa o formato do JSON de retorno exigido.
type Response struct {
	Nome    string `json:"nome"`
	Horario string `json:"horario"`
}

var (
	startTime      = time.Now()
	registry       = prometheus.NewRegistry()
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "korp_http_requests_total",
			Help: "Quantidade total de requisições HTTP recebidas por este serviço.",
		},
		[]string{"path", "method", "status"},
	)
)

func init() {
	// Registrar métricas customizadas no registro customizado
	registry.MustRegister(requestCounter)

	// Registrar métrica de uptime dinâmica calculada sob demanda no scraping
	uptimeGauge := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "korp_service_uptime_seconds",
			Help: "Tempo de atividade do serviço http-server-projeto-korp em segundos.",
		},
		func() float64 {
			return time.Since(startTime).Seconds()
		},
	)
	registry.MustRegister(uptimeGauge)

	// Registrar coletores padrão do Go (Threads, Goroutines, GC, Memória)
	registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
}

func handleKorp(w http.ResponseWriter, r *http.Request) {
	status := "200"
	defer func() {
		requestCounter.WithLabelValues("/projeto-korp", r.Method, status).Inc()
	}()

	if r.Method != http.MethodGet {
		status = "405"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Método não permitido"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := Response{
		Nome:    "Projeto Korp",
		Horario: time.Now().UTC().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		status = "500"
		log.Printf("Erro ao serializar resposta: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	// Roteamento
	http.HandleFunc("/projeto-korp", handleKorp)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	log.Println("Iniciando servidor http-server-projeto-korp na porta 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Falha crítica no servidor HTTP: %v", err)
	}
}
