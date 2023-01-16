package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	basic "github.com/luoruofeng/go-jwt-demo/basic"
	"github.com/luoruofeng/go-jwt-demo/thirdjwt"
)

func SignupHandler(rw http.ResponseWriter, r *http.Request) {
	// extra error handling should be done at server side to prevent malicious attacks
	if _, ok := r.Header["Email"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Email Missing"))
		return
	}
	if _, ok := r.Header["Username"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Username Missing"))
		return
	}
	if _, ok := r.Header["Passwordhash"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Passwordhash Missing"))
		return
	}
	if _, ok := r.Header["Fullname"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Fullname Missing"))
		return
	}

	// validate and then add the user
	check := basic.AddUserObject(r.Header["Email"][0], r.Header["Username"][0], r.Header["Passwordhash"][0],
		r.Header["Fullname"][0], 0)
	// if false means username already exists
	if !check {
		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte("Email or Username already exists"))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("User Created"))
}

func SigninHandler(rw http.ResponseWriter, r *http.Request) {
	// validate the request first.
	if _, ok := r.Header["Email"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Email Missing"))
		return
	}
	if _, ok := r.Header["Passwordhash"]; !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Passwordhash Missing"))
		return
	}
	// let’s see if the user exists
	valid, err := basic.ValidateUser(r.Header["Email"][0], r.Header["Passwordhash"][0])
	if err != nil {
		// this means either the user does not exist
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("User Does not Exist"))
		return
	}

	if !valid {
		// this means the password is wrong
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Incorrect Password"))
		return
	}
	user := basic.User{Email: r.Header["Email"][0], Passwordhash: r.Header["Passwordhash"][0]}
	tokenString, err := getSignedToken(user)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Internal Server Error"))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(tokenString))
}

func getSignedToken(user basic.User) (string, error) {
	// we make a JWT Token here with signing method of ES256 and claims.
	// claims are attributes.
	claimsMap := thirdjwt.Payload{
		Phone: "18280025374",
		User:  user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 30)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                       // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                       // 生效时间
		},
	}
	// here we provide the shared secret. It should be very complex.
	// Also, it should be passed as a System Environment variable

	tokenString, err := thirdjwt.GenerateToken(claimsMap)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("home"))
}

func main() {
	mainRouter := mux.NewRouter()
	authRouter := mainRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", SignupHandler)
	// The Signin will send the JWT back as we are making microservices.
	// The JWT token will make sure that other services are protected.
	// So, ultimately, we would need a middleware
	authRouter.HandleFunc("/signin", SigninHandler)

	// Add the middleware to different subrouter
	otherRouter := mainRouter.PathPrefix("/other").Subrouter()
	otherRouter.HandleFunc("/home", HomeHandler)
	otherRouter.Use(thirdjwt.TokenValidationMiddleware)

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
