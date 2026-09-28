package auth

import (
	"context"
	"log"
	"net/http"
)

type authHandler struct {
	code     string
	exitChan chan struct{}
}

func GetAuthCode() (code string, err error) {
	ah := authHandler{
		exitChan: make(chan struct{}),
	}
	port := "8080"

	mux := http.NewServeMux()
	mux.HandleFunc("/", ah.handleGetCode)

	tokenServer := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("Getting auth code on port %s", port)
		err = tokenServer.ListenAndServe()
		log.Printf("Server finished: %v", err)
	}()

	<-ah.exitChan
	log.Println("Received shutdown trigger")
	err = tokenServer.Shutdown(context.Background())
	log.Println(err)

	if err == http.ErrServerClosed {
		err = nil
		code = ah.code
	}

	return
}

func (a *authHandler) handleGetCode(w http.ResponseWriter, r *http.Request) {
	a.code = r.URL.Query().Get("code")

	go func() {
		a.exitChan <- struct{}{}
	}()
}
