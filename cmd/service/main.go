package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func info(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintln(w, "Happily serving")
		fmt.Fprintln(w)

		httpMux := reflect.ValueOf(http.DefaultServeMux).Elem()
		finList := httpMux.FieldByName("m")
		fmt.Fprintln(w, finList)
	}
}

func main() {
	http.HandleFunc("/", info)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Fatal(http.ListenAndServe(":"+port, Log(http.DefaultServeMux)))

}
