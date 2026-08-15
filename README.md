# 动物园亲子游线路

这是一个单一 Go module 的纯 backend 示例。它用固定内存 fixture 表示园区，家长可以指定入口、动物区、午餐点、出口、出发分钟和是否使用儿童车，系统会依据距离、坡度、拥挤程度、儿童车友好度及表演时间生成半日路线。`AdminService` 提供节点、连接和拥挤程度维护能力。

## 环境

- Go 1.22.12
- `CGO_ENABLED=0`
- 无数据库、网络服务、时钟或随机数依赖

## 运行

```bash
CGO_ENABLED=0 go run ./cmd/zoo-route
```

列出固定 fixture 的节点：

```bash
CGO_ENABLED=0 go run ./cmd/zoo-route -list
```

自定义路线：

```bash
CGO_ENABLED=0 go run ./cmd/zoo-route -entrance entrance-south -animals savanna,penguin -lunch family-cafe -exit exit-west -start 600 -stroller=false
```

业务链路测试命令为 `go test -count=1 ./...`。其中持久化失败传播用例固定验证调用方应看到落库错误，并用于复现当前注入缺陷。

## 代码结构

- `fixture.go`：固定节点、连线和表演时间
- `planner.go`：确定性的路径、顺序、等待和半日时长算法
- `service.go`：先校验、规划后落库的服务编排及管理员维护接口
- `memory_store.go`：可注入错误的内存存储
- `cmd/zoo-route`：可执行入口
