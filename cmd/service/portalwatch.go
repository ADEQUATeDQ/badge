package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/the42/badge"
)

var AdequateGit = "https://data.adequate.at/api/v4/projects/"

func portalwatch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Maybe we need to know where we come from
		// data.gv.at vs. opendataportal.at
		// and render only for those hosts
		//
		// It may be better though to operate that service behind Nginx or like
		// and perform the port mapping / filtering there

		// host := r.Host
		if err := r.ParseForm(); err != nil {
			s := fmt.Sprintf("Error parsing form: %s", err)
			log.Printf(s)
			http.Error(w, s, http.StatusInternalServerError)
			return
		}

		// expecting the url to be of the form <basepath>portal_id/dataset_id
		// therefore split at '/'; first part is the portal_id, second part is the dataset_id
		parameters := strings.Split(r.URL.Path[len(basepath):], "/")

		if len(parameters) < 2 {
			s := fmt.Sprintf("Not enough parameters in call to %s badge service", basepath)
			log.Printf(s)
			http.Error(w, s, http.StatusInternalServerError)
			return
		}
		portal := parameters[0]
		id := parameters[1]

		// do not serve a badge if there is no indication for what ID or portal to retrieve information
		if len(id) > 0 && len(portal) > 0 {
			p := AdequateGit + portal + "%2F" + id

			// perform the Portalwatch quality check call. For now, the result interpretation is very easy.
			// If a HTTP-code of 200 is returned, we assume a quality check has been performed and render a badge
			resp, err := http.Head(p)

			if err != nil {
				s := err.Error()
				log.Printf(s)
				http.Error(w, s, http.StatusInternalServerError)
				return
			}

			switch resp.StatusCode {
			case http.StatusOK:
				w.Header().Set("Content-Type", "image/svg+xml")
				badge.Render("ADEQUATe", "Checked "+"\xE2\x9C\x94", "brightgreen", w)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

const basepath = "/adequate/portalwatch/"

func init() {
	http.HandleFunc(basepath, portalwatch)
}
