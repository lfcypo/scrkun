# scrkun

`scrkun` 是一个基于屏幕截图和视觉语言模型的桌面活动监测程序。它会定时检查屏幕画面，识别当前活动，将检测结果保存到 SQLite，并在检测到视频、游戏或小说相关活动时发送通知。项目也提供用于查看记录的网页和 JSON API。

## 功能

- 定时截取主显示器画面；与上次截图相似度较高时跳过本轮分析
- 将截图发送给阿里云百炼 DashScope 的 `qwen3-vl-flash` 模型，识别学习、编程、视频、游戏、文档、网页、小说或闲置等活动
- 按置信度阈值筛选识别结果，并保存活动类型、判断理由、权重和事件时间
- 对视频、游戏和小说活动发送 Bark 或「虾推啥」通知
- 提供活动记录网页及 `/api/data` JSON 接口

## 环境要求

- Go `1.24.6` 或兼容的较新版本
- 运行程序的设备需要有可用的图形桌面和显示器
- DashScope API Key；启用对应通知时还需要 Bark Token 或「虾推啥」Token

## 配置

程序从当前工作目录读取 `config.toml`。复制示例配置后，按需填写密钥：

```sh
cp config-example.toml config.toml
```

示例配置：

```toml
[ai]
apikey = "你的 DashScope API Key"

[bark]
token = "你的 Bark Token"

[xtuis]
token = "你的虾推啥 Token"

[machine]
name = "设备名称"

[monitor]
sleep = 5
interval = 300

[detect]
threshold = 0.5

[server]
port = 8890

[db]
path = "data.db"

[log]
mode = 3
file = "./logs/scrkun.log"
maxSize = 100
maxBackups = 10
maxAge = 30
compress = true
```

`bark.token` 和 `xtuis.token` 可按需配置；留空时对应通知渠道不会启用。`monitor.sleep` 是检查循环的间隔秒数，`monitor.interval` 是发出异常通知后的冷却秒数，`detect.threshold` 是结果置信度阈值。`machine.name` 未设置时使用系统主机名。日志模式 `1` 表示标准输出、`2` 表示文件、`3` 表示同时输出到两处。

首次启动前请确认 `ai.apikey` 已设置为可调用 DashScope 兼容模式 API 的有效密钥。截图会以 JPEG 图片形式发送到该服务进行分析。

## 运行

直接使用 Go 启动：

```sh
go run .
```

或编译并运行：

```sh
make build
./build/scrkun
```

`make build` 会在 `build/` 下生成当前系统和架构对应的二进制文件。启动时请从包含 `config.toml` 的目录运行程序。

仅启动记录查询网页、不进行截图监测：

```sh
./build/scrkun -q
```

## 查看记录

默认网页地址为 [http://localhost:8890](http://localhost:8890)，端口可通过 `server.port` 修改。

记录接口：

```http
GET /api/data
```

接口返回所有活动记录，结构为 `{"data": [...]}`。每条记录包含事件时间、事件类型、活动类型、判断理由和权重等字段。`eventType` 中 `0` 表示普通事件、`1` 表示异常事件；`usageType` 按照 `model/usage_type.go` 中的枚举顺序编码。

## 通知配置

- **Bark**：在 `bark.token` 中填写 Bark Token，通知会通过 Bark 推送。
- **虾推啥**：关注「虾推啥」微信公众号，复制收到的 Token 并填入 `xtuis.token`，通知会发送到绑定的微信。

两个渠道可以同时启用。程序仅在活动类型为视频、游戏或小说时发送通知。

## 开发

格式化 Go 代码：

```sh
make fmt
```

构建当前平台：

```sh
make build
```

构建 Makefile 中列出的 Linux、macOS 和 Windows 目标：

```sh
make build-all
```

## 项目结构

```text
capture/   屏幕截图、相似度比较与图像编码
config/    TOML 配置加载
model/     活动记录及类型定义
notifiy/   Bark 和虾推啥通知
server/    查询网页与 HTTP API
store/     SQLite 数据库
usagedt/   截图活动识别及 DashScope 请求
```

## 许可证

本项目采用 [MIT License](LICENSE)。
