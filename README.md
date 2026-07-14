# Study_Traker (学习检索与进度管理助手)

基于 Go 原生标准库实现的双轨制学习进度管理系统，同时支持面向人类用户的 HTML Dashboard 和面向 AI Agent 的 OpenAPI 接口。

## 核心功能

* **双轨制路由**：
  * `/dashboard` (GET) - 网页端 HTML 仪表盘，直观展示学习历史。
  * `/api/logs` (GET/POST) - 原始学习日志的查询与提交（支持单条 JSON 写入）。
  * `/api/skill` (GET) - 面向大模型 Agent 的数据摘要接口，返回总学时、学科列表等去重预计算数据。
* **SRE 级加固**：
  * **原子写入**：写文件采用 `.tmp` 临时文件 + `os.Rename` 原子替换，防止中途崩溃导致数据损坏。
  * **并发安全**：使用 `sync.RWMutex` 进行读写锁隔离保护，缩小锁粒度，防止死锁。
  * **输入校验与防御**：使用 `http.MaxBytesReader` 限制请求体大小，并防范多余 JSON 注入。
  * **安全认证**：实现 Bearer Token 拦截中间件，防止接口在内网无防护暴露。
  * **优雅关闭**：监听 `SIGINT` 和 `SIGTERM` 信号，通过 `server.Shutdown` 保证在退出时处理完现有连接。

## 项目结构

```text
Study_Traker/
├── main.go          # 服务启动入口、配置与优雅关闭
├── handler.go       # 协议传输层，包含所有路由 HTTP Handler 与认证中间件
├── storage.go       # 数据存储与加工层，负责读写锁控制、原子写入与数据过滤
├── model.go         # 核心领域数据模型定义
├── validate.go      # 业务级输入合法性校验
├── openapi.yaml     # 供大模型（Dify）导入的 OpenAPI 规范说明书
└── templates/
    └── dashboard.tmpl # 网页端模板文件
```

## 运行与部署

### 1. 启动 Go 服务

系统读取 `Secret_Token_Key` 环境变量作为 API 鉴权密钥。启动时必须监听 `0.0.0.0` 以供容器网络访问：

```bash
export Secret_Token_Key="your-secret-key"
go run . -addr 0.0.0.0:8082
```

### 2. 本地接口验证

```bash
# 查询日志 (需带 Token)
curl -i -H "Authorization: Bearer your-secret-key" http://127.0.0.1:8082/api/logs

# 条件查询
curl -i -H "Authorization: Bearer your-secret-key" "http://127.0.0.1:8082/api/logs?subject=高数"

# 新增日志
curl -i -X POST -H "Content-Type: application/json" -H "Authorization: Bearer your-secret-key" \
  -d '{"date":"2026-07-13","subject":"高数","theme":"积分","duration":90,"challenge_level":3}' \
  http://127.0.0.1:8082/api/logs
```

## 接入大模型平台 (Dify)

本玩具设计之初即为 LLM Agent 的自定义工具。以 Dify 平台为例：

1. **导入规格**：
   在 Dify 的“工具/插件”页面选择 **Swagger API**，导入项目根目录下的 `openapi.yaml` 文件。
2. **注意网络地址**：
   YAML 顶部的 `servers[0].url` 必须指向 Dify 容器可达的宿主机 IP（在 Docker 网桥中通常为 `http://172.17.0.1:8082`）。
3. **安全配置**：
   在 Dify 的鉴权设置中选择 **Bearer** 模式，填入你启动 Go 服务时设置 of `Secret_Token_Key` 的值。
4. **SSRF 防护加固**：
   如果 Dify 提示限制对本地 API 的访问，需在 Dify 的 `.env` 配置文件中将宿主机 IP 写入白名单：
   `SSRF_PROXY_ALLOW_PRIVATE_IPS=172.17.0.1`。
5. **创建 Agent**：
   在 Dify 中创建 **Agent 助手**，绑定你的中转模型，添加 `Study_Traker API` 即可完成对接。
