package logger

import (
	"github.com/potapovsv/bismuthCube/config"
	"log"
	"os"
	"sync"
)

type Logger struct {
	*log.Logger
	file *os.File
}

var (
	instance *Logger
	once     sync.Once
)

// Init инициализирует логгер
func Init(logToFile bool) *Logger {
	once.Do(func() {
		instance = createLogger(logToFile)
	})
	return instance
}

func createLogger(logToFile bool) *Logger {
	logger := &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lmsgprefix),
	}

	if logToFile {
		file, err := os.OpenFile("bismuth.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Printf("⚠️ Failed to open log file: %v", err)
		} else {
			logger.file = file
			logger.SetOutput(file)
		}
	}

	return logger
}

// Close освобождает ресурсы
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

// Get возвращает экземпляр логгера
func Get() *Logger {
	if instance == nil {
		return Init(false) // По умолчанию только в консоль
	}
	return instance
}

func (l *Logger) Debug(v ...interface{}) {
	if config.GetConfig().Logging.Level == "debug" {
		l.Printf("[DEBUG] %v", v)
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.Printf("[INFO] %v", v)
}

func (l *Logger) Warn(v ...interface{}) {
	l.Printf("[WARN] %v", v)
}

func (l *Logger) Error(v ...interface{}) {
	l.Printf("[ERROR] %v", v)
}
func (l *Logger) Rotate() error {
	if l.file == nil {
		return nil
	}

	newFile, err := os.OpenFile(l.file.Name(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	oldFile := l.file
	l.SetOutput(newFile)
	l.file = newFile
	oldFile.Close()

	return nil
}
func colorize(level string, msg string) string {
	colors := map[string]string{
		"INFO":  "\033[36m", // Cyan
		"WARN":  "\033[33m", // Yellow
		"ERROR": "\033[31m", // Red
	}
	return colors[level] + msg + "\033[0m"
}
