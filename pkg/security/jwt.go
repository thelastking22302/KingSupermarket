package security

import (
	"errors"
	"fmt"
	"os"
	"time"

	usermodels "github.com/KingSupermarket/model/userModels"
	"github.com/KingSupermarket/pkg/logger"
	"github.com/KingSupermarket/pkg/redisDB"
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

func UpdateToken(refreshtoken string) (string, error) {
	//new token
	claims, err := ValidateToken(refreshtoken)
	if err != nil {
		return "", errors.New("refreshtoken khong hop le")
	}
	//kiem tra refresh token con han hay khong
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return "", errors.New("refresh token het han ban phai dang nhap lai")
	}
	newAccessToken := &usermodels.Token{
		User_Id: claims.User_Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessToken).SignedString([]byte(keyJwt))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ValidateToken(userToken string) (*usermodels.Token, error) {
	//xac thuc token
	log := logger.GetLogger()
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
			log.Warnf("token expiration")
		}
		//check token redis
		redisInstance := redisDB.GetInstanceRedis()
		exists, err := redisInstance.CheckRefreshToken()
		if err != nil {
			log.Errorf("error checking refresh token: %v", err)
			return nil, err
		}

		if !exists {
			log.Errorf("Refresh token invalid on Redis.")
			return nil, fmt.Errorf("refresh token not found")
		}
		return claims, nil
	} else {
		panic("invalid claims")
	}
}
