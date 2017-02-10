package router

import "net/http"

// Route structure defines the standard interface for API interaction
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(w http.ResponseWriter, r *http.Request) (status int, body []byte, err error)
}

// Routes are a collection of Route structures
type Routes []Route
