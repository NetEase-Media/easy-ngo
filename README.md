# 简介
[文档地址](https://netease-media.github.io/easy-ngo-website/)

[快速使用例子](https://github.com/NetEase-Media/easy-ngo-examples)

## 什么是 `easy-ngo`
easy-ngo是由网易传媒开发的基于Go语言的开发工具包，基于easy-ngo工具包，开发者可以快速构建高可用、大并发的应用。easy-ngo已经在传媒内部大量使用，节省了业务开发的大量时间。
使用easy-ngo带来的好处：
* 开箱即用，用户基于easy-ngo直接开发代码，无需进行SDK的选型及适配工作
* 在第三方SDK的基础上支持了Metrics、Tracer、Logger，用户只需要配置即可实现可观测性
* 支持文件、ENV、Param、Apollo等配置源，用户可直接使用
* 支持gorm、redis、kafka、memcache、zk等大部分的中间件SDK，并在这些SDK的基础上封装了Metrics、Tracer、Logger等可观测性的特性，用户可以直接使用，无需再进行封装
* 支持gin、grpc等常用的server，并在这些server的基础上封装了Metrics、Tracer、Logger等可观测性的特性，用户可以直接使用，无需再进行封装
* 用户可以快速构建微服务架构，easy-ngo框架支持服务注册、服务发现、负载均衡等特性
* easy-ngo基于插件和hook机制，用户可以很方便的扩展自己想要的插件，并注册到整个框架中

## easy-ngo框架建设背景
网易传媒在2020年底开始尝试使用Go语言做业务开发，目前大部分业务都基于Go语言开发，并在线上为用户提供服务

### 背景
网易传媒的主要开发语言是Java，在业务全部接入容器后，在线业务也面临着以下一些问题：
1. 在线业务内存使用量偏高：传媒主要开发语言是Java，使用SpringBoot框架，最少使用2G内存，普遍内存使用量都在4G以上，还有8G、16G、32G等内存使用的应用。
2. 在线业务编译速度和启动速度偏慢：使用maven编译、打包、打镜像、传镜像都比较耗时，拖慢了整个CI的流程。
3. 占用空间较大：由于使用Java，JVM在镜像实例都需要上百兆（400M以上）的空间，拉取，上传都比较耗时。
网易传媒于2020年将核心业务全部迁入容器，在容器和微服务的大背景下，应用的小而快显得就格外的重要，Go语言就比较适合于我们的需求，目前已经有很多互联网厂商都在积极推进Go语言的应用，于是，网易传媒在2020年底开始尝试Go语言的探索，并在2021年使用Go语言重构核心业务，目前大部分业务都基于Go语言开发，并在线上为用户提供服务。
### Go语言简介
Go语言于2009年11月正式宣布推出，它是Google开发的一种静态强类型、编译型、并发型、并具有垃圾回收功能的编程语言，它的特性包括：

* 编译速度快
* 语法简单
* 像动态语言一样开发
* 资源消耗少
* 为并发IO而生
* 可运维性好
* 与C/C++兼容
* 统一而完备的工具集
### easy-ngo介绍
如果我们要开发一个应用，除了应用核心业务代码外，还需要很多的底层支持，可以参见如下图：
![应用的依赖](https://netease-media.github.io/easy-ngo-website/assets/images/easy-ngo-1-366cb15746dd0d0d2e1dc2ffcb023845.png)

在传媒技术团队中推广Go语言，亟需一个Web框架提供给业务开发同事使用，内含业务开发常用库，避免重复造轮子影响效率，并且需要无感知的自动监控数据上报，能在框架层面支持业务的优雅上下线，对云原生监控的支持，支持服务注册，服务发现，服务调用等能力，于是就孕育出easy-ngo框架。
easy-ngo的主要目标如下：
* 提供比原有Java框架更高的性能和更低的资源占用率
* 尽量为业务开发者提供所需的全部工具库
* 嵌入云原生监控，自动上传监控数据
* 嵌入全链路监控，提供标准的opentracering协议，和第三方的全链路监控系统结合(Jaeger、Zipkin)
* 自动加载配置和初始化程序环境，开发者能直接使用各种库
* 与线上的健康检查、运维接口等运行环境匹配，无需用户手动开发配置
* 能够方便的扩展第三方的SDK到框架中

easy-ngo避免重复造轮子，所有模块都是在多个开源库中对比并挑选其一，然后增加部分必需功能，easy-ngo支持的能力如下图所示：
![easy-ngo的能力](https://netease-media.github.io/easy-ngo-website/assets/images/easy-ngo-2-705e4cec580d238bc19bb24b628aa539.png)

easy-ngo框架为业务选择并包装了用到的中间件和基础服务，让业务可以快速的进入到业务开发的阶段，省去了研究和比较一些基础组件的时间，大大节省了业务的开发周期。


## 快速开始
让我们从一个最简单的HelloWorld开始easy-ngo之旅吧
首先，将代码从github上clone下来
```
git clone https://github.com/NetEase-Media/easy-ngo.git
```
其次，进入sample目录
```
cd examples/application
```
查看代码如下：
```
package main

import (
	"net/http"

	"github.com/NetEase-Media/easy-ngo/application"
	_ "github.com/NetEase-Media/easy-ngo/application/r/rconfig"
	"github.com/NetEase-Media/easy-ngo/application/r/rgin"
	_ "github.com/NetEase-Media/easy-ngo/examples/application/include"
	"github.com/gin-gonic/gin"
)

func main() {
	app := application.Default()
	app.Initialize(xgin)
	app.Startup()
}

func xgin() error {
	g := rgin.Gin()
	g.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world!")
	})
	return nil
}
```
查看配置如下：
```
[ngo.app]
name = "quickstart-demo"
[ngo.server.gin]
port = 8888
enabledMetric = false
[ngo.app.healthz]
port = 10000
```
运行以下命令，最简单的服务便启动了
```
go run . -c ./app.toml
```
So Cool！更多示例，我们可以进入examples目录查看。
easy-ngo访问地址如下
```
https://github.com/NetEase-Media/easy-ngo-examples
```

## 微信交流群
欢迎大家扫描二维码加入我们
![微信群](https://github.com/NetEase-Media/easy-ngo-website/blob/gh-pages/images/Wechateasyngo.jpeg)
