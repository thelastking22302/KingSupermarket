package middleware

import (
	"net/http"

	usermodels "github.com/KingSupermarket/model/userModels"
	"github.com/gin-gonic/gin"
)

func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userAdmin usermodels.Users
		if err := c.ShouldBind(&userAdmin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bind UserAdmin faild",
			})
			return
		}

		if userAdmin.Role != "ADMIN" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "You do not have permission to access the data",
			})
			return
		}
		c.Next()
	}
}
