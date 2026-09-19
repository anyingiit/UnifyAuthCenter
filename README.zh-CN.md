[English](README.md) · **简体中文**

> 英文版是规范版本。本页与 [README.md](README.md) 不一致时，以英文版为准。

<!-- translation-of: README.md sha256:94feecd76b787605 -->

<!-- Source: Best-README-Template BLANK_README (Unlicense) — https://github.com/othneildrew/Best-README-Template -->
<a id="readme-top"></a>

# UnifyAuthCenter

一个自托管的 Go 语言认证中心：为终端用户签发经 TOTP 验证的会话，用共享密码守护一个管理后台和一个内部状态查询接口，并附带一个生成新 TOTP 密钥及其二维码的工具。

[![CI](https://github.com/anyingiit/UnifyAuthCenter/actions/workflows/ci.yml/badge.svg)](https://github.com/anyingiit/UnifyAuthCenter/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/anyingiit/UnifyAuthCenter)](LICENSE)

[报告问题](https://github.com/anyingiit/UnifyAuthCenter/issues/new?template=bug_report.yml) · [提出需求](https://github.com/anyingiit/UnifyAuthCenter/issues/new?template=feature_request.yml)

<details>
  <summary>目录</summary>
  <ol>
    <li><a href="#about-the-project">关于本项目</a></li>
    <li><a href="#getting-started">开始使用</a></li>
    <li><a href="#usage">用法</a></li>
    <li><a href="#contributing">参与贡献</a></li>
    <li><a href="#license">许可证</a></li>
    <li><a href="#contact">联系方式</a></li>
  </ol>
</details>

## 关于本项目

UnifyAuthCenter（`main.go`）是一个独立的小型 HTTP 服务，用于架在若干个其他站点前面，回答同一个问题：这位访客是否已登录？终端用户在 `/auth_center/login` 用 TOTP 验证码登录，验证成功后服务端会创建一条数据库中的会话记录并以 Cookie 的形式下发（`handles/authCenter/handle.go`）。另有一个由固定密码守护的 `/admin` 区域，供管理员查看并注销任意一条会话（`handles/admin/login/handle.go`）。其他后端服务可以在请求头中携带约定好的 `Internal-Password`，调用 `/internal/authorization_user_auth_status` 来确认某个会话是否仍然有效（`handles/internalApi/userAuthStatus.go`）——这也是"统一"（Unify）二字的由来：一个服务，供多处调用。

第二个、与登录无关的功能位于 `/tool/generation_TOTP`：一个用于生成全新 TOTP 密钥并渲染其二维码的页面，方便用户注册一个新的身份验证器（`handles/tool/generationTOTP/handle.go`）。

计划中的功能与已知问题，见 [open issues](https://github.com/anyingiit/UnifyAuthCenter/issues)。

## 开始使用

### 环境要求

- Go 1.18 或更高版本（`go.mod` 中声明的最低版本）
- 用于 CGO 的 C 工具链，因为 `gorm.io/driver/sqlite` 依赖 `github.com/mattn/go-sqlite3`；也可以改用 Docker，省去在本地安装 Go 和 C 编译器，直接构建容器镜像

### 安装

```sh
git clone https://github.com/anyingiit/UnifyAuthCenter.git
cd UnifyAuthCenter
go build -v -o app ./main.go
```

若改用 Docker 构建镜像，可使用仓库自带的构建脚本：

```sh
./scripts/build.sh
```

## 用法

请从仓库根目录运行该可执行文件，这样它才能按相对路径找到 `database/sessions.db`、`template/` 和 `static/`：

```sh
./app
```

服务会监听 `0.0.0.0:8066`，并在首次启动时创建 `database/sessions.db`（如果该文件尚不存在）。`scripts/run.sh` 展示了 Docker 镜像下的等价用法，发布同一个端口：

```sh
docker run -d -p 8066:8066 anyingiit/unify_auth_center:V0.6
```

## 参与贡献

欢迎参与。[CONTRIBUTING.md](CONTRIBUTING.md) 说明如何提交 issue 或 pull request，[CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) 说明对所有参与者的行为要求。

请不要在公开的 issue 或 pull request 中报告安全问题。[SECURITY.md](SECURITY.md) 说明了私下报告的方式。

## 许可证

以 MIT 许可证分发。详见 [LICENSE](LICENSE)。

## 联系方式

项目地址：[https://github.com/anyingiit/UnifyAuthCenter](https://github.com/anyingiit/UnifyAuthCenter)

<p align="right">(<a href="#readme-top">back to top</a>)</p>
