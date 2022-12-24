package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Movie struct {
	Name     string    `json:"name"`
	Id       string    `json:"id"`
	Code     string    `json:"code"`
	Director *Director `json:"director"`
}
type Director struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

var Movies []Movie

func main() {
	r := mux.NewRouter()
	Movies = append(Movies, Movie{Id: "1", Code: "123142", Name: "Matrix", Director: &Director{"Nika", "Chighladze"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// metrices

	r.Handle("/metrices",promhttp.Handler())

	fmt.Printf("Server is starting:\n")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatalf("Error in listening %v", err)
	}

}
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestedId := mux.Vars(r)
	for i, Movie := range Movies {
		if Movie.Id == requestedId["id"] {
			Movies = append(Movies[:i], Movies[i+1:]...)
			json.NewEncoder(w).Encode(Movies)
			return
		}
	}
	fmt.Fprintf(w, "Movie with requested ID not exists in database")
}
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requstedId := mux.Vars(r)
	for i, _ := range Movies {
		if Movies[i].Id == requstedId["id"] {
			return
		}
	}
	fmt.Fprintf(w, "Movie with requested ID not exists in database")
}
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		log.Fatalf("Error in Decoding movie json : %v", err)
	}
	movie.Id = strconv.Itoa(rand.Intn(99999))
	Movies = append(Movies, movie)
	err = json.NewEncoder(w).Encode(Movies)
	if err != nil {
		log.Fatalf("Error in encoding movie json : %v", err)
	}
}
func updateMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	requestedId := mux.Vars(r)
	for i, _ := range Movies {
		if Movies[i].Id == requestedId["id"] {
			var movie Movie
			movie.Id = Movies[i].Id
			Movies = append(Movies[:i], Movies[i+1:]...)
			_ = json.NewDecoder(r.Body).Decode(&movie)
			Movies = append(Movies, movie)
		}
	}
}
