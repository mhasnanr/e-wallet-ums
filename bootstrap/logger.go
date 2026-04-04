package bootstrap

import "go.uber.org/zap"

var Log Logger

func SetupZapLogger() error {
	if Log != nil {
		return nil
	}

	opts := []zap.Option{}
	l, err := zap.NewProduction(opts...)
	if err != nil {
		return err
	}
	sl := l.Sugar()

	Log = &zapSugaredLogger{log: sl}
	return nil
}
