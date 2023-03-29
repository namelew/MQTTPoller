package logs

import (
	"log"
	"os"
	"sync"

	"github.com/namelew/mqtt-bm-latency/utils"
)

type Log struct {
	filename string
	m sync.Mutex
}

func Build(filename string) *Log{
	return &Log{
		filename: filename,
		m: sync.Mutex{},
	}
}

func (l *Log) Create() {
	if !utils.FileExists("worker.log"){
		f,_ := os.Create("worker.log")
		f.Close()
	} else{
		os.Truncate("worker.log", 0)
	}
}

func (l *Log) Register(msg string) {
	l.m.Lock()
	log.Println(msg)
	f,_ := os.OpenFile(l.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(msg + "\n")
	f.Close()
	l.m.Unlock()
}