package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	artisturl    = "https://groupietrackers.herokuapp.com/api/artists"
	relationsurl = "https://groupietrackers.herokuapp.com/api/relation"
	locationurl  = "https://groupietrackers.01edu.ai/api/locations"
	datesurl     = "https://groupietrackers.01edu.ai/api/dates"
)

var (
	artists    []Artist
	tmpl       *template.Template
	artistTmpl *template.Template
	relations  map[int]Relation
	client     = &http.Client{Timeout: 10 * time.Second}
)

func fetchArtists() ([]Artist, error) {
	resp, err := client.Get(artisturl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected code : %s", resp.Status)
	}
	defer resp.Body.Close()
	var result []Artist
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func fetchRelations() ([]Relation, error) {
	resp, err := client.Get(relationsurl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected code : %s", resp.Status)
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
	if r.Method != http.MethodGet {
		http.Error(w, "405 methodnotallowed", http.StatusMethodNotAllowed)
	}

	if err := tmpl.Execute(w, artists); err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	}

}
func artistHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/artist" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "405 methodnotallowed", http.StatusMethodNotAllowed)
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
func findArtistByID(id int) *Artist {
	for i := range artists {
		if artists[i].ID == id {
			return &artists[i]
		}
	}
	return nil
}
