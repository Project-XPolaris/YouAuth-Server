# YouAuth 配置说明

YouAuth 支持通过配置文件和环境变量两种方式进行配置。环境变量的优先级高于配置文件。

## 配置文件

配置文件路径可以通过环境变量 `YOUAUTH_CONFIG_PATH` 指定。如果未指定，将使用默认配置文件。

## 配置项说明

### 基础配置

| 配置项 | 环境变量 | 类型 | 默认值 | 说明 |
|--------|----------|------|--------|------|
| addr | - | string | ":8000" | 服务监听地址 |
| application | - | string | "You Auth Service" | 应用名称 |
| instance | - | string | "main" | 实例名称 |

### JWT 配置

| 配置项 | 环境变量 | 类型 | 说明 |
|--------|----------|------|------|
| token.secret | YOUAUTH_TOKEN_SECRET | string | JWT 签名密钥 |
| token.issuer | YOUAUTH_TOKEN_ISSUER | string | JWT 签发者 |
| token.accessTokenExpiresIn | YOUAUTH_TOKEN_ACCESS_EXPIRES | int64 | 访问令牌过期时间（秒） |
| token.refreshTokenExpiresIn | YOUAUTH_TOKEN_REFRESH_EXPIRES | int64 | 刷新令牌过期时间（秒） |
| token.authCodeExpiresIn | YOUAUTH_TOKEN_AUTH_CODE_EXPIRES | int64 | 授权码过期时间（秒） |
| token.appTokenExpiresIn | YOUAUTH_TOKEN_APP_EXPIRES | int64 | 应用令牌过期时间（秒） |
| token.url | YOUAUTH_TOKEN_URL | string | JWT URL |

### 外部登录配置

| 配置项 | 环境变量 | 类型 | 说明 |
|--------|----------|------|------|
| externalLoginPage | YOUAUTH_EXTERNAL_LOGIN_PAGE | string | 外部登录页面 URL |

## 配置文件示例

```yaml
addr: ":8000"
application: "You Auth Service"
instance: "main"

token:
  secret: "your-secret-key"
  issuer: "youauth"
  accessTokenExpiresIn: 3600
  refreshTokenExpiresIn: 604800
  authCodeExpiresIn: 600
  appTokenExpiresIn: 31536000
  url: "https://auth.example.com"

externalLoginPage: "https://login.example.com"
```

## 环境变量示例

```bash
# 基础配置
export YOUAUTH_CONFIG_PATH="/path/to/config.yaml"

# JWT 配置
export YOUAUTH_TOKEN_SECRET="your-secret-key"
export YOUAUTH_TOKEN_ISSUER="youauth"
export YOUAUTH_TOKEN_ACCESS_EXPIRES="3600"
export YOUAUTH_TOKEN_REFRESH_EXPIRES="604800"
export YOUAUTH_TOKEN_AUTH_CODE_EXPIRES="600"
export YOUAUTH_TOKEN_APP_EXPIRES="31536000"
export YOUAUTH_TOKEN_URL="https://auth.example.com"

# 外部登录配置
export YOUAUTH_EXTERNAL_LOGIN_PAGE="https://login.example.com"
```

## 注意事项

1. 所有时间相关的配置项（如过期时间）均以秒为单位
2. 环境变量的优先级高于配置文件中的值
3. 如果环境变量未设置，将使用配置文件中的值
4. 如果配置文件和环境变量都未设置，将使用默认值 