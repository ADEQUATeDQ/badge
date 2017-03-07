package main

import (
	"log"
	"net/http"
	"os"

	"github.com/the42/badge"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func renderbadge(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Maybe we need to know where we come from
		// data.gv.at vs. opendataportal.at
		// and render only for those hosts
		// Maybe better though to operate that service behind Nginx or like
		// and perform the port mapping / filtering there

		// host := r.Host

		id := r.URL.Path[1:]

		// for now, do not serve a badge if there is no indication for what ID to retrieve information
		if len(id) > 0 {
			// do here all the wizardry to retrieve information about this data set
			w.Header().Set("Content-Type", "image/svg+xml")
			badge.Render("ADEQUATe", id, "#5272B4", w)
		}
	}
}

func main() {
	http.HandleFunc("/", renderbadge)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Fatal(http.ListenAndServe(":"+port, Log(http.DefaultServeMux)))

}
