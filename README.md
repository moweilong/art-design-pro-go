# Project art-design-pro-go

art-design-pro-go 是一个基于 Go 语言开发的现代化微服务应用，采用简洁架构设计，具有代码质量高、扩展能力强、符合 Go 编码及最佳实践等特点。

art-design-pro-go 具有以下特性：
- 软件架构：采用简洁架构设计，确保项目结构清晰、易维护；
- 高频 Go 包：使用了 Go 项目开发中常用的包，如 gin、otel、gorm、gin、uuid、cobra、viper、pflag、resty、govalidator、slog、protobuf、casbin、onexstack 等；
- 目录结构：遵循 [project-layout](https://github.com/golang-standards/project-layout) 规范，采用标准化的目录结构；
- 认证与授权：实现了基于 JWT 的认证和基于 Casbin 的授权功能；
- 错误处理：设计了独立的错误包及错误码管理机制；
- 构建与管理：使用高质量的 Makefile 对项目进行管理；
- 代码质量：通过 golangci-lint 工具对代码进行静态检查，确保代码质量；
- 测试覆盖：包含单元测试、性能测试、模糊测试和示例测试等多种测试案例；
- 丰富的 Web 功能：支持 Trace ID、优雅关停、中间件、跨域处理、异常恢复等功能；
- 多种数据交换格式：支持 JSON 和 Protobuf 数据格式的交换；
- 开发规范：遵循多种开发规范，包括代码规范、版本规范、接口规范、日志规范、错误规范以及提交规范等；
- API 设计：接口设计遵循 RESTful API 规范；
- 项目具有 Dockerfile，并且 Dockerfile 符合最佳实践；

## Getting Started

### Prerequisites

在开始之前，请确保您的开发环境中安装了以下工具：

**必需工具：**
- [Go](https://golang.org/dl/) 1.25.3 或更高版本
- [Git](https://git-scm.com/) 版本控制工具
- [Make](https://www.gnu.org/software/make/) 构建工具

**可选工具：**
- [Docker](https://www.docker.com/) 容器化部署
- [golangci-lint](https://golangci-lint.run/) 代码静态检查

**验证安装：**

```bash
$ go version  
go version go1.25.3 linux/amd64  
$ make --version  
GNU Make 4.3  
```

### Building

> 提示：项目配置文件配置项 `metadata.makefileMode` 不能为 `none`，如果为 `none` 需要自行构建。

在项目根目录下，执行以下命令构建项目：

**1. 安装依赖工具和包**

```bash
$ make deps  # 安装项目所需的开发工具  
$ go mod tidy # 下载 Go 模块依赖  
```

**2. 生成代码**

```bash
$ make protoc # generate gRPC code  
$ go get cloud.google.com/go/compute@latest cloud.google.com/go/compute/metadata@latest  
$ go mod tidy # tidy dependencies  
$ go generate ./... # run all go:generate directives  
```

**3. 构建应用**

```bash
$ make build # build all binary files locate in cmd/  
```

**构建结果：**

```bash
_output/platforms/  
├── linux/  
│   └── amd64/  
│       └── art-apiserver  # apiserver 服务二进制文件  
└── darwin/  
    └── amd64/  
        └── art-apiserver  
```

### Running

启动服务有多种方式：

**1. 使用构建的二进制文件运行**

```bash  
# 启动 apiserver 服务  
$ _output/platforms/linux/amd64/art-apiserver --config configs/art-apiserver.yaml  
# 服务将在以下端口启动：  
# - HTTP API: http://localhost:5555
# - Health Check: http://localhost:5555/healthz  
# - Metrics: http://localhost:5555/metrics  
$ curl http://localhost:5555/healthz # 测试：打开另外一个终端，调用健康检查接口  
```

**2. 使用 Docker 运行**

```bash
# 构建镜像  
$ make image
$ docker run --name art-apiserver -v configs/art-apiserver.yaml:/etc/art-apiserver.yaml -p 5555:5555 docker.io/moweilong/art-apiserver:latest -c /etc/art-apiserver.yaml
```

**配置文件示例：**  

art-apiserver 配置文件 `configs/art-apiserver.yaml`：

```yaml
addr: 0.0.0.0:5555 # 服务监听地址
timeout: 30s # 服务端超时
otel:
  endpoint: 127.0.0.1:4327
  service-name: art-apiserver
  output-mode: otel
  level: debug
  add-source: true
  use-prometheus-endpoint: true
  slog: # 改配置项只有 output-mod 为 slog 时生效
    format: text
    time-format: "2006-01-02 15:04:05"
    output: stdout
mysql:  
  addr: 127.0.0.1:3306
  username: onex
  password: "onex(#)666"
  database: onex
  max-connection-life-time: 10s
  max-idle-connections: 100
  max-open-connections: 100
```  

## Versioning

本项目遵循 [语义版本控制](https://semver.org/lang/zh-CN/) 规范。

## Authors

### 主要贡献者

- **moweilong** - *项目创建者和维护者* - [kalandramo@gmail.com](mailto:kalandramo@gmail.com)
  - 项目架构设计
  - 核心功能开发
  - 技术方案制定

### 贡献者列表

感谢所有为本项目做出贡献的开发者们！

<!-- 这里会自动显示贡献者头像 -->
<a href="github.com/moweilong/art-design-pro-go/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=github.com/moweilong/art-design-pro-go" />
</a>

*贡献者列表由 [contrib.rocks](https://contrib.rocks) 生成*

## 附录

### 项目结构

```bash
art-design-pro-go/  
├── cmd/                     # 应用程序入口  
│   └── art-apiserver/       # apiserver 服务  
│       └── main.go          # 主函数  
├── internal/                # 私有应用程序代码  
│   └── apiserver/             # apiserver 内部包  
│       ├── biz/             # 业务逻辑层  
│       ├── handler/         # gin 处理器  
│       ├── model/           # GORM 数据模型  
│       ├── pkg/             # 内部工具包  
│       └── store/           # 数据访问层  
├── pkg/                     # 公共库代码  
│   ├── api/                 # API 定义  
├── examples/                # 示例代码  
│   └── client/              # 客户端示例  
├── configs/                 # 配置文件  
├── docs/                    # 项目文档  
├── build/                   # 构建配置  
│   └── docker/              # Docker 文件  
├── scripts/                 # 构建和部署脚本  
├── third_party/             # 第三方依赖  
├── Makefile                 # 构建配置  
├── go.mod                   # Go 模块文件  
├── go.sum                   # Go 模块校验文件  
└── README.md                # 项目说明文档  
```

### 相关链接

- [项目文档](docs/)
- [问题追踪](github.com/moweilong/art-design-pro-go/issues)
- [讨论区](github.com/moweilong/art-design-pro-go/discussions)
- [项目看板](github.com/moweilong/art-design-pro-go/projects)
- [发布页面](github.com/moweilong/art-design-pro-go/releases)

### 支持

如果这个项目对您有帮助，请考虑给我们一个 ⭐️ 来支持项目发展！

[![Star History Chart](https://api.star-history.com/svg?repos=github.com/moweilong/art-design-pro-go&type=Date)](https://star-history.com/#github.com/moweilong/art-design-pro-go&Date)
