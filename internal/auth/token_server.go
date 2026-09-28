package auth

import (
	"context"
	"fmt"
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
	if err != nil {
		return "", fmt.Errorf("unable to capture authorization code: %w", err)
	}

	code = ah.code
	fmt.Printf("Authorization code captured: %s\n", code)

	return
}

func (a *authHandler) handleGetCode(w http.ResponseWriter, r *http.Request) {
	a.code = r.URL.Query().Get("code")

	//TODO: Return JS to attempt to close window

	go func() {
		a.exitChan <- struct{}{}
	}()
}
