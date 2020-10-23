package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"khin-api/helper"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Welcome to the unofficial Khinsider-Ripper API!\n"+
		"")
	fmt.Println("Endpoint Hit: " + r.RequestURI)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/api/search/{term}", search)
	router.HandleFunc("/api/tracks/{code}", getTracks)
	router.HandleFunc("/api/get/{code}/{track}/{format}", getDirect)
	server := http.Server{
		Addr:    ":5100",
		Handler: router,
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
		},
	}

	fmt.Printf("Server listening on %s", server.Addr)
	if err := /*server.ListenAndServeTLS("certs/cert.pem", "certs/privkey.pem");*/ server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	println("Khinsider-Ripper API has started!\nPress CTRL+C to end.")
	handleRequests()
}

func search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !helper.Validate_key(r.Header.Get("API-Key")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	key := vars["term"]
	_ = json.NewEncoder(w).Encode(helper.SearchFor(key))
}

func getTracks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !helper.Validate_key(r.Header.Get("API-Key")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	key := vars["code"]
	_ = json.NewEncoder(w).Encode(helper.TrackGetter(key))
}

func getDirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !helper.Validate_key(r.Header.Get("API-Key")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	key := vars["code"]
	keyTr := vars["track"]
	keyFor := vars["format"]
	_ = json.NewEncoder(w).Encode(helper.DirectLinkGrab(key, keyTr, keyFor))
}
