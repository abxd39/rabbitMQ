package common

import (
	"bufio"
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"sctek.com/typhoon/th-platform-gateway/hook"
	"time"
)

func InitLogger(dir, name, level string)error {
	if len(dir)==0||len(name)==0{
		return fmt.Errorf("日志文件目录或文件名的配置文件为空！！")
	}
	baseLogPath := fmt.Sprintf("%s%s", dir, name)
	writer, err := rotatelogs.New(
		baseLogPath+"_%Y%m%d.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)

	if err != nil {
		panic("config local file system logger error. ")
	}
	switch level {

	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stderr)
	case "info":
		setNull()
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		setNull()
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		setNull()
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		setNull()
		logrus.SetLevel(logrus.InfoLevel)
	}

	g := &logrus.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, g)
	logrus.AddHook(lfHook)
	return nil
}

func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	logrus.SetOutput(writer)
}

/*
type Logger struct {
	LogFile    string
	TraceLevel int
	trace      *log.Logger
	info       *log.Logger
	warn       *log.Logger
	error      *log.Logger
}

func NewLogger(logfile string, tracelevel int) (*Logger, error) {
	logger := new(Logger)
	logger.LogFile = logfile
	logger.TraceLevel = tracelevel
	if w, err := logger.getWriter(); err != nil {
		return logger, err
	} else {
		logger.trace = log.New(w, "[T] ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.info = log.New(w, "[I] ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.warn = log.New(w, "[W] ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.error = log.New(w, "[E] ", log.Ldate|log.Ltime|log.Lshortfile)
		return logger, err
	}
}

func (l *Logger) Traceln(v ...interface{}) {
	l.outputln(l.trace, l.TraceLevel, v...)
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.outputf(l.trace, l.TraceLevel, format, v...)
}

func (l *Logger) Infoln(v ...interface{}) {
	l.outputln(l.info, l.TraceLevel, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.outputf(l.info, l.TraceLevel, format, v...)
}

func (l *Logger) Warnln(v ...interface{}) {
	l.outputln(l.warn, l.TraceLevel, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.outputf(l.warn, l.TraceLevel, format, v...)
}

func (l *Logger) Errorln(v ...interface{}) {
	l.outputln(l.error, l.TraceLevel, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.outputf(l.error, l.TraceLevel, format, v...)
}

func (l *Logger) outputln(logger *log.Logger, tracelevel int, v ...interface{}) {
	s := fmt.Sprintln(v...) + l.getTraceInfo(tracelevel)
	logger.Output(3, s)
}

func (l *Logger) outputf(logger *log.Logger, tracelevel int, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...) + l.getTraceInfo(tracelevel)
	logger.Output(3, s)
}

func (l *Logger) getWriter() (io.Writer, error) {
	lf := l.LogFile
	if lf == "" {
		return os.Stdout, nil
	}
	return os.OpenFile(l.LogFile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

func (l *Logger) getTraceInfo(level int) string {
	t := ""
	//for i := 0; i < level; i++ {
	//	_, file, line, ok := runtime.Caller(3 + i)
	//	if !ok {
	//		break
	//	}
	//	//t += fmt.Sprintln("in", file, line)
	//}
	return t
}
*/
