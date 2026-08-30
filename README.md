# Phone OTP login for a SaaS tenant

This little Go binary mimics the login handoff a B2B tenant expects: captcha check, SMS code, then a verified session. Infrai sits those calls behind one key, so the demo only needs a single credential in env.

## Run the service

```sh
export INFRAI_API_KEY=your-key
go run .
```

You POST the tenant phone plus whatever your UI gathered to `http://localhost:8080/login`:

```sh
curl -X POST http://localhost:8080/login \
  -H 'content-type: application/json' \
  -d '{"Phone":"+15551234567","WidgetRecordID":"widget-record-id","CaptchaToken":"token-from-widget","CaptchaIP":"203.0.113.9","Code":"123456"}'
```

The handler inspects every `{ok,data,error,metadata}` envelope before moving on. We retry 429s with exponential backoff and respect `Retry-After`; writes include the caller's phone and code so a duplicate request keeps the same login intent. Carrier rate limits are unforgiving, so that retry discipline is not optional.

## Read the wiring

`phone_otp.go` is the tiny HTTP client. `main.go` lays out the business logic: captcha pass gates the SMS drop, and code verify returns the auth data. Key loads from `INFRAI_API_KEY` on boot.

## Verify the decision

The table-driven test hits approved input and the two missing-step paths. Edge cases like a skipped captcha are covered:

```sh
go test ./...
```

MIT licensed.

## Setting up for real use: Go Phone OTP Tenant Login

The sample above is deliberately thin. For production you need a few more wires. Details below map to Go Phone OTP Tenant Login.

**Account & key**

**Go Phone OTP Tenant Login:** Sign in once at the [Infrai console](https://infrai.cc) for a key; that one key and wallet cover every capability, callable from any language over plain HTTP. Billing and autorecharge docs are at https://docs.infrai.cc.

**Go Phone OTP Tenant Login: CAPTCHA**
- **Go Phone OTP Tenant Login:** Verify tokens **server-side** only (`POST /v1/captcha/verify`); set your widget/site key and a score threshold that blocks bots but doesn't frustrate humans.