# Gophernoon - Building A Link Shortner Service, Solution Part 1

## Listening For HTTP Requests

``` go
package main

import (
    "log"
    "net/http"
)

func main() {
    log.Print("starting http server on port 80...")
    log.Fatal(http.ListenAndServe(":80", nil))
}
```

The above code will log "starting http server on port 80..." and then attempt
to start listening on port `80` for HTTP traffic. The `http.ListenAndServe`
function is blocking code and will only stop blocking once an error has been
detected and returned, once this happens the application will log it to the
`stderr`. It's worth noting that the `log` package always outputs to `stderr`.

## Adding HTTP Handlers

``` go
package main

import (
	"log"
	"net/http"
)

type router struct {
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/create":
		handleCreateURLAliasPair(w, req)
	default:
		handleRedirect(w, req)
	}
}

func handleCreateURLAliasPair(w http.ResponseWriter, req *http.Request) {
}

func handleRedirect(w http.ResponseWriter, req *http.Request) {
}

func main() {
	r := new(router)

	log.Print("starting http server on port 80...")
	log.Fatal(http.ListenAndServe(":80", r))
}

```

Here we add a router struct and attach the `ServeHTTP` method to it. This 
method and signature is required to use it as a `http.Handler` in the
`http.ListenAndServe` function.

Inside of the `ServeHTTP` function we build some basic routing based on the
path. If the path is `/create` then we handle the creation of a new URL alias 
pair, otherwise we will try to redirect.

## Creating New URL:Alias Pairs

``` go
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
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
}

func main() {
	r := new(router)

	log.Print("starting http server on port 80...")
	log.Fatal(http.ListenAndServe(":80", r))
}

```

Here we have added a `map[string]string` to the router struct, to store our
`url:alias` pairs, as well as add a `newRouter` function to handle the 
initialisation of the map. We also fleshed out the `handleCreateURLAliasPair`
function to actually do some work.

First it checks that request is a `POST` request, if not a HTTP status 
`405 Method Not Allowed` code is written to the response and the function
hits an early return statement. Secondly we parse the form data to make it 
accessible to our code, if an error is returned we respond with a http status
code of `500 Internal Server Error`. Thirdly we validate that the submitted 
form has `alias` in it and that it isn't empty, if it is we respond with a 
`400 Bad Request` status code. Then we go on to parse the `url` parameter from
the query string, again responding with a `400 Bad Request` HTTP status code if
validation fails. Finally we assign an item in the map where the `alias` is the
key and the `url` is the value. We'll use this later when handling the 
redirect.

Note in the `main` function we have implemented the `newRouter` constructor to
ensure the map is intialised. If we write to an uninitialised map we will get a
segmentation fault error.

### Redirecting

``` go
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
    alias := strings.TrimLeft(req.URL.Path, "/")
	redirectURL, ok := r.urlAliasPairs[alias]
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
```

Finally we need to handle redirecting if the alias is provided. for this we get
the path element from the request and trim off the trailing slash, using the
`strings.TrimLeft` function. Next we attempt to get the url value from the 
stored map, in the event that the `alias` key can not be found in the map we
respond with a http status code `404 Not Found` and hit an early return. 
Finally we redirect the http request to the `url` value from the map.