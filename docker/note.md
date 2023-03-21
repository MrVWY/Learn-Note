### docker file 多个 FROM 指令的意义  
多阶段构建的最大意义:每一条 FROM 指令都是一个构建阶段，多条 FROM 表示多阶段构建，最后生成的镜像只能是最后一个阶段的结果，最大的使用场景是将编译环境和运行环境分离，
因此将前置阶段中的编译好的文件拷贝到后边的阶段中，能使得最后生成的镜像里面不包含编译时所需要的文件，减少镜像体积。  

```dockerfile
# 编译阶段
FROM golang:1.10.3

COPY server.go /build/

WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=6 go build -ldflags '-w -s' -o server

# 运行阶段
FROM scratch

# 从编译阶段的中拷贝编译结果到当前镜像中 --from=0 参数，从前边的阶段中拷贝文件到当前阶段中，多个FROM语句时，0代表第一个阶段
COPY --from=0 /build/server /

ENTRYPOINT ["/server"]
```

```dockerfile
# 编译阶段 命名为 builder
FROM golang:1.10.3 as builder

# ... 省略

# 运行阶段
FROM scratch

# 从编译阶段的中拷贝编译结果到当前镜像中
COPY --from=builder /build/server /
```

