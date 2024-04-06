package logger

import (
	"github.com/jagerente/discord-lootbox-opener/internal/gui"
	"sync"
	"time"
)

type Logger struct {
	*sync.RWMutex
	cap  int
	logs []gui.Log
}

func New(cap int) *Logger {
	return &Logger{
		cap:     cap,
		logs:    make([]gui.Log, 0, cap),
		RWMutex: &sync.RWMutex{},
	}
}

func (l *Logger) Log(content string) {
	l.Lock()
	defer l.Unlock()

	if len(l.logs) >= l.cap {
		l.logs = l.logs[1:]
	}

	l.logs = append(l.logs, gui.Log{
		CreatedAt: time.Now(),
		Content:   content,
	})
}

func (l *Logger) GetLogs() []gui.Log {
	l.RLock()
	defer l.RUnlock()

	return l.logs
}
