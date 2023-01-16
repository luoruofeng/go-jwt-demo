package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/luoruofeng/go-jwt-demo/basic"
)

func main() {
	mainRouter := mux.NewRouter()
	authRouter := mainRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", basic.SignupHandler)
	// The Signin will send the JWT back as we are making microservices.
	// The JWT token will make sure that other services are protected.
	// So, ultimately, we would need a middleware
	authRouter.HandleFunc("/signin", basic.SigninHandler)

	// Add the middleware to different subrouter
	otherRouter := mainRouter.PathPrefix("/other").Subrouter()
	otherRouter.HandleFunc("/home", basic.HomeHandler)
	otherRouter.Use(basic.TokenValidationMiddleware)

	// HTTP server
	// Add time outs
	server := &http.Server{
		Addr:    "0.0.0.0:8888",
		Handler: mainRouter,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error Booting the Server")
	}
}
