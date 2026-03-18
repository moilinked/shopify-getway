## shopify-gateway

一个基于 Go 的 Shopify 网关服务，用于承接应用侧请求校验、App Proxy 鉴权以及 Shopify webhook 监听等能力。

### Doc

- [shopify-session-tokens](https://shopify.dev/docs/apps/build/authentication-authorization/session-tokens)
- [shopify-webhook](https://shopify.dev/docs/apps/build/webhooks/subscribe/https#step-2-validate-the-origin-of-your-webhook-to-ensure-its-coming-from-shopify)

### 技术栈

- `Go 1.25`：项目主语言与运行时
- `chi v5`：HTTP 路由与中间件组织
- `golang-jwt/jwt v5`：Shopify Session Token(JWT) 校验
- `yaml.v3`：读取 `config/config.yaml` 配置
- `zerolog`：结构化日志输出
- `Docker` / `docker-compose`：容器化构建与本地部署

### 项目结构

- `cmd/server`：服务启动入口，负责加载配置、初始化日志并启动 HTTP Server
- `config`：项目配置文件
- `internal/config`：配置加载与字段校验
- `internal/httpapi`：路由注册、健康检查与 webhook 接收 handler
- `internal/logger`：日志初始化与日志落盘相关能力
- `internal/middleware`：CORS、请求日志、Shopify Session Token、App Proxy 签名、Webhook 验签及上下文处理中间件
