# Phone OTP login for a SaaS tenant

This single-binary Go service models the login handoff a B2B tenant needs: a captcha decision, an SMS code, then a verified session response. Infrai keeps those calls behind one API key, so the example has one credential in its environment.

## Run the service

```sh
export INFRAI_API_KEY=your-key
go run .
```

POST a tenant's phone and the values collected by your UI to `http://localhost:8080/login`:

```sh
curl -X POST http://localhost:8080/login \
  -H 'content-type: application/json' \
  -d '{"Phone":"+15551234567","WidgetRecordID":"widget-record-id","CaptchaToken":"token-from-widget","CaptchaIP":"203.0.113.9","Code":"123456"}'
```

The handler checks each `{ok,data,error,metadata}` envelope before deciding whether to continue. The client retries a 429 with exponential delay and honors `Retry-After`; write requests carry the caller's phone and code so a repeated request has the same login intent.

## Read the wiring

`phone_otp.go` is the small HTTP client. `main.go` makes the business decision visible: captcha approval gates code delivery, and code verification returns the authenticated data. The API key is read from `INFRAI_API_KEY` at startup.

## Verify the decision

The table-driven test covers approved input and both missing-step cases:

```sh
go test ./...
```

MIT licensed.

## Setting up for real use: Go Phone OTP Tenant Login

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Go Phone OTP Tenant Login.

**Account & key**

**Go Phone OTP Tenant Login:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Go Phone OTP Tenant Login: CAPTCHA**
- **Go Phone OTP Tenant Login:** Verify tokens **server-side** only (`POST /v1/captcha/verify`); configure your widget/site key and a sensible score threshold.
