package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"service/internal/commons"
)

type PaymentEnvironment commons.HandlerEnvironment

func (pe PaymentEnvironment) Payment(ctx *gin.Context) {
	var paymentRequest commons.PaymentRequest
	err := ctx.BindJSON(&paymentRequest)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if commons.ContainsEmptyValues(paymentRequest) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if commons.ContainsEmptyValues(paymentRequest.Address) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if commons.ContainsEmptyValues(paymentRequest.PaymentIntent) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if paymentRequest.ShippingMethod != "DHL" {
		// we don't support any other shipping method than DHL, in this showcase
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if paymentRequest.PaymentIntent.Method != "Stripe" {
		// we don't support any other payment methods than Stripe, in this showcase
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	pi, err := paymentintent.Get(paymentRequest.PaymentIntent.Id, nil)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response, err := http.Get(fmt.Sprintf("%s/reservation/%s", pe.StorageUrl, paymentRequest.CartId))
	if err != nil {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	switch response.StatusCode {
	case http.StatusOK:
		// all good, we move forward
	case http.StatusNoContent:
		// there is no reservation in the database
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	default:
		// something went wrong, but status code is unexpected, it shouldn't happen
		log.Printf("unexpected status code from storage service: %d", response.StatusCode)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var orderResponse []commons.CartItem
	err = json.Unmarshal(body, &orderResponse)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// for the showcase we assume that all the values are in USD
	var price float64
	for _, item := range orderResponse {
		price += item.Price * float64(item.Amount)
	}
	if pi.Amount != int64(math.Ceil(price*100)) {
		// the amount in the payment intent doesn't match the amount in the order
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if pi.Status != "requires_capture" { // we will capture the funds once the order is shipped
		// the payment intent is not in the correct state
		ctx.AbortWithStatus(http.StatusPaymentRequired)
		return
	}

	bodyRaw, err := json.Marshal(commons.OrderRequest{
		Address:        paymentRequest.Address,
		ShippingMethod: paymentRequest.ShippingMethod,
	})
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	response, err = http.Post(fmt.Sprintf("%s/order/%s", pe.StorageUrl, paymentRequest.CartId), "application/json", bytes.NewReader(bodyRaw))
	if err != nil {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}
