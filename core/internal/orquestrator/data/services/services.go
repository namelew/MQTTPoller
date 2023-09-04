package services

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/data/models"
	"github.com/namelew/mqtt-bm-latency/packages/logs"
	"github.com/tidwall/btree"
	"golang.org/x/exp/constraints"
)

type TableData interface {
	models.Worker | models.Experiment
}

type Table[I constraints.Ordered, K TableData] struct {
	data     btree.Map[I, K]
	filename string
	mutex    sync.Mutex
	log      *logs.Log
	FreeSlot uint64
}

const syncer_timeout time.Duration = time.Second * 3

func Build[I constraints.Ordered, K TableData](filename string, l *logs.Log) *Table[I, K] {
	new := Table[I, K]{
		log:      l,
		filename: filename,
	}

	go new.syncer()

	return &new
}

func (h *Table[I, K]) Add(id I, data K) {
	h.mutex.Lock()

	h.data.Set(id, data)
	h.FreeSlot++

	h.mutex.Unlock()
}

func (h *Table[I, K]) Remove(id I) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, found := h.data.Delete(id)

	if !found {
		return errors.New("unable to find data in the tree")
	}

	return nil
}

func (h *Table[I, K]) Update(id I, new K) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, exists := h.data.Get(id)

	if !exists {
		return errors.New("unable to find data in the tree")
	}

	h.data.Set(id, new)

	return nil
}

func (h *Table[I, K]) List() []K {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.data.Values()
}

func (h *Table[I, K]) Get(id I) (K, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	item, found := h.data.Get(id)

	if !found {
		return item, errors.New("unable to find element in tree")
	}

	return item, nil
}

func (h *Table[I, K]) syncer() {
	for {
		<-time.After(syncer_timeout)

		h.mutex.Lock()

		f, err := os.OpenFile(h.filename, os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			h.log.Register("Unable to open data file: " + h.filename + ". " + err.Error())
			h.mutex.Unlock()
			continue
		}

		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)

		if err := encoder.Encode(h.data.Values()); err != nil {
			h.log.Register("unable to enconde data to write in " + h.filename + ". " + err.Error())
			h.mutex.Unlock()
			continue
		}

		_, err = f.Write(buffer.Bytes())

		if err != nil {
			h.log.Register("unable to write in the datafile " + h.filename + ". " + err.Error())
			h.mutex.Unlock()
			continue
		}

		h.mutex.Unlock()
	}
}
