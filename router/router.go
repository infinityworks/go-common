package router

import (
	"io"
	"io/ioutil"
	"net/http"

	"fmt"
	"time"

	"github.com/gorilla/mux"
	"github.com/infinityworksltd/go-common/metrics"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

type appRequest struct {
	Log *log.Logger
	Route
	H func(w http.ResponseWriter, r *http.Request) (status int, body []byte, err error)
}

type errorResponse struct {
	Error string    `json:"error"`
	Time  time.Time `json:"time"`
}

func (ar appRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	status, body, err := ar.H(w, r)

	defer func(begun time.Time) {
		metrics.Instrument(
			time.Since(begun).Seconds(),
			status,
			ar.Route.Method,
			ar.Route.Name,
		)
	}(time.Now())

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(status)

	ar.Log.WithFields(log.Fields{
		"Error":       err,
		"Type":        "request.run",
		"Path":        r.RequestURI,
		"RequestType": r.Method,
		"RespCode":    status,
		"LogDate":     start,
	}).Info(ar.Route.Name)

	if err != nil {
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
			return
		case http.StatusInternalServerError:
			ar.Log.Errorf("Status returning internal error: %d", status)
			writeError(err, w)
			return
		default:
			ar.Log.Errorf("Status returning something else error: %d", status)
			writeError(err, w)
			return
		}
	}

	w.Write(body)
}

func writeError(e error, w http.ResponseWriter) {

	errResp := errorResponse{Error: fmt.Sprintf("%s", e), Time: time.Now()}
	b, err := json.Marshal(errResp)

	if err != nil {
		m, _ := json.Marshal("System failure, could not marshall out error")

		w.Write(m)
		return
	}

	w.Write(b)
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

	registerMetrics(router)

	return router
}

func registerMetrics(router *mux.Router) {
	metrics.Init()

	handler := promhttp.Handler()

	router.
		Methods("GET").
		Path("/metrics").
		Name("Metrics").
		Handler(handler)
}

// UnmarshalBody Accepts an io.ReadCloser (usually a HTTP request body) and an interface to unmarshal the request into.
func UnmarshalBody(body io.ReadCloser, s interface{}) error {

	b, err := ioutil.ReadAll(io.LimitReader(body, 1048576))

	if err != nil {
		return fmt.Errorf("Could not read the JSON request body. Error: %s", err)
	}

	if err := body.Close(); err != nil {
		return fmt.Errorf("Could not close the JSON request body. Error: %s", err)
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("Could not unmarshal the request body into the struct you gave us. Error: %s", err)
	}

	return nil
}

// MarshalBody takes an interface to Marshal as JSON and returns it, it also handles returning of an error state
func MarshalBody(s interface{}) (status int, body []byte, err error) {
	out, err := json.Marshal(s)

	if err != nil {
		err = errors.Wrap(err, "Could not conver the response into JSON")
		return http.StatusInternalServerError, []byte(""), err
	}

	return http.StatusOK, out, nil
}
