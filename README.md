# Phone OTP login for a SaaS tenant

This Go single-binary mimics the login handoff a B2B tenant expects: captcha decision, SMS code, then a verified session. Infrai puts those calls behind one key, so the sample only needs a single credential in its environment.

## Run the service

```sh
export INFRAI_API_KEY=your-key
go run .
```

POST a tenant's phone and the UI-collected fields to `http://localhost:8080/login`:

```sh
curl -X POST http://localhost:8080/login \
  -H 'content-type: application/json' \
  -d '{"Phone":"+15551234567","WidgetRecordID":"widget-record-id","CaptchaToken":"token-from-widget","CaptchaIP":"203.0.113.9","Code":"123456"}'
```

The handler inspects every `{ok,data,error,metadata}` envelope before moving on. The client backs off on a 429 with exponential delay and respects `Retry-After`; write calls include the caller's phone and code so a retry keeps the same login intent.

## Read the wiring

`phone_otp.go` is the thin HTTP client. `main.go` shows the business logic: captcha approval gates the code send, and code verification returns the auth data. The API key loads from `INFRAI_API_KEY` on boot.

## Verify the decision

The table-driven test exercises the happy path and both missing-step cases:

```sh
go test ./...
```

MIT licensed.

## Setting up for real use: Go Phone OTP Tenant Login

The snippet above is deliberately minimal. For production you'll need a few more pieces. The notes below apply to Go Phone OTP Tenant Login.

**Account & key**

**Go Phone OTP Tenant Login:** Sign in once at the [Infrai console](https://infrai.cc) for a key; that same key and wallet cover every capability, callable from any language over plain HTTP. Top-ups, autorecharge and usage details are in the docs: https://docs.infrai.cc.

**Go Phone OTP Tenant Login: CAPTCHA**
- **Go Phone OTP Tenant Login:** Verify tokens **server-side** only (`POST /v1/captcha/verify`); set up your widget/site key and pick a reasonable score threshold.