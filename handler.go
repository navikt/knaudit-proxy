package knaudit_proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Server struct {
	sender SendCloser
}

func NewServer(sender SendCloser) *Server {
	return &Server{
		sender: sender,
	}
}

func jsonResponse(w http.ResponseWriter, msg *Message) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(msg.Code)
	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		slog.Error("json encoding response", "error", err)
	}
}

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (h *Server) ReportHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("handling report request", "method", r.Method, "url", r.URL.String())

	data, err := io.ReadAll(r.Body)
	if err != nil {
		jsonResponse(w, &Message{
			Status:  "bad request",
			Message: fmt.Errorf("unable to read audit data: %w", err).Error(),
			Code:    http.StatusBadRequest,
		})

		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	err = h.sender.Send(string(data))
	if err != nil {
		jsonResponse(w, &Message{
			Status:  "bad request",
			Message: fmt.Errorf("storing audit data: %w", err).Error(),
			Code:    http.StatusInternalServerError,
		})

		return
	}

	jsonResponse(w, &Message{
		Status:  "ok",
		Message: "audit data stored",
		Code:    http.StatusOK,
	})
}

func (h *Server) StatusHandler(w http.ResponseWriter, _ *http.Request) {
	jsonResponse(w, &Message{
		Status:  "ok",
		Message: "knaudit-proxy is running",
		Code:    http.StatusOK,
	})
}
