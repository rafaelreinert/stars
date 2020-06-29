package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rafaelreinert/stars/pkg/config"
	"github.com/rafaelreinert/stars/pkg/planet/repository"
	"github.com/rafaelreinert/stars/pkg/planet/retriever"
)

type Server struct {
	PlanetRepository repository.PlanetRepository
	CountRetriever   retriever.PlanetAppearancesOnMoviesCounter
	Cfg              config.Config
}

func (s *Server) ListenAndServe() {
	log.Println("Listening ", fmt.Sprintf(":%d", s.Cfg.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Cfg.Port), s.handler()))
}
