package main

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
    "math/rand"
    "time"
)

// 创建 zap logger，配置文件切割（lumberjack）
func NewLogger() *zap.Logger {
    // 配置 lumberjack 来进行日志轮转
    lumberJackLogger := &lumberjack.Logger{
        Filename:   "./logfile.log",  // 日志文件路径
        MaxSize:    10,               // 每个日志文件的大小（MB）
        MaxBackups: 3,                // 保留的旧日志文件数
        MaxAge:     28,               // 保留的天数
        Compress:   true,             // 是否压缩旧日志文件
    }

    // zap 的核心配置
    writeSyncer := zapcore.AddSync(lumberJackLogger)
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(encoderConfig), // JSON 编码器
        writeSyncer,                           // 日志输出目的地
        zap.InfoLevel,                         // 日志级别，日志级别设置为 zap.InfoLevel，这意味着所有 等于或高于 info 级别 的日志（如 info、warn、error、fatal 等）都会写入日志文件中。因此，除了 info 级别之外，其他更高优先级的日志（warn、error、fatal 等） 也会写入日志文件。
    )

    logger := zap.New(core, zap.AddCaller()) // 加入调用者信息

    return logger
}

func main() {
    // 初始化随机种子
    rand.Seed(time.Now().UnixNano())

    // 创建 logger
    logger := NewLogger()

    // 输出不同级别的日志
    logger.Debug("This is a debug message")  // 不会写入文件
    logger.Info("This is an info message")    // 会写入文件
    logger.Warn("This is a warning message")  // 会写入文件
    logger.Error("This is an error message")  // 会写入文件

    // 记得同步日志缓冲区
    defer logger.Sync()
}
