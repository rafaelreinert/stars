package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rafaelreinert/stars/pkg/config"
	"github.com/rafaelreinert/stars/pkg/planet/repository"
	"github.com/rafaelreinert/stars/pkg/planet/retriever"
)

// Server is the struct which  initializes and control the HTTP server and all API handles
type Server struct {
	PlanetRepository repository.PlanetRepository
	CountRetriever   retriever.PlanetAppearancesOnMoviesCounter
	Cfg              config.Config
}

// ListenAndServe starts an HTTP server with API handler loaded, it uses PORT env variable or port 8080
func (s *Server) ListenAndServe() {
	log.Println("Listening ", fmt.Sprintf(":%d", s.Cfg.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Cfg.Port), s.handler()))
}
