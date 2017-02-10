package router

import (
	"net/http"

	"fmt"
	"time"

	"github.com/gorilla/mux"

	log "github.com/Sirupsen/logrus"
)

type appRequest struct {
	AppConfig interface{}
	Log       *log.Logger
	Route
	H func(w http.ResponseWriter, r *http.Request) (status int, body []byte, err error)
}

func (ar appRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	status, body, err := ar.H(w, r)

	if err != nil {
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			ar.Log.Info(fmt.Sprintf("Status returning internal error: %d", status))
			http.Error(w, http.StatusText(status), status)
		default:
			ar.Log.Info(fmt.Sprintf("Status returning something else error: %d", status))
			http.Error(w, http.StatusText(status), status)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(status)
		w.Write(body)
	}

	ar.Log.WithFields(log.Fields{
		"Error":       err,
		"Type":        "request.run",
		"Path":        r.RequestURI,
		"RequestType": r.Method,
		"RespCode":    status,
		"LogDate":     start,
	}).Info(ar.Route.Name)
}

func NewRouter(logger *log.Logger, routes Routes) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		logger.Info(fmt.Sprintf("Adding route %s with type %s", route.Pattern, route.Method))

		ar := appRequest{
			Log:   logger,
			Route: route,
			H:     route.HandlerFunc,
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(ar)

	}

	return router
}
