package userHandler

import (
	"net/http"
	"time"

	usermodels "github.com/KingSupermarket/model/userModels"
	requsermodel "github.com/KingSupermarket/model/userModels/reqUserModel"
	"github.com/KingSupermarket/pkg/security"
	"github.com/KingSupermarket/repository"
	repouseriml "github.com/KingSupermarket/repository/repo_user_iml"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUpHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataUser usermodels.Users
		if err := c.ShouldBind(&dataUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "faild bind data user",
			})
			return
		}
		validate := validator.New()
		if err := validate.Struct(dataUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "faild validator data user",
			})
			return
		}
		idUser, err := uuid.NewUUID()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "faild uuid data user",
			})
			return
		}
		role := usermodels.MEMBER.String()
		pwd := security.HashAndSalt([]byte(dataUser.Password))
		time := time.Now().UTC()
		newDataUser := usermodels.Users{
			Id:           primitive.NewObjectID(),
			User_id:      idUser.String(),
			Name:         dataUser.Name,
			Avatar:       dataUser.Avatar,
			Age:          dataUser.Age,
			Email:        dataUser.Email,
			Password:     pwd,
			Role:         role,
			Phone_Number: dataUser.Phone_Number,
			Created_At:   &time,
			Updated_At:   &time,
		}

		bus := repository.NewUserRepoIml(repouseriml.NewDB(db))
		data, err := bus.NewSignUp(c.Request.Context(), &newDataUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "faild db data user",
			})
			return
		}
		accessToken, refreshToken, err := security.JwtToken(&newDataUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "faild token",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result":        "signup successful",
			"data":          &data,
			"access token":  accessToken,
			"refresh token": refreshToken,
		})
	}
}

func SignInHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqSignin requsermodel.SigninModel
		if err := c.ShouldBind(&reqSignin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "faild bind data reqSignin",
			})
			return
		}
		validate := validator.New()
		if err := validate.Struct(reqSignin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "faild validator data reqSignin",
			})
			return
		}
		biz := repository.NewUserRepoIml(repouseriml.NewDB(db))
		foundUsers, err := biz.NewSignIn(c.Request.Context(), &reqSignin)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "faild db data signin",
			})
			return
		}
		//compare password
		isValidPassword := security.ComparePasswords(foundUsers.Password, []byte(reqSignin.Password))
		if !isValidPassword {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "faild password data signin",
			})
			return
		}
		//updated token
		acToken, reqToken, _ := security.UpdateToken(foundUsers)
		c.JSON(http.StatusOK, gin.H{
			"result":        "signin successful",
			"access token":  acToken,
			"refresh token": reqToken,
		})
	}
}

func ProfileUserHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdInterface, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user ID"})
			return
		}

		userId, ok := userIdInterface.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
			return
		}

		// Lấy thông tin người dùng từ database dựa trên userId
		userRepo := repository.NewUserRepoIml(repouseriml.NewDB(db))
		user, err := userRepo.NewProfileUser(c.Request.Context(), userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	}
}
func UpdateUserHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var UpdateUser usermodels.Users
		if err := c.ShouldBind(&UpdateUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "faild bind update data",
			})
			return
		}
		tokenData, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad user id",
			})
			return
		}
		claims, _ := tokenData.(string)
		bus := repository.NewUserRepoIml(repouseriml.NewDB(db))
		if err := bus.NewUpdateUser(c.Request.Context(), claims, &UpdateUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "faild update user",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"comment": "updated user successfully",
		})
	}
}

func DeleteUserHandler(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenData, ok := c.Get("userId")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad user id",
			})
			return
		}
		dataClaims := tokenData.(string)
		bus := repository.NewUserRepoIml(repouseriml.NewDB(db))
		if err := bus.NewDeleteUser(c.Request.Context(), dataClaims); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "faild dalete user",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"comment": "delete user successfully",
		})
	}
}
