package middlewares

import (
	"log"
	"strings"
	"time"
	"todo/infrastructure"
	"todo/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
)

func GetTokenString(user *model.User) (string, string, error) {
	claim := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"password": user.Password,
		"role":     user.Role,
	}

	log.Println(user)
	refreshClaim := jwt.MapClaims{
		"id": user.ID,
	}

	jwtauth.SetExpiry(claim, time.Now().Local().Add(time.Hour*time.Duration(infrastructure.GetExtendAccessHour())))
	jwtauth.SetExpiry(refreshClaim, time.Now().Local().Add(time.Hour*time.Duration(infrastructure.GetExtendRefreshHour())))
	_, tokenString, _ := infrastructure.GetEncodeAuth().Encode(claim)
	_, refreshToken, _ := infrastructure.GetEncodeAuth().Encode(refreshClaim)
	log.Println(tokenString)
	tokenString = "Bearer " + tokenString
	refreshToken = "Bearer " + refreshToken
	return tokenString, refreshToken, nil

}

func GetClaimsData(tokenStr string) (*model.User, error) {
	user := model.User{}
	words := strings.Fields(tokenStr)
	if len(words) == 1 {
		token, err := jwt.ParseWithClaims(words[0], &user, func(token *jwt.Token) (interface{}, error) {
			return infrastructure.GetPublicKey(), nil
		})
		if err != nil || !token.Valid {
			infrastructure.ErrLog.Println("Problem getting claims data: ", err)
			return nil, err
		}
	} else {
		token, err := jwt.ParseWithClaims(words[1], &user, func(token *jwt.Token) (interface{}, error) {
			return infrastructure.GetPublicKey(), nil
		})
		if err != nil {
			infrastructure.ErrLog.Println("Problem getting claims data: ", err)
			return nil, err
		}
		if !token.Valid {
			infrastructure.ErrLog.Println("Problem getting claims data: ", err)
			return nil, err
		}
		infrastructure.InfoLog.Println("get claims 1: ", token.Valid, user.Username, user.Role)
	}

	return &user, nil
}
