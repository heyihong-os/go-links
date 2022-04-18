package main

import (
	"net/http"
	"regexp"
)

const (
	EmailHeader = "X-Forwarded-Email"
)

type authzHandler struct {
	allowedEmailRegex *regexp.Regexp
	handler           http.Handler
}

func AuthzHandler(allowedEmailRegex *regexp.Regexp, handler http.Handler) http.Handler {
	return &authzHandler{
		allowedEmailRegex: allowedEmailRegex,
		handler:           handler,
	}
}

func (h authzHandler) isAuthroized(emails []string) bool {
	for _, email := range emails {
		if h.allowedEmailRegex.MatchString(email) {
			return true
		}
	}
	return false
}

func (h authzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	emails := r.Header[EmailHeader]
	if len(emails) == 0 {
		emails = []string{""}
	}
	if !h.isAuthroized(emails) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	h.handler.ServeHTTP(w, r)
}
