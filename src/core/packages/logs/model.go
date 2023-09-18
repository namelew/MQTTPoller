package logs

import (
	"log"
	"os"
	"sync"

	"github.com/namelew/mqtt-poller/src/core/packages/utils"
)

type Log struct {
	file  string
	mutex *sync.Mutex
}

func Build(f string) *Log {
	return &Log{
		file:  f,
		mutex: &sync.Mutex{},
	}
}

func (l *Log) Create() {
	l.mutex.Lock()
	if !utils.FileExists(l.file) {
		f, err := os.Create(l.file)
		if err != nil {
			log.Panic(err.Error())
		}
		f.Close()
	} else {
		os.Truncate(l.file, 0)
	}
	l.mutex.Unlock()
}

func (l *Log) Register(s string) {
	l.mutex.Lock()
	f, err := os.OpenFile(l.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	log.Println(s)

	if err != nil {
		log.Panic(err.Error())
	}
	defer f.Close()

	f.WriteString(s + "\n")
	l.mutex.Unlock()
}

func (l *Log) Fatal(s string) {
	l.Register(s)
	os.Exit(1)
}
