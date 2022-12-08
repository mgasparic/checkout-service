# Checkout Service

The exact description of the task is provided in https://github.com/mgasparic/checkout-service/blob/main/docs/task.pdf.

## System

The entire e-commerce sample system is composed of multiple services. The checkout service is one of them, and it is the
sole focus of this task. The checkout service communicates with the storage synchronizer and with payment systems. It is
expected to receive http requests directly from clients. The components diagram is provided below, for easier
understanding.

![Components diagram](docs/checkout.jpg)

## Endpoints

| Path      | Method | Data                                                                                                                                                                                                                                                                    | Success Response                                  | Possible Error Codes                                                                                                        |
|:----------|:-------|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:--------------------------------------------------|:----------------------------------------------------------------------------------------------------------------------------|
| /checkout | POST   | {<br>cartId: string<br>address: {<br>firstName: string<br>lastName: string<br>street: string<br>houseNr: string<br>city: string<br>state: string<br>country: string<br>zip: string}<br>}                                                                                | {<br>shipping: []string<br>payment: []string<br>} | 204 - No Content<br>400 - Bad Request<br>500 - Internal Server Error<br>503 - Service Unavailable                           |
| /payment  | POST   | {<br>cartId: string<br>address: {<br>firstName: string<br>lastName: string<br>street: string<br>houseNr: string<br>city: string<br>state: string<br>country: string<br>zip: string}<br>shippingMethod: string<br>paymentIntent: {<br>id: string<br>method: string}<br>} | NA                                                | 204 - No Content<br>400 - Bad Request<br>402 - Payment Required<br>500 - Internal Server Error<br>503 - Service Unavailable |
