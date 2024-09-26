package userHandler

import (
	"net/http"
	"strings"
	"time"

	usermodels "github.com/KingSupermarket/model/userModels"
	requsermodel "github.com/KingSupermarket/model/userModels/reqUserModel"
	"github.com/KingSupermarket/pkg/redisDB"
	"github.com/KingSupermarket/pkg/security"
	"github.com/KingSupermarket/repository"
	repouseriml "github.com/KingSupermarket/repository/repo_user_iml"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUpHandler(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var dataUser usermodels.Users

		if err := c.BodyParser(&dataUser); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to bind data user",
			})
		}

		validate := validator.New()
		if err := validate.Struct(dataUser); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to validate data user",
			})
		}

		idUser, err := uuid.NewUUID()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to generate uuid for data user",
			})
		}

		role := usermodels.MEMBER.String()
		pwd := security.HashAndSalt([]byte(dataUser.Password))
		now := time.Now().UTC()

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
			Created_At:   &now,
			Updated_At:   &now,
		}

		bus := repository.NewUserRepoIml(repouseriml.NewDB(db))
		data, err := bus.NewSignUp(c.Context(), &newDataUser)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to save data user to db",
			})
		}

		accessToken, refreshToken, err := security.JwtToken(&newDataUser)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to generate token",
			})
		}

		r := redisDB.GetInstanceRedis()
		if err := r.SaveRefreshToken(refreshToken); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to save refresh token",
			})
		}

		return c.JSON(fiber.Map{
			"result":        "signup successful",
			"data":          data,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func SignInHandler(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		refreshToken := c.Get("Authorization")
		if refreshToken != "" {
			refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")
		}

		// Nếu có refresh token, thử cập nhật access token
		if refreshToken != "" {
			newAccessToken, err := security.UpdateToken(refreshToken)
			if err == nil {
				// Nếu cập nhật thành công, trả về access token mới
				return c.Status(http.StatusOK).JSON(fiber.Map{
					"result":       "signin successful",
					"access_token": newAccessToken,
				})
			}
		}

		var reqSignin requsermodel.SigninModel
		if err := c.BodyParser(&reqSignin); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to bind data reqSignin",
			})
		}

		validate := validator.New()
		if err := validate.Struct(reqSignin); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to validate data reqSignin",
			})
		}

		biz := repository.NewUserRepoIml(repouseriml.NewDB(db))
		foundUsers, err := biz.NewSignIn(c.Context(), &reqSignin)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to retrieve user data from database",
			})
		}

		// So sánh mật khẩu
		isValidPassword := security.ComparePasswords(foundUsers.Password, []byte(reqSignin.Password))
		if !isValidPassword {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "failed password verification",
			})
		}

		// Cập nhật token
		acToken, reqToken, _ := security.JwtToken(foundUsers)
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"result":        "signin successful",
			"access_token":  acToken,
			"refresh_token": reqToken,
		})
	}
}

func ProfileUserHandler(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIdInterface := c.Locals("userId") // Lấy ID người dùng
		userId, ok := userIdInterface.(string)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid user ID"})
		}

		// Lấy thông tin người dùng từ database
		userRepo := repository.NewUserRepoIml(repouseriml.NewDB(db))
		user, err := userRepo.NewProfileUser(c.Context(), userId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"data": user})
	}
}
func UpdateUserHandler(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var updateUser usermodels.Users
		if err := c.BodyParser(&updateUser); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "failed to bind update data"})
		}

		tokenData := c.Locals("userId")  // Lấy giá trị từ Locals
		claims, ok := tokenData.(string) // Kiểm tra kiểu dữ liệu
		if !ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "bad user id"})
		}

		bus := repository.NewUserRepoIml(repouseriml.NewDB(db))
		if err := bus.NewUpdateUser(c.Context(), claims, &updateUser); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "failed to update user"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"comment": "updated user successfully"})
	}
}

func DeleteUserHandler(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenData := c.Locals("userId") // Lấy ID người dùng
		dataClaims, ok := tokenData.(string)
		if !ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "bad user id",
			})
		}

		bus := repository.NewUserRepoIml(repouseriml.NewDB(db))
		if err := bus.NewDeleteUser(c.Context(), dataClaims); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to delete user",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"comment": "delete user successfully",
		})
	}
}
func HandlerHistory(db *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenData := c.Locals("userId") // Lấy ID người dùng
		dataClaims, ok := tokenData.(string)
		if !ok {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "bad user id",
			})
		}

		bus := repository.NewUserRepoIml(repouseriml.NewDB(db))
		users, history, err := bus.NewHistoryPurchases(c.Context(), dataClaims)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to retrieve user history",
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"comment":        "history retrieved successfully",
			"User":           users,
			"historyPurchas": history,
		})
	}
}
