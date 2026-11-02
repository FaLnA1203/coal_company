package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	HTTPHandler *HTTPHandler
	HTTPServer  *http.Server
}

func NewServer(handlers *HTTPHandler) *Server {
	return &Server{
		HTTPHandler: handlers,
	}
}

func (server *Server) Run() error {
	router := chi.NewRouter()

	router.Get("/miners/price", server.HTTPHandler.PriceMinersHandler)
	router.Get("/equipment/price", server.HTTPHandler.PriceEquipmentHandler)
	router.Post("/miners/buy", server.HTTPHandler.BuyMinerHandler)
	router.Post("/equpment/buy", server.HTTPHandler.BuyEquipmentHandler)
	router.Get("/miners", server.HTTPHandler.HandlerGetMiners)
	router.Get("/equipment", server.HTTPHandler.HandlerGetEquipment)
	router.Get("/company", server.HTTPHandler.HandlerInfoCompany)
	router.Patch("/company/close", server.HTTPHandler.HandlerCloseCompany)

	s := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.HTTPHandler.SetCloseServerFunc(s.Close)

	err := s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err

}
