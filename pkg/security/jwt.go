package security

import (
	"fmt"
	"os"
	"time"

	usermodels "github.com/KingSupermarket/model/userModels"
	"github.com/dgrijalva/jwt-go"
)

var keyJwt string = os.Getenv("SECRET_KEY")

func JwtToken(data *usermodels.Users) (string, string, error) {
	//thanh phan cua 1 token
	newClaimsAccess := &usermodels.Token{
		User_Id: data.User_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		},
	}
	newClaimsRefresh := &usermodels.Token{
		User_Id: data.User_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24*7)).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		},
	}
	//khoi tao token
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaimsAccess).SignedString([]byte(keyJwt))
	if err != nil {
		panic(err)
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaimsRefresh).SignedString([]byte(keyJwt))
	if err != nil {
		panic(err)
	}
	return accessToken, refreshToken, nil
}

func UpdateToken(data *usermodels.Users) (string, string, error) {
	//new token
	newAccessToken, newRefreshToken, err := JwtToken(data)
	if err != nil {
		return "", "", err
	}
	return newAccessToken, newRefreshToken, nil
}

func ValidateToken(userToken string) (*usermodels.Token, error) {
	//xac thuc token
	token, err := jwt.ParseWithClaims(
		userToken,
		&usermodels.Token{},
		func(token *jwt.Token) (interface{}, error) { //callback cung cap khoa bi mat(keyJwt) xac thuc token
			return []byte(keyJwt), nil
		},
	)
	if err != nil {
		return nil, err
	}
	//trich xuat cac claims duoc token xac thuc
	if claims, ok := token.Claims.(*usermodels.Token); ok {
		// Check token expiration
		if claims.ExpiresAt < time.Now().Local().Unix() {
			fmt.Println("claims token expires")
		}
		return claims, nil
	} else {
		panic("invalid claims")
	}
}
