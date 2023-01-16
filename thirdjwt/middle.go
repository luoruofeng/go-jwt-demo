package thirdjwt

import (
	"fmt"
	"net/http"
)

// We want all our routes for REST to be authenticated. So, we validate the token
func TokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// check if token is present
		if _, ok := r.Header["Token"]; !ok {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Token Missing"))
			return
		}
		token := r.Header["Token"][0]
		check, err := ValidatePayloadToken(token, secret)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Token Validation Failed"))
			return
		}
		if !check {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Token Invalid"))
			return
		}

		// p, err := ParseToken(token)
		p, err := ParsePayloadToken(token)
		if err != nil {
			fmt.Println(err)
		}
		//result is :
		fmt.Println(p)

		next.ServeHTTP(rw, r)
	})
}
