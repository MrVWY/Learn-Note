package main

import (
    "os"
    "math/rand"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

// 创建 zap logger，配置文件切割（lumberjack）和终端输出
func NewLogger() *zap.Logger {
    // 配置 lumberjack 来进行日志轮转
    lumberJackLogger := &lumberjack.Logger{
        Filename:   "./logfile.log",  // 日志文件路径
        MaxSize:    10,               // 每个日志文件的大小（MB）
        MaxBackups: 3,                // 保留的旧日志文件数
        MaxAge:     28,               // 保留的天数
        Compress:   true,             // 是否压缩旧日志文件
    }

    // 创建文件写入器（lumberjack）
    fileWriteSyncer := zapcore.AddSync(lumberJackLogger)

    // 创建终端写入器（标准输出）
    consoleWriteSyncer := zapcore.AddSync(os.Stdout)

    // zap 的编码器配置
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    // 创建 JSON 编码器和 Console 编码器
    fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
    consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

    // 创建两个 Core，一个用于文件，一个用于终端
    core := zapcore.NewTee(
        zapcore.NewCore(fileEncoder, fileWriteSyncer, zap.InfoLevel),     // 输出到文件
        zapcore.NewCore(consoleEncoder, consoleWriteSyncer, zap.InfoLevel), // 输出到终端
    )

    // 创建 Logger 并返回
    return zap.New(core, zap.AddCaller())
}

// 按比例输出 info 日志
func LogInfoWithProbability(logger *zap.Logger, message string, probability float64) {
    if rand.Float64() < probability {
        logger.Info(message)
    }
}

func main() {
    // 创建 logger
    logger := NewLogger()

    // 示例：以 30% 的概率输出 info 级别的日志
    for i := 0; i < 10; i++ {
        LogInfoWithProbability(logger, "This is an info message", 0.3)
    }
    
    // 输出不同级别的日志
    logger.Debug("This is a debug message")  // 不会输出
    logger.Info("This is an info message")    // 会输出到文件和终端
    logger.Warn("This is a warning message")  // 会输出到文件和终端
    logger.Error("This is an error message")  // 会输出到文件和终端

    // 记得同步日志缓冲区
    defer logger.Sync()
}
