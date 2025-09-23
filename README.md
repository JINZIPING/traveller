# traveller
网络质量评估All in One

## 项目结构

```
.
├── client/               # 探测客户端
│   ├── config/           # Client 配置文件
│   ├── internal/         
│   │   ├── infra/        # RabbitMQ & Logger 初始化
│   │   ├── probe/        # ICMP/TCP 探测实现
│   │   └── service/      # 任务消费逻辑
│   ├── logs/             # Client 日志
│   └── main.go           # Client 入口
│
├── server/               # 服务端
│   ├── config/           # Server 配置 (MySQL, ClickHouse, RabbitMQ, Prometheus)
│   ├── internal/
│   │   ├── adapter/      # 适配器 (Prometheus Publisher)
│   │   ├── dao/          # 数据访问层 (MySQL / ClickHouse)
│   │   ├── infra/        # 基础设施 (DB/MQ/Logger 初始化)
│   │   ├── router/       # 路由注册
│   │   └── service/      # 任务下发 & 结果消费逻辑
│   ├── logs/             # Server 日志
│   └── main.go           # Server 入口
│
├── pkg/                  # 公共库 (Client/Server 共享)
│   ├── model/            # 任务 & 结果数据模型
│   └── mq/               # MQ 抽象接口 & RabbitMQ 实现
│
├── config/               # 全局配置 (YAML)
├── go.mod
├── go.sum
└── README.md
```
## 架构设计与工厂模式的应用
### 公共模块pkg
在项目演进过程中，Server 与 Client 都需要使用相同的数据模型（Task、Result）以及相同的消息队列交互逻辑（Producer、Consumer）。
为了避免重复代码并提升可维护性，公共部分被抽取到了`pkg/`中：

- `pkg/model/`：定义探测任务（Task）与探测结果（Result）的数据结构，供 Server 与 Client 共享。
- `pkg/mq/`：抽象消息队列的通用接口（生产者/消费者），并提供 RabbitMQ 的具体实现。

这样，Server 和 Client 不再直接依赖 RabbitMQ 的细节，而是通过统一的接口来交互。

### 抽象工厂模式在 MQ 中的应用
在消息队列部分，项目引入了 抽象工厂模式（Abstract Factory Pattern）。
因为在 Server 和 Client 中，都存在 Producer 和 Consumer，但它们的职责不同：

- Server： 
  - Producer：下发探测任务 
  - Consumer：消费探测结果
- Client： 
  - Producer：上报探测结果 
  - Consumer：消费探测任务

为了避免重复代码，并且支持未来可能的扩展（例如替换 RabbitMQ到Kafka），抽象出了一组统一接口。

### 设计思路
在`pkg/mq`中定义了三个核心接口：
```go
// 任务发布器
type TaskPublisher interface {
    Publish(task any) error
}

// 结果消费者
type ResultConsumer interface {
    Consume() (<-chan []byte, error)
}

// 抽象工厂
type Factory interface {
    CreatePublisher() TaskPublisher
    CreateConsumer() ResultConsumer
}
```
然后在`pkg/mq/rabbitmq/`下实现了 RabbitMQ 版本的工厂：
```go
type RabbitMQFactory struct { ... }

func (f *RabbitMQFactory) CreatePublisher() mq.TaskPublisher { ... }
func (f *RabbitMQFactory) CreateConsumer() mq.ResultConsumer { ... }
```
Server 和 Client 通过调用`factory.CreatePublisher()`/`factory.CreateConsumer()`来获取所需的 MQ 组件，不需要关心底层是 RabbitMQ 还是未来可能的 Kafka

## 任务调度器（Scheduler）

本项目内置了一个轻量级的任务调度器，用于管理网络探测任务（TCP/ICMP）。

* **核心功能**

    * 支持 **一次性任务** 和 **周期性任务**（通过 `interval_sec` 控制执行间隔）。
    * 使用 `map[string]*tickJob` 管理任务，每个任务都有唯一 `jobID`（如 `icmp:8.8.8.8`、`tcp:1.1.1.1:80`）。
    * 基于 `time.Ticker` 周期触发，任务逻辑与调度解耦。
    * 提供 `Remove(jobID)` 方法，可优雅地停止并清理任务。

* **设计优势**

    * 无需依赖第三方调度服务，避免外部组件不可用时造成数据中断。
    * 内置在 Server 中，保证探测任务 **持续性和可控性**。
    * 通过 API 可灵活下发、暂停和删除任务，便于统一管理。


