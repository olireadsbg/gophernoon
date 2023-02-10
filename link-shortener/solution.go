package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type router struct {
	urlAliasPairs map[string]string
}

func newRouter() *router {
	return &router{
		urlAliasPairs: make(map[string]string),
	}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/create":
		r.handleCreateURLAliasPair(w, req)
	default:
		r.handleRedirect(w, req)
	}
}

func (r *router) handleCreateURLAliasPair(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := req.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if req.Form.Get("alias") == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid alias: " + req.Form.Get("alias")))
		return
	}

	if _, err := url.Parse(req.Form.Get("url")); err != nil {

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid url: " + req.URL.Query().Get("url")))
		return
	}

	r.urlAliasPairs[req.Form.Get("alias")] = req.Form.Get("url")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(
		fmt.Sprintf("/%s redirects to %s",
			req.Form.Get("alias"),
			req.Form.Get("url"),
		),
	))
}

func (r *router) handleRedirect(w http.ResponseWriter, req *http.Request) {
	redirectURL, ok := r.urlAliasPairs[strings.TrimLeft(req.URL.Path, "/")]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("invalid alias"))
		return
	}

	http.Redirect(w, req, redirectURL, http.StatusMovedPermanently)
}

func main() {
	r := newRouter()

	log.Print("starting http server on port 80...")
	log.Fatal(http.ListenAndServe(":80", r))
}
