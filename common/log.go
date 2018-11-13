package common

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

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
		logger.trace = log.New(w, "[T] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		logger.info = log.New(w, "[I] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		logger.warn = log.New(w, "[W] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		logger.error = log.New(w, "[E] ", log.Ldate|log.Ltime|log.Lmicroseconds)
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
	for i := 0; i < level; i++ {
		_, file, line, ok := runtime.Caller(3 + i)
		if !ok {
			break
		}
		t += fmt.Sprintln("in", file, line)
	}
	return t
}
