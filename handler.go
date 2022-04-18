package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	UpdateGoLinkHTML = `
		<h1>Update Go Links</h1>
		<form action="/_api/update" method="POST">
			<label for="shortLink">Short Link: go/</label>
  			<input type="text" id="shortLink" name="shortLink"><br><br>
  			<label for="originalLink">Original Link:</label>
  			<input type="text" id="originalLink" name="originalLink"><br><br>
			<input type="submit" value="Update">
		</form>
	`
)

type RequestHandler struct {
	ls *LinkStore
}

func NewRequestHandler(ls *LinkStore) *RequestHandler {
	return &RequestHandler{ls: ls}
}

func (rh *RequestHandler) handleNotFound(w http.ResponseWriter, r *http.Request) {
	originalLink, err := rh.ls.GetLink(r.Context(), r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if originalLink == nil {
		fmt.Fprint(w, UpdateGoLinkHTML)
		return
	}
	http.Redirect(w, r, *originalLink, http.StatusFound)
}

func (rh *RequestHandler) handleHome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, UpdateGoLinkHTML)
}

func (rh *RequestHandler) handleApiUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "<h1>Error: %s</h1>\n", err)
		return
	}

	if err := rh.ls.PutLink(r.Context(), r.Form.Get("shortLink"), r.Form.Get("originalLink")); err != nil {
		fmt.Fprintf(w, "<h1>Error: %s</h1>\n", err)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}
