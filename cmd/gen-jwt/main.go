package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	secret := os.Args[1]
	userID, _ := strconv.ParseUint(os.Args[2], 10, 64)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),
		"email":   "probe@edairy.africa",
		"roles":   []string{},
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	})
	s, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	fmt.Print(s)
}
