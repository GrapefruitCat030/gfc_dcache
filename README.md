# gfc_dcache 手搓一个golang分布式缓存

## 一、基于HTTP的内存缓存服务

我们现在要写一个基于HTTP的内存缓存服务。核心为【内存缓存结构】+ 【HTTP服务接口】。

### 1.1 内存缓存

我们将缓存封装为接口的形式，提供 get、set、delete三种方法。为什么用接口？我们当前要实现的是基于内存的缓存，后面我们还会有更多种实现。

### 1.2 HTTP服务

注意，我们要实现的是一个【服务】，那么我们就按照一般服务的实现形式来做项目布局。

这一节其实也不难，我们将缓存结构封装好，然后在此之上加入HTTP/REST接口的形式，即可将其我们的缓存升级为一个网络服务。我们项目的关键不是搓HTTP框架的轮子，不需要花太多时间在HTTP处理上，所以在golang原生http库之外用 gorilla/mux 库即可。

当前项目总体结构：

```
├── api
│   ├── handler
│   │   ├── delete.go
│   │   ├── get.go
│   │   ├── set.go
│   │   └── status.go
│   └── route
│       └── router.go
├── go.mod
├── go.sum
├── main.go
├── pkg
│   ├── cache
│   │   ├── cache.go
│   │   ├── mcache.go
│   │   └── stat.go
│   └── server
│       └── server.go
└── README.md
```









































































































## 二、改用 HTTP/REST 协议和 TCP 协议混合接口



## 三、RocksDB：缓存持久化与突破内存限制



## 四、缓存性能提升



## 五、分布式实现

 




