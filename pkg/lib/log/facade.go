package log

func Debugln(args ...interface{}) {
	logger.Debugln(args...)
}

func Infoln(args ...interface{}) {
	logger.Infoln(args...)
}

func Warnln(args ...interface{}) {
	logger.Warnln(args...)
}

func Errorln(args ...interface{}) {
	logger.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	logger.Fatalln(args...)
}
