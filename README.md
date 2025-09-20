<p align="center">
    <a href="https://pocketbase.io" target="_blank" rel="noopener">
        <img src="https://i.imgur.com/5qimnm5.png" alt="PocketBase - open source backend in 1 file" />
    </a>
</p>

<p align="center">
    <a href="https://github.com/pocketbase/pocketbase/actions/workflows/release.yaml" target="_blank" rel="noopener"><img src="https://github.com/pocketbase/pocketbase/actions/workflows/release.yaml/badge.svg" alt="build" /></a>
    <a href="https://github.com/pocketbase/pocketbase/releases" target="_blank" rel="noopener"><img src="https://img.shields.io/github/release/pocketbase/pocketbase.svg" alt="Latest releases" /></a>
    <a href="https://pkg.go.dev/github.com/pocketbase/pocketbase" target="_blank" rel="noopener"><img src="https://godoc.org/github.com/pocketbase/pocketbase?status.svg" alt="Go package documentation" /></a>
</p>

[PocketBase](https://pocketbase.io) is an open source Go backend that includes:

- embedded database (_SQLite_) with **realtime subscriptions**
- built-in **files and users management**
- convenient **Admin dashboard UI**
- and simple **REST-ish API**

## 🚀 Enhanced Features

This fork extends PocketBase with additional features:

- **📱 WhatsApp Business OTP Integration** - Send OTP codes via WhatsApp Business API
- **📞 Phone Field Support** - Automatic phone field addition to users collection
- **🔧 Multi-channel OTP Delivery** - Support for email, WhatsApp, or both delivery methods

> [!NOTE]
> **Testing Status**: These enhanced features are currently in development and have not been fully tested in production environments. Please use with caution and test thoroughly before deploying to production.

**For documentation and examples, please visit https://pocketbase.io/docs.**

> [!WARNING]
> Please keep in mind that PocketBase is still under active development
> and therefore full backward compatibility is not guaranteed before reaching v1.0.0.

## API SDK clients

The easiest way to interact with the PocketBase Web APIs is to use one of the official SDK clients:

- **JavaScript - [pocketbase/js-sdk](https://github.com/pocketbase/js-sdk)** (_Browser, Node.js, React Native_)
- **Dart - [pocketbase/dart-sdk](https://github.com/pocketbase/dart-sdk)** (_Web, Mobile, Desktop, CLI_)

You could also check the recommendations in https://pocketbase.io/docs/how-to-use/.


## Overview

### Use as standalone app

You could download the prebuilt executable for your platform from the [Releases page](https://github.com/pocketbase/pocketbase/releases).
Once downloaded, extract the archive and run `./pocketbase serve` in the extracted directory.

The prebuilt executables are based on the [`examples/base/main.go` file](https://github.com/pocketbase/pocketbase/blob/master/examples/base/main.go) and comes with the JS VM plugin enabled by default which allows to extend PocketBase with JavaScript (_for more details please refer to [Extend with JavaScript](https://pocketbase.io/docs/js-overview/)_).

### Use as a Go framework/toolkit

PocketBase is distributed as a regular Go library package which allows you to build
your own custom app specific business logic and still have a single portable executable at the end.

Here is a minimal example:

0. [Install Go 1.23+](https://go.dev/doc/install) (_if you haven't already_)

1. Create a new project directory with the following `main.go` file inside it:
    ```go
    package main

    import (
        "log"

        "github.com/pocketbase/pocketbase"
        "github.com/pocketbase/pocketbase/core"
    )

    func main() {
        app := pocketbase.New()

        app.OnServe().BindFunc(func(se *core.ServeEvent) error {
            // registers new "GET /hello" route
            se.Router.GET("/hello", func(re *core.RequestEvent) error {
                return re.String(200, "Hello world!")
            })

            return se.Next()
        })

        if err := app.Start(); err != nil {
            log.Fatal(err)
        }
    }
    ```

2. To init the dependencies, run `go mod init myapp && go mod tidy`.

3. To start the application, run `go run main.go serve`.

4. To build a statically linked executable, you can run `CGO_ENABLED=0 go build` and then start the created executable with `./myapp serve`.

_For more details please refer to [Extend with Go](https://pocketbase.io/docs/go-overview/)._

## 📱 WhatsApp OTP Integration

This fork includes WhatsApp Business API integration for sending OTP codes. Here's how to set it up:

### Prerequisites

1. **WhatsApp Business API Account** - You need a Meta Business account with WhatsApp Business API access
2. **Access Token** - Get your access token from Meta for Developers
3. **Phone Number ID** - Your WhatsApp Business phone number ID

### Configuration

1. **Admin Dashboard Setup:**
   - Go to Settings > Application
   - Fill in the WhatsApp Business API section:
     - **WhatsApp Access Token**: Your Meta access token
     - **Phone Number ID**: Your WhatsApp Business phone number ID

2. **Collection OTP Settings:**
   - Go to Collections > users (or your auth collection)
   - Enable OTP in the Auth options
   - Set **Delivery Method** to:
     - `email` - Send OTP via email only
     - `whatsapp` - Send OTP via WhatsApp only  
     - `both` - Send OTP via both email and WhatsApp

3. **Customize WhatsApp Template:**
   - Edit the WhatsApp message template in collection settings
   - Use placeholders: `{OTP}`, `{APP_NAME}`, `{RECORD_EMAIL}`, `{RECORD_ID}`

### API Usage

```bash
# Request OTP via WhatsApp
curl -X POST http://localhost:8090/api/collections/users/auth-with-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "otp": "12345678"
  }'
```

### Example Implementation

See `examples/whatsapp_otp_example.go` for a complete implementation example.

> [!WARNING]
> **Production Readiness**: This WhatsApp integration is experimental and requires thorough testing before production use. Ensure you have proper error handling and fallback mechanisms in place.

### Building and running the repo main.go example

To build the minimal standalone executable, like the prebuilt ones in the releases page, you can simply run `go build` inside the `examples/base` directory:

0. [Install Go 1.23+](https://go.dev/doc/install) (_if you haven't already_)
1. Clone/download the repo
2. Navigate to `examples/base`
3. Run `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build`
   (_https://go.dev/doc/install/source#environment_)
4. Start the created executable by running `./base serve`.

Note that the supported build targets by the pure Go SQLite driver at the moment are:

```
darwin  amd64
darwin  arm64
freebsd amd64
freebsd arm64
linux   386
linux   amd64
linux   arm
linux   arm64
linux   ppc64le
linux   riscv64
linux   s390x
windows amd64
windows arm64
```

### Testing

PocketBase comes with mixed bag of unit and integration tests.
To run them, use the standard `go test` command:

```sh
go test ./...
```

Check also the [Testing guide](http://pocketbase.io/docs/testing) to learn how to write your own custom application tests.

## Security

If you discover a security vulnerability within PocketBase, please send an e-mail to **support at pocketbase.io**.

All reports will be promptly addressed and you'll be credited in the fix release notes.

## Contributing

PocketBase is free and open source project licensed under the [MIT License](LICENSE.md).
You are free to do whatever you want with it, even offering it as a paid service.

### Enhanced Features Development

This fork focuses on extending PocketBase with WhatsApp Business integration and enhanced OTP capabilities. We welcome contributions for:

- **WhatsApp Business API improvements** - Better error handling, rate limiting, template management
- **Multi-channel OTP enhancements** - SMS, Telegram, or other messaging platforms
- **Phone number validation** - International phone number formatting and validation
- **Testing and documentation** - Comprehensive tests for WhatsApp integration
- **UI/UX improvements** - Better admin interface for OTP configuration

### How to Contribute

You could help continuing its development by:

- [Contribute to the source code](CONTRIBUTING.md)
- [Suggest new features and report issues](https://github.com/pocketbase/pocketbase/issues)
- **Test WhatsApp integration** - Help us test the WhatsApp OTP functionality
- **Documentation improvements** - Help improve setup guides and API documentation

### Original PocketBase Contributing

For the original PocketBase project:

PRs for new OAuth2 providers, bug fixes, code optimizations and documentation improvements are more than welcome.

But please refrain creating PRs for _new features_ without previously discussing the implementation details.
PocketBase has a [roadmap](https://github.com/orgs/pocketbase/projects/2) and I try to work on issues in specific order and such PRs often come in out of nowhere and skew all initial planning with tedious back-and-forth communication.

Don't get upset if I close your PR, even if it is well executed and tested. This doesn't mean that it will never be merged.
Later we can always refer to it and/or take pieces of your implementation when the time comes to work on the issue (don't worry you'll be credited in the release notes).
