package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/esslamb/golang-hex/pkg/routes"
	"github.com/esslamb/golang-hex/pkg/user"

	"github.com/esslamb/golang-hex/pkg/database"
	"github.com/esslamb/golang-hex/pkg/utils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	port     = utils.GetEnv("PORT", "5000")
	mongoURI = utils.GetEnv("MONGOURI", "mongodb://mongo:27017")
)

func main() {
	// create mongodb connection
	m, err := database.OpenClientConnection(mongoURI)
	if err != nil {
		log.Error(errors.New("unable to connect to DB"))
		log.Error(err)
	}

	// create instance of any services required
	u := user.NewService(m)

	// create new instance of handler with services attached
	h := routes.NewHandler(u)

	// create new router and pass it the handler methods
	r := routes.CreateRouter(h)

	// create http server and pass it any services it needs
	srv := createHTTPServer(r)

	// start serve process
	log.Fatal(srv.ListenAndServe())
}

func createHTTPServer(r *mux.Router) *http.Server {
	log.Info("Creating HTTP Server")
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 20,
		ReadTimeout:  time.Second * 20,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	return srv
}
