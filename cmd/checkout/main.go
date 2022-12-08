package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/stripe/stripe-go/v74"
	"log"
	"net/http"
	"service/internal/commons"
	"service/internal/handlers"
)

type envVars struct {
	ServicePort        int    `envconfig:"SERVICE_PORT" default:"8080"`
	ServiceEnvironment string `envconfig:"SERVICE_ENVIRONMENT" default:"local"`
	StorageUrl         string `envconfig:"STORAGE_URL"`
	StripeKey          string `envconfig:"STRIPE_KEY"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var envVars envVars
	err := envconfig.Process("", &envVars)
	if err != nil {
		log.Fatal(err)
	}

	stripe.Key = envVars.StripeKey

	if envVars.ServiceEnvironment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	handler := commons.HandlerEnvironment{StorageUrl: envVars.StorageUrl}
	ce := handlers.CheckoutEnvironment(handler)
	pe := handlers.PaymentEnvironment(handler)

	router := gin.Default()
	router.POST("/checkout", ce.Checkout)
	router.POST("/payment", pe.Payment)
	router.GET("/", func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusOK)
	})

	log.Fatal(router.Run(fmt.Sprintf(":%d", envVars.ServicePort)))
}
