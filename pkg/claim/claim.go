package claim

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Claim struct {
	// add the standard fields from the library
	jwt.StandardClaims
	// add the user's id
	ID int `json:"id"`
	// Expiration time defined in claim
	ExpiresAt int64 `json:"exp"`
}

// GenerateJWT returns a token with the claim.
func (c *Claim) GenerateJWT(signingString string) (string, error) {
	// set the expiration time to 15'
	c.ExpiresAt = time.Now().Add(time.Minute * 15).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(signingString))
}

// GetFromToken returns a claim from a token.
func GetFromToken(tokenString, signingString string) (*Claim, error) {
	token, err := jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) {
		return []byte(signingString), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claim")
	}

	iID, ok := claim["id"]
	if !ok {
		return nil, errors.New("user id not found")
	}

	id, ok := iID.(float64)
	if !ok {
		return nil, errors.New("invalid user id")
	}

	return &Claim{ID: int(id)}, nil
}