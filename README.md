<!-- Source: Best-README-Template BLANK_README (Unlicense) — https://github.com/othneildrew/Best-README-Template -->
<a id="readme-top"></a>

# UnifyAuthCenter

A self-hosted Go authentication center that issues TOTP-verified sessions for end users, gates an admin panel and an internal status API behind shared passwords, and includes a tool for generating new TOTP secrets with QR codes.

**English** · [简体中文](README.zh-CN.md)

[![CI](https://github.com/anyingiit/UnifyAuthCenter/actions/workflows/ci.yml/badge.svg)](https://github.com/anyingiit/UnifyAuthCenter/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/anyingiit/UnifyAuthCenter)](LICENSE)

[Report a bug](https://github.com/anyingiit/UnifyAuthCenter/issues/new?template=bug_report.yml) · [Request a feature](https://github.com/anyingiit/UnifyAuthCenter/issues/new?template=feature_request.yml)

<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#getting-started">Getting Started</a></li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

## About The Project

UnifyAuthCenter (`main.go`) is a small standalone HTTP service, built to sit in front of a handful of other sites and answer one question: is this visitor logged in? An end user authenticates at `/auth_center/login` with a TOTP code; on success the server creates a database-backed session and sets it as a cookie (`handles/authCenter/handle.go`). A separate `/admin` area, guarded by its own fixed password, lets an operator list and revoke any of those sessions (`handles/admin/login/handle.go`). Other backends can check whether a session is still valid by calling `/internal/authorization_user_auth_status` with a shared `Internal-Password` header (`handles/internalApi/userAuthStatus.go`), which is how "unify" is meant: one service, several call sites.

A second, unrelated feature lives under `/tool/generation_TOTP`: a page that generates a brand-new TOTP secret and renders its QR code so a person can enroll a fresh authenticator (`handles/tool/generationTOTP/handle.go`).

See the [open issues](https://github.com/anyingiit/UnifyAuthCenter/issues) for planned features and known issues.

## Getting Started

### Prerequisites

- Go 1.18 or newer (the floor declared in `go.mod`)
- A C toolchain for CGO, since `gorm.io/driver/sqlite` builds on `github.com/mattn/go-sqlite3` -- or Docker, to skip installing Go and a C compiler locally and build the container image instead

### Installation

```sh
git clone https://github.com/anyingiit/UnifyAuthCenter.git
cd UnifyAuthCenter
go build -v -o app ./main.go
```

To build the Docker image instead, use the repository's own build script:

```sh
./scripts/build.sh
```

## Usage

Run the binary from the repository root, so it can find `database/sessions.db`, `template/` and `static/` by their relative paths:

```sh
./app
```

The server listens on `0.0.0.0:8066` and creates `database/sessions.db` on first start if it does not already exist. `scripts/run.sh` shows the equivalent for the Docker image, publishing the same port:

```sh
docker run -d -p 8066:8066 anyingiit/unify_auth_center:V0.6
```

## Contributing

Contributions are welcome. Read [CONTRIBUTING.md](CONTRIBUTING.md) for how to open an issue or a pull request, and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for the standards expected of everyone taking part.

Please do not report security issues in public issues or pull requests. [SECURITY.md](SECURITY.md) explains how to report them privately.

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for details.

## Contact

Project link: [https://github.com/anyingiit/UnifyAuthCenter](https://github.com/anyingiit/UnifyAuthCenter)

<p align="right">(<a href="#readme-top">back to top</a>)</p>
