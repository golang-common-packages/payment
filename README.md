# Payment

## PayPal

### Unit test
```go
go test -v -cover -covermode=atomic ./...
```

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

### OpenID identity v1

* GET /v1/identity/openidconnect/userinfo/?schema=:schema
* POST /v1/identity/openidconnect/tokenservice (oauth or refresh token)

### Payment experience v1

* GET /v1/payment-experience/web-profiles
* POST /v1/payment-experience/web-profiles
* GET /v1/payment-experience/web-profiles/:id
* PUT /v1/payment-experience/web-profiles/:id
* DELETE /v1/payment-experience/web-profiles/:id

### Reporting v1

* POST /v1/reporting/transactions

### Vault v1

* POST /v1/vault/credit-cards
* DELETE /v1/vault/credit-cards/:id
* PATCH /v1/vault/credit-cards/:id
* GET /v1/vault/credit-cards/:id
* GET /v1/vault/credit-cards

### Checkout v2

* POST /v2/checkout/orders
* GET /v2/checkout/orders/:id
* PATCH /v2/checkout/orders/:id
* POST /v2/checkout/orders/:id/authorize
* POST /v2/checkout/orders/:id/capture

### Notification v1

* POST /v1/notifications/webhooks
* GET /v1/notifications/webhooks
* GET /v1/notifications/webhooks/:id
* PATCH /v1/notifications/webhooks/:id
* DELETE /v1/notifications/webhooks/:id
* POST /v1/notifications/verify-webhook-signature

### Catalog v1 (Product)

* POST /v1/catalogs/products
* PATCH /v1/catalogs/products/:id
* GET /v1/catalogs/products/:id
* GET /v1/catalogs/products

### Billing Plan v1

* POST  /v1/billing/plans
* PATCH /v1/billing/plans/:id
* GET   /v1/billing/plans/:id
* GET   /v1/billing/plans
* POST  /v1/billing/plans/:id/activate
* POST  /v1/billing/plans/:id/deactivate
* POST  /v1/billing/plans/:id/update-pricing-schemes

### Billing Subscription v1

* POST /v1/billing/subscriptions
* PATCH /v1/billing/subscriptions/:id
* GET /v1/billing/subscriptions/:id
* POST /v1/billing/subscriptions/:id/activate
* POST /v1/billing/subscriptions/:id/cancel
* POST /v1/billing/subscriptions/:id/revise
* POST /v1/billing/subscriptions/:id/capture
* POST /v1/billing/subscriptions/:id/suspend
* GET /v1/billing/subscriptions/:id/transactions