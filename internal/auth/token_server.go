package auth

import (
	"log"
	"net/http"
)

type authHandler struct {
	code string
}

func StartServer() (err error) {
	var ch authHandler
	port := "8080"

	mux := http.NewServeMux()
	mux.HandleFunc("*", ch.handleGetToken)

	tokenServer := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Getting auth code on port %s", port)
	log.Fatal(tokenServer.ListenAndServe())
	return
}

func (a *authHandler) handleGetToken(w http.ResponseWriter, r *http.Request) {
	a.code = r.URL.Query().Get("code")

}
