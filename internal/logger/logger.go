package logger

import (
	"log"
	"os"
	"sync"
)

var Log *Logger
var once sync.Once

// Logger 封装了标准日志记录器
type Logger struct {
	*log.Logger
}

func init() {
	once.Do(func() {
		Log = &Logger{
			Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		}
	})
}

// Info 记录信息消息
func (l *Logger) Info(v ...interface{}) {
	l.Print(append([]interface{}{"[INFO] "}, v...)...)
}

// Infof 记录带格式的信息消息
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Printf("[INFO] "+format, v...)
}

// Error 记录错误消息
func (l *Logger) Error(v ...interface{}) {
	l.Print(append([]interface{}{"[ERROR] "}, v...)...)
}

// Errorf 记录带格式的错误消息
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Printf("[ERROR] "+format, v...)
}

// Warn 记录警告消息
func (l *Logger) Warn(v ...interface{}) {
	l.Print(append([]interface{}{"[WARN] "}, v...)...)
}

// Warnf 记录带格式的警告消息
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Printf("[WARN] "+format, v...)
}

// Debug 记录调试消息
func (l *Logger) Debug(v ...interface{}) {
	l.Print(append([]interface{}{"[DEBUG] "}, v...)...)
}

// Debugf 记录带格式的调试消息
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Printf("[DEBUG] "+format, v...)
}
