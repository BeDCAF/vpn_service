package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(httpHandler *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandler,
	}
}

func (s *HTTPServer) StartServer(addr, port string) error {
	router := mux.NewRouter()

	router.Path("/clients").Methods("POST").HandlerFunc(s.httpHandlers.HandleAddUser)
	router.Path("/clients").Methods("DELETE").HandlerFunc(s.httpHandlers.HandleDeleteUser)

	if err := http.ListenAndServe(addr+":"+port, router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}

	return nil
}
