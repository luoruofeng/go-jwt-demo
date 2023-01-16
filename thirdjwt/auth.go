package thirdjwt

import (
	"errors"
	"fmt"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/luoruofeng/go-jwt-demo/basic"
)

var (
	secret     = "Secure_Random_String"
	signMethod = jwt.SigningMethodHS256
)

type Payload struct {
	jwt.RegisteredClaims // 注意!这是jwt-go的v4版本新增的，原先是jwt.StandardClaims
	basic.User
	Phone string `json:"phone"`
}

func GenerateToken(claims jwt.Claims) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(signMethod, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, err
}

func ParseToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, Secret())
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return map[string]interface{}(claims), nil
	} else {
		return nil, errors.New("Convert token failed(type jwt.MapClaims)")
	}
}

func ParsePayloadToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, Secret())
	if err != nil {
		return nil, err
	}
	if payload, ok := token.Claims.(*Payload); ok {
		return payload, nil
	} else {
		return nil, errors.New("Convert token failed(type jwt.MapClaims)")
	}
}

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	}
}

func ValidateToken(tokenString string, secret string) (bool, error) {
	token, err := jwt.Parse(tokenString, Secret())

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return false, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return false, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return false, errors.New("token not active yet")
			} else {
				return false, errors.New("couldn't handle this token")
			}
		}
	}

	return token.Valid, nil
}

func ValidatePayloadToken(tokenString string, secret string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, Secret())

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return false, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return false, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return false, errors.New("token not active yet")
			} else {
				return false, errors.New("couldn't handle this token")
			}
		} else {
			return false, err
		}
	}

	return token.Valid, nil
}
