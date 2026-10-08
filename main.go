package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	artisturl    = "https://groupietrackers.herokuapp.com/api/artists"
	relationsurl = "https://groupietrackers.herokuapp.com/api/relation"
)

var (
	artists    []Artist
	tmpl       *template.Template
	artistTmpl *template.Template
	relations  map[int]Relation
)

func fetchArtists() ([]Artist, error) {
	resp, err := http.Get(artisturl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}
func fetchRelations() ([]Relation, error) {
	resp, err := http.Get(relationsurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var RELATION RelationsIndex
	err = json.NewDecoder(resp.Body).Decode(&RELATION)
	if err != nil {
		return nil, err
	}
	return RELATION.Index, nil
}
func artistdetails(a *Artist) ArtistDetail {
	detail := ArtistDetail{Artist: a}

	if relation, ok := relations[a.ID]; ok {
		detail.Concerts = relation.DatesLocations
	}

	return detail
}

func implementfetch() error {
	var (
		artistsData   []Artist
		relationsData []Relation
		artistErr     error
		relationErr   error
	)
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go func() {
		defer waitGroup.Done()
		artistsData, artistErr = fetchArtists()
	}()
	go func() {
		defer waitGroup.Done()
		relationsData, relationErr = fetchRelations()
	}()
	waitGroup.Wait()
	if artistErr != nil {
		return artistErr
	}
	if relationErr != nil {
		return relationErr
	}

	artists = artistsData

	relations = make(map[int]Relation)

	for _, relation := range relationsData {
		formattedLocations := make(map[string][]string)

		for location, dates := range relation.DatesLocations {
			location = strings.ReplaceAll(location, "_", " ")
			location = strings.ReplaceAll(location, "-", " ")

			formattedLocations[location] = dates
		}

		relation.DatesLocations = formattedLocations
		relations[relation.ID] = relation
	}

	return nil
}

func indexhandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if err := tmpl.Execute(w, artists); err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	}

}
func findArtistByID(id int) *Artist {
	for i := range artists {
		if artists[i].ID == id {
			return &artists[i]
		}
	}
	return nil
}
func artistHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/artist" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	artistid := r.URL.Query().Get("id")
	Artistidint, err := strconv.Atoi(artistid)
	if err != nil {
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}
	artist := findArtistByID(Artistidint)
	if artist == nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	detail := artistdetails(artist)
	if err := artistTmpl.Execute(w, detail); err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	}

}

func main() {
	var err error
	err = implementfetch()
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	artistTmpl, err = template.ParseFiles("templates/artist.html")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/styling.css", http.FileServer(http.Dir("templates")))
	http.HandleFunc("/", indexhandler)
	http.HandleFunc("/artist", artistHandler)

	log.Println("server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
