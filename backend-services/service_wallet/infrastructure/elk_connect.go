package infrastructure

import (
	"fmt"
	"net"
	"os"
	"service-wallet/utils"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewConnectionELK() (*zap.Logger, error) {
	serviceName := os.Getenv("SERVICE_NAME")
	serviceEnv := os.Getenv("SERVICE_ENV")
	urlLogstash := fmt.Sprintf("%s:%s", os.Getenv("LOGSTASH_HOST"), os.Getenv("LOGSTASH_PORT"))
	conn, err := net.Dial("tcp", urlLogstash)
	if err != nil {
		errMsg := utils.MessageError("infrastructure::NewConnectionELK", fmt.Errorf("failed dial logstash. err : %v", err))
		return nil, errMsg
	}

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// 🔥 ASYNC BUFFER
	bufferedWriter := &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(conn),
		Size:          256 * 1024,      // 256KB buffer
		FlushInterval: 5 * time.Second, // flush tiap 5 detik
	}

	core := zapcore.NewCore(
		encoder,
		bufferedWriter,
		zap.InfoLevel,
	)

	logger := zap.New(
		core, zap.AddCaller(), zap.AddCallerSkip(1),
	).With(
		zap.String("service", serviceName),
		zap.String("env", serviceEnv),
	)

	return logger, nil
}
