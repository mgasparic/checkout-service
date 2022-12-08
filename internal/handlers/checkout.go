package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"service/internal/commons"
)

type CheckoutEnvironment commons.HandlerEnvironment

func (ce CheckoutEnvironment) Checkout(ctx *gin.Context) {
	var checkoutRequest commons.CheckoutRequest
	err := ctx.BindJSON(&checkoutRequest)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if commons.ContainsEmptyValues(checkoutRequest) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if commons.ContainsEmptyValues(checkoutRequest.Address) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response, err := http.Post(fmt.Sprintf("%s/reservation/%s", ce.StorageUrl, checkoutRequest.CartId), "text/html", nil)
	if err != nil {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	switch response.StatusCode {
	case http.StatusOK:
		// all good, we move forward
	case http.StatusNoContent:
		// there is no shopping cart in the database
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	default:
		// something went wrong, but status code is unexpected, it shouldn't happen
		log.Printf("unexpected status code from storage service: %d", response.StatusCode)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
		// in a more complex scenario, we could have more cases here - for example, if the items are out of stock,
		// we could return a different status code, such as 409
	}

	checkoutResponse := commons.CheckoutResponse{Shipping: []string{"DHL"}, Payment: []string{"Stripe"}}
	// we could have a more complex logic here for selecting the payment method and shipping method, we could even
	// include the list of items in that logic; however, for the showcase, simple logic is enough
	if checkoutRequest.Address.Country == "US" {
		checkoutResponse = commons.CheckoutResponse{Shipping: []string{"UPS", "DHL"}, Payment: []string{"Credit Card", "PayPal"}}
	}
	ctx.AbortWithStatusJSON(http.StatusOK, checkoutResponse)
}
