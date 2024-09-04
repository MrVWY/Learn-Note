package main

import (
    "math/rand"
    "time"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
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
        zap.InfoLevel,                         // 日志级别
    )

    logger := zap.New(core, zap.AddCaller()) // 加入调用者信息

    return logger
}

// 按比例输出 info 日志
func LogInfoWithProbability(logger *zap.Logger, message string, probability float64) {
    if rand.Float64() < probability {
        logger.Info(message)
    }
}

func main() {
    // 初始化随机种子
    rand.Seed(time.Now().UnixNano())

    // 创建 logger
    logger := NewLogger()

    // 示例：以 30% 的概率输出 info 级别的日志
    for i := 0; i < 10; i++ {
        LogInfoWithProbability(logger, "This is an info message", 0.3)
    }

    // 记得同步日志缓冲区
    defer logger.Sync()
}

