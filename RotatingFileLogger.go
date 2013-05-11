package RotatingFileLogger

import (
	"fmt"
	"log"
	"os"
)

//--------------------
// LOG LEVEL
//--------------------

// Log levels to control the logging output.
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

type MyLogger struct {
	*log.Logger           // package log
	level        int      // log output level
	baseFileName string   // log file name
	maxBytes     int64    // log output file max size byte
	backupCount  int      // log output max number
	baseFilePrt  *os.File // log file point
}

// create logger
// fileLogger, _ := NewLog(0, "abc.log", 10*1024, 5)
// stderrLogger, _ := NewLog(0, "", 0, 0)
func NewLog(level int, baseFileName string, maxBytes int64,
	backupCount int) (*MyLogger, error) {
	logger := new(MyLogger)
	if baseFileName != "" {
		fprt, err := os.OpenFile(baseFileName,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			return nil, err
		}
		logger.Logger = log.New(fprt, "", log.Ldate|log.Ltime)
		logger.level = level
		logger.baseFileName = baseFileName
		logger.maxBytes = maxBytes
		logger.backupCount = backupCount
		logger.baseFilePrt = fprt
		return logger, nil
	}
	logger.Logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
	return logger, nil
}

// output logger msg
func (m *MyLogger) emit(record ...interface{}) error {
	ok, err := m.shouldRollover(record...)
	if err != nil {
		return err
	}
	if ok {
		m.doRollover()
	}
	m.Println(record...)
	return nil
}

// Do a rollover
func (m *MyLogger) doRollover() {
	if m.backupCount > 0 {
		for i := m.backupCount - 1; i >= 0; i-- {
			sfn := fmt.Sprintf("%s.%d", m.baseFileName, i)
			dfn := fmt.Sprintf("%s.%d", m.baseFileName, i+1)
			f, err := os.Open(dfn)
			if err == nil {
				f.Close()
				os.Remove(dfn)
			}
			os.Rename(sfn, dfn)
		}
		dfn := m.baseFileName + ".1"
		f, err := os.Open(dfn)
		if err == nil {
			f.Close()
			os.Remove(dfn)
		}
		m.baseFilePrt.Close()
		os.Rename(m.baseFileName, dfn)
		fprt, _ := os.OpenFile(m.baseFileName,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		m.Logger = log.New(fprt, "", log.Ldate|log.Ltime)
		m.baseFilePrt = fprt
	}
}

// Determine if rollover should occur.
// Basically, see if the supplied record would cause the file to exceed
// the size limit we have.
func (m *MyLogger) shouldRollover(record ...interface{}) (ok bool,
	err error) {
	if m.maxBytes > 0 {
		fileinfo, err := m.baseFilePrt.Stat()
		if err != nil {
			return false, err
		}
		var recordSize = len([]byte(fmt.Sprintf("", record...)))
		if fileinfo.Size()+int64(recordSize) >= m.maxBytes {
			return true, nil
		}
	}
	return false, nil
}

// Trace logs a message at trace level.
func (m *MyLogger) Trace(v ...interface{}) {
	if m.level <= LevelTrace {
		var msg = []interface{}{"[T]"}
		msg = append(msg, v...)
		m.emit(msg...)
	}
}

// Debug logs a message at debug level.
func (m *MyLogger) Debug(v ...interface{}) {
	if m.level <= LevelDebug {
		var msg = []interface{}{"[D]"}
		msg = append(msg, v...)
		m.emit(msg...)
	}
}

// Info logs a message at info level.
func (m *MyLogger) Info(v ...interface{}) {
	if m.level <= LevelInfo {
		var msg = []interface{}{"[I]"}
		msg = append(msg, v...)
		m.emit(msg...)
	}
}

// Warning logs a message at warning level.
func (m *MyLogger) Warn(v ...interface{}) {
	if m.level <= LevelWarning {
		var msg = []interface{}{"[W]"}
		msg = append(msg, v...)
		m.emit(msg...)
	}
}

// Error logs a message at error level.
func (m *MyLogger) Error(v ...interface{}) {
	if m.level <= LevelError {
		var msg = []interface{}{"[E]"}
		msg = append(msg, v...)
		m.emit(msg...)
	}
}

// Critical logs a message at critical level.
func (m *MyLogger) Critical(v ...interface{}) {
	if m.level <= LevelCritical {
		var msg = []interface{}{"[C]"}
		msg = append(msg, v...)
		m.emit(msg...)
	}
}
