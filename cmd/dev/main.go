package main

import (
	"github.com/KingSupermarket/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	router.NewRouter(r)
	r.Run(":3251")
}
