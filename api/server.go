package api

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"xcodeengine/executor"
	"xcodeengine/model"
	"xcodeengine/problems"
	"xcodeengine/service"

	compilergrpc "github.com/lijuuu/GlobalProtoXcode/Compiler"
)

// ExecuteRequest captures payloads from the UI.
type ExecuteRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	Mode     string `json:"mode"` // "problem" for relaxed limits, anything else is standard
	Input    string `json:"input"`
}

// ExecuteResponse mirrors the compiler response with HTTP friendly error reporting.
type ExecuteResponse struct {
	Output        string `json:"output"`
	Error         string `json:"error,omitempty"`
	StatusMessage string `json:"status_message"`
	Success       bool   `json:"success"`
	ExecutionTime string `json:"execution_time,omitempty"`
}

// StartServer boots a simple HTTP server that exposes the execution API and serves the static UI.
func StartServer(addr string, workerPool *executor.WorkerPool) {
	mux := http.NewServeMux()
	compilerService := service.NewCompilerService(workerPool)

	mux.HandleFunc("/api/execute", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ExecuteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		req.Language = strings.TrimSpace(req.Language)
		req.Code = strings.TrimSpace(req.Code)

		if req.Language == "" || req.Code == "" {
			http.Error(w, "code and language are required", http.StatusBadRequest)
			return
		}

		var (
			resp *compilergrpc.CompileResponse
			err  error
		)

		if strings.EqualFold(req.Mode, "problem") {
			resp, err = compilerService.ExecuteProblemCode(req.Code, req.Language)
		} else {
			encoded := base64.StdEncoding.EncodeToString([]byte(req.Code))
			resp, err = compilerService.Compile(encoded, req.Language, req.Input)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, ExecuteResponse{
			Output:        resp.Output,
			Error:         resp.Error,
			StatusMessage: resp.StatusMessage,
			Success:       resp.Success,
			ExecutionTime: resp.ExecutionTime,
		})
	})

	mux.HandleFunc("/api/problems", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, http.StatusOK, problems.ListProblems())
	})

	mux.HandleFunc("/api/problems/submit", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req model.ProblemSubmissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if req.ProblemID == "" {
			http.Error(w, "problem_id is required", http.StatusBadRequest)
			return
		}

		resp, err := compilerService.JudgeProblem(req.Code, req.Language, req.ProblemID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	})

	fileServer := http.FileServer(http.Dir("web"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "web/index.html")
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	log.Printf("HTTP UI server listening on %s", addr)
	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Printf("http server error: %v", err)
		}
	}()
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	setCORSHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
