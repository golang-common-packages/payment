# Payment

## PayPal

### Auth

* POST /v1/oauth2/token

### /v1/payments

* POST /v1/payments/payouts
* GET /v1/payments/payouts/:id
* GET /v1/payments/payouts-item/:id
* POST /v1/payments/payouts-item/:id/cancel
* GET /v1/payments/sale/:id
* POST /v1/payments/sale/:id/refund

### /v2/payments

* GET /v2/payments/authorizations/:id
* GET /v2/payments/captures/:id
* POST /v2/payments/authorizations/:id/capture
* POST /v2/payments/authorizations/:id/void
* POST /v2/payments/authorizations/:id/reauthorize
* GET /v2/payments/refund/:id