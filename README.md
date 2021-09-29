# Payment

## PayPal

### Auth v1

* POST /v1/oauth2/token

### Payment v1

* POST /v1/payments/payouts
* GET /v1/payments/payouts/:id
* GET /v1/payments/payouts-item/:id
* POST /v1/payments/payouts-item/:id/cancel
* GET /v1/payments/sale/:id
* POST /v1/payments/sale/:id/refund
* GET /v1/payments/billing-plans
* POST /v1/payments/billing-plans
* PATCH /v1/payments/billing-plans/:id
* POST /v1/payments/billing-agreements
* POST /v1/payments/billing-agreements/:token/agreement-execute

### Payment v2

* GET /v2/payments/authorizations/:id
* POST /v2/payments/authorizations/:id/capture
* POST /v2/payments/authorizations/:id/void
* POST /v2/payments/authorizations/:id/reauthorize
* GET /v2/payments/captures/:id
* GET /v2/payments/refund/:id

### Identity v1
* GET /v1/identity/openidconnect/userinfo/?schema=:schema
* POST /v1/identity/openidconnect/tokenservice (oauth or refresh token)