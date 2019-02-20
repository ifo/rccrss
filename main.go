package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/feeds"
)

func main() {
	// Setup
	// Get the server port.
	portStr := os.Getenv("PORT")
	portDefault := 3000
	var err error
	if portStr != "" {
		portDefault, err = strconv.Atoi(portStr)
	}
	if err != nil {
		log.Fatal(err)
	}
	port := flag.Int("port", portDefault, "Port to run the server on")

	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/", indexHandler)
	r.Get("/feed.xml", feedHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index page"))
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	feed := &feeds.Feed{
		Title: "Recurse Center Code Review Simple Syndicator",
		Link:  &feeds.Link{Href: "https://rccrss.recurse.com"},
	}

	atom, err := feed.ToAtom()
	if err != nil {
		http.Error(w, "We couldn't make a feed 😭", 500)
	}

	w.Write([]byte(atom))
}
