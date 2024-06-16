package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
    Logger *zap.Logger
)

func Setup(isSlave bool, masterAddr string) {
    encoder := zap.NewProductionEncoderConfig()
    encoder.EncodeTime = zapcore.RFC3339TimeEncoder
    encoder.EncodeLevel = zapcore.CapitalLevelEncoder

    config := zap.Config{
        Level: zap.NewAtomicLevelAt(zap.DebugLevel),
        Development: false,
        OutputPaths: []string{
            "stdout",
        },
        ErrorOutputPaths: []string{
            "stderr",
        },
        Encoding: "json",
        EncoderConfig: encoder,
        InitialFields: map[string]interface{}{
            "cluster_info": map[string]string{
                "replica_status": "master",
                "master_address": masterAddr,
            },
        },
    }

    if isSlave {
        config.InitialFields["cluster_info"].
            (map[string]string)["replica_status"] = "slave"
    }

    Logger = zap.Must(config.Build())
}

