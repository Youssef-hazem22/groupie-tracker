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
	artists   []Artist
	tmpl      *template.Template
	relations map[int]Relation
	client    = &http.Client{Timeout: 10 * time.Second}
)

func fetchArtists() ([]Artist, error) {
	resp, err := client.Get(artisturl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected code : %s", resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}
func fetchRelations() ([]Relation, error) {
	resp, err := client.Get(relationsurl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected code : %s", resp.Status)
	}
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
	q := r.URL.Query().Get("q")

	data := IndexData{
		Artists: searchArtists(q),
		Query:   q,
	}

	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	}
}

func searchArtists(query string) []Artist {
	if len(query) > 100 {
		query = query[:100]
	}
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return artists
	}

	var results []Artist
	for _, a := range artists {
		if artistMatches(a, query) {
			results = append(results, a)
		}
	}
	return results
}

func artistMatches(a Artist, q string) bool {
	if strings.Contains(strings.ToLower(a.Name), q) ||
		strings.Contains(strings.ToLower(a.FirstAlbum), q) ||
		strings.Contains(strconv.Itoa(a.CreationDate), q) {
		return true
	}

	for _, member := range a.Members {
		if strings.Contains(strings.ToLower(member), q) {
			return true
		}
	}

	for location := range relations[a.ID].DatesLocations {
		if strings.Contains(strings.ToLower(location), q) {
			return true
		}
	}

	return false
}
func artistHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := tmpl.ExecuteTemplate(w, "artist.html", detail); err != nil {
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
