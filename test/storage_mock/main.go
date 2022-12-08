package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"time"
)

type envVars struct {
	ServicePort int `envconfig:"SERVICE_PORT" default:"9000"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var envVars envVars
	err := envconfig.Process("", &envVars)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/reservation/:cartId", func(ctx *gin.Context) {
		switch ctx.Param("cartId") {
		case "uuid_existing":
			ctx.AbortWithStatusJSON(http.StatusOK, []interface{}{
				gin.H{
					"itemId": "sku_123",
					"amount": 1,
					"price":  0.99,
				},
				gin.H{
					"itemId": "sku_456",
					"amount": 4,
					"price":  2.50,
				},
			})
		case "uuid_non_existing":
			ctx.AbortWithStatus(http.StatusNoContent)
		case "uuid_timeout":
			time.Sleep(time.Minute)
		default:
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	})
	router.POST("/reservation/:cartId", func(ctx *gin.Context) {
		switch ctx.Param("cartId") {
		case "uuid_existing":
			ctx.AbortWithStatus(http.StatusOK)
		case "uuid_non_existing":
			ctx.AbortWithStatus(http.StatusNoContent)
		case "uuid_timeout":
			time.Sleep(time.Minute)
		default:
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	})
	router.POST("/order/:cartId", func(ctx *gin.Context) {
		switch ctx.Param("cartId") {
		case "uuid_existing":
			ctx.AbortWithStatus(http.StatusOK)
		case "uuid_timeout":
			time.Sleep(time.Minute)
		default:
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	})

	log.Fatal(router.Run(fmt.Sprintf(":%d", envVars.ServicePort)))
}
