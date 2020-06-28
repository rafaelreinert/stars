package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rafaelreinert/starts/pkg/planet"
	"github.com/rafaelreinert/starts/pkg/planet/retriever"
)

func (s *Server) handler() http.Handler {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "X-Session-Token"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	credentialsOk := handlers.AllowCredentials()
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "DELETE", "PUT", "OPTIONS"})

	r := mux.NewRouter()
	r.HandleFunc("/planets", s.getPlanetByNameHandler).Methods("GET").Queries("name", "")
	r.HandleFunc("/planets", s.listPlanetHandler).Methods("GET")
	r.HandleFunc("/planets", s.createPlanetHandler).Methods("POST")
	r.HandleFunc("/planets/{id}", s.getPlanetHandler).Methods("GET")
	r.HandleFunc("/planets/{id}", s.updatePlanetHandler).Methods("PUT")
	r.HandleFunc("/planets/{id}", s.deletePlanetHandler).Methods("DELETE")

	return handlers.CORS(headersOk, originsOk, methodsOk, credentialsOk)(r)
}

func (s *Server) listPlanetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	planets, err := retriever.RetriveAllPlanets(ctx, s.CountRetriever, s.PlanetRepository)
	if err != nil {
		log.Println("Error retriving all planets", err)
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := json.Marshal(planets)
	if err != nil {
		log.Println("Error Marshaling the result planets", err)
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Println("Error to write the response", err)
	}
}

func (s *Server) createPlanetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

	var newPlanet planet.Planet
	err := json.NewDecoder(r.Body).Decode(&newPlanet)
	if err != nil {
		log.Println("Error Decoding the planet", err)
		handleError(w, http.StatusBadRequest, "Planet JSON is Invalid")
		return
	}
	savedPlanet, err := s.PlanetRepository.Create(ctx, newPlanet)
	if err != nil {
		log.Println("Error Creating a planet", err)
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := json.Marshal(savedPlanet)
	if err != nil {
		log.Println("Error Marshaling a planet", err)
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		log.Println("Error to write the response", err)
	}
}

func (s *Server) getPlanetByNameHandler(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	planetName := r.URL.Query().Get("name")
	planet, err := retriever.RetrivePlanetByName(ctx, planetName, s.CountRetriever, s.PlanetRepository)
	if err != nil {
		log.Println("Error retriving the planet", err)
		handleError(w, http.StatusNotFound, err.Error())
		return
	}

	response, err := json.Marshal(planet)
	if err != nil {
		log.Println("Error Marshaling the result planet", err)
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Println("Error to write the response", err)
	}
}

func (s *Server) getPlanetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	vars := mux.Vars(r)
	planet, err := retriever.RetrivePlanet(ctx, vars["id"], s.CountRetriever, s.PlanetRepository)
	if err != nil {
		log.Println("Error retriving the planet", err)
		handleError(w, http.StatusNotFound, err.Error())
		return
	}

	response, err := json.Marshal(planet)
	if err != nil {
		log.Println("Error Marshaling the result planet", err)
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Println("Error to write the response", err)
	}
}

func (s *Server) updatePlanetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	vars := mux.Vars(r)
	var newPlanet planet.Planet
	err := json.NewDecoder(r.Body).Decode(&newPlanet)
	if err != nil {
		log.Println("Error Decoding the planet", err)
		handleError(w, http.StatusBadRequest, "Planet JSON is Invalid")
		return
	}
	newPlanet.ID = vars["id"]
	planet, err := s.PlanetRepository.Update(ctx, newPlanet)
	if err != nil {
		log.Println("Error updating the planet", err)
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := json.Marshal(planet)
	if err != nil {
		log.Println("Error Marshaling the result planet", err)
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Println("Error to write the response", err)
	}
}

func (s *Server) deletePlanetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	vars := mux.Vars(r)

	err := s.PlanetRepository.Delete(ctx, vars["id"])
	if err != nil {
		log.Println("Error deleting the planet", err)
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func handleError(w http.ResponseWriter, statusCode int, errorMensage string) {
	response, _ := json.Marshal(map[string]string{"error": errorMensage})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write(response)
	if err != nil {
		log.Println("Error to write the response", err)
	}
}
