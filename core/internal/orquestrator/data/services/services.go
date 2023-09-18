package services

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/namelew/mqtt-poller/core/internal/orquestrator/data/models"
	"github.com/namelew/mqtt-poller/core/packages/logs"
	"github.com/namelew/mqtt-poller/core/packages/utils"
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

	if utils.FileExists(filename) && utils.FileExists(filename+".index") {
		new.loadFiles()
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

		datafile, err := os.OpenFile(h.filename, os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			h.log.Register("Unable to open data file: " + h.filename + ". " + err.Error())
			h.mutex.Unlock()
			continue
		}

		indexfile, err := os.OpenFile(h.filename+".index", os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			h.log.Register("Unable to open index file: " + h.filename + ".index. " + err.Error())
			h.mutex.Unlock()
			continue
		}

		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)

		if err := encoder.Encode(h.data.Values()); err != nil {
			datafile.Close()
			h.log.Register("unable to enconde data to write in " + h.filename + ". " + err.Error())
			h.mutex.Unlock()
			continue
		}

		_, err = datafile.Write(buffer.Bytes())

		if err != nil {
			datafile.Close()
			h.log.Register("unable to write in the datafile " + h.filename + ". " + err.Error())
			h.mutex.Unlock()
			continue
		}

		buffer.Reset()

		if err := encoder.Encode(h.data.Keys()); err != nil {
			indexfile.Close()
			h.log.Register("unable to enconde data to write in " + h.filename + ".index. " + err.Error())
			h.mutex.Unlock()
			continue
		}

		_, err = indexfile.Write(buffer.Bytes())

		if err != nil {
			h.log.Register("unable to write in the indexfile " + h.filename + ".index. " + err.Error())
			indexfile.Close()
			h.mutex.Unlock()
			continue
		}

		h.mutex.Unlock()
		datafile.Close()
		indexfile.Close()
	}
}

func (h *Table[I, K]) loadFiles() {
	var buffer bytes.Buffer
	decoder := gob.NewDecoder(&buffer)
	keys := make([]I, 0)
	values := make([]K, 0)

	h.mutex.Lock()
	defer h.mutex.Unlock()

	datafileBytes, err := os.ReadFile(h.filename)

	if err != nil {
		h.log.Fatal("Unable to open data file: " + h.filename + ". " + err.Error())
	}

	indexfileBytes, err := os.ReadFile(h.filename + ".index")

	if err != nil {
		h.log.Fatal("Unable to open index file: " + h.filename + "index. " + err.Error())
	}

	_, err = buffer.Write(datafileBytes)

	if err != nil {
		h.log.Fatal("Unable to store the data of the datafiel in a buffer. " + err.Error())
	}

	err = decoder.Decode(&values)

	if err != nil {
		h.log.Fatal("Unable to decode datafile data. " + err.Error())
	}

	buffer.Reset()

	_, err = buffer.Write(indexfileBytes)

	if err != nil {
		h.log.Fatal("Unable to store the data of the indexfile in a buffer. " + err.Error())
	}

	err = decoder.Decode(&keys)

	if err != nil {
		h.log.Fatal("Unable to decode indexfile data. " + err.Error())
	}

	if len(keys) != len(values) {
		h.log.Register("Inconsistent data. " + err.Error() + "\nDiscarting pre-session data")
		return
	}

	for i := range keys {
		h.data.Set(keys[i], values[i])
	}
}
