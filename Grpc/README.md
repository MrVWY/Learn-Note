### grpc

```
A high-performance, open-source universal RPC framework
```
&ensp;&ensp;所谓RPC(remote procedure call 远程过程调用)框架实际是提供了一套机制，使得应用程序之间可以进行通信，而且也遵从server/client模型。使用的时候客户端调用server端提供的接口就像是调用本地的函数一样。

### 性能
&ensp;&ensp;gRPC 消息使用 Protobuf（一种高效的二进制消息格式）进行序列化。 Protobuf 在服务器和客户端上可以非常快速地序列化。 Protobuf 序列化产生的有效负载较小，这在移动应用等带宽有限的方案中很重要。

&ensp;&ensp;gRPC 专为 HTTP/2（HTTP 的主要版本）而设计，与 HTTP 1.x 相比，HTTP/2 具有巨大性能优势：

- 二进制组帧和压缩。 HTTP/2 协议在发送和接收方面均紧凑且高效。
- 在单个 TCP 连接上多路复用多个 HTTP/2 调用。 多路复用可消除队头阻塞。

&ensp;&ensp;HTTP/2 不是 gRPC 独占的。 许多请求类型（包括具有 JSON 的 HTTP API）都可以使用 HTTP/2，并受益于其性能改进。

### grpc(Protobuf)与http(json)对比
|    功能    |       gRPC        | 具有 JSON 的 HTTP API |
|:--------:|:-----------------:|:------------------:|
|    协定    |    必需 (.proto)    |    可选 (OpenAPI)    |
|    协议    |      HTTP/2       |        HTTP        |
| Payload  | 	Protobuf（小型，二进制） |  	JSON（大型，人工可读取）   |
|   规定性    |     	严格规范	宽松      |    任何 HTTP 均有效     |
|   流式处理   |   	客户端、服务器，双向	    |      客户端、服务器       |
|  浏览器支持   | 	无（需要 grpc-web）	  |         是          |
|   安全性	   |     传输 (TLS)	     |      传输 (TLS)      |
| 客户端代码生成	 |        是	         |  OpenAPI + 第三方工具   |


