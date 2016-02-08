package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	addr       = flag.String("addr", ":8080", "address to listen on")
	emailsFile = flag.String("emails-file", "emails.txt", "path to email file")
)

func writeEmails(emails <-chan string) {
	f, err := os.OpenFile(*emailsFile, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	csvw := csv.NewWriter(f)
	for e := range emails {
		now := time.Now()
		err = csvw.Write([]string{e, now.Format(time.RFC3339)})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	emails := make(chan string)
	go writeEmails(emails)
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		dec := json.NewDecoder(req.Body)
		e := make(map[string]string)
		if err := dec.Decode(e); err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if email, ok := e["email"]; ok {
			emails <- email
		} else {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		rw.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}
