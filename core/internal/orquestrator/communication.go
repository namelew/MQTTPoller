package orquestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/data"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/data/models"
	"github.com/namelew/mqtt-bm-latency/packages/logs"
	"github.com/namelew/mqtt-bm-latency/packages/messages"
	tout "github.com/namelew/mqtt-bm-latency/packages/timeout"
	"github.com/namelew/mqtt-bm-latency/packages/waitgroup"
	"golang.org/x/exp/slices"
)

type queue struct {
	items []messages.ExperimentResult
	m     *sync.Mutex
}

type Orquestrator struct {
	log       *logs.Log
	client    mqtt.Client
	waitGroup *waitgroup.WaitGroup
	response  *queue
	tolerance int
}

func Build(c mqtt.Client, l *logs.Log, t int, hki int) *Orquestrator {
	return &Orquestrator{
		log:       l,
		waitGroup: waitgroup.New(),
		response: &queue{
			items: []messages.ExperimentResult{},
			m:     &sync.Mutex{},
		},
		client:    c,
		tolerance: t,
	}
}

func (o Orquestrator) ListWorkers() []models.Worker {
	return data.WorkersTable.List()
}

func (o Orquestrator) GetWorker(id string) *models.Worker {
	worker, _ := data.WorkersTable.Get(id)
	return &worker
}

func (o Orquestrator) timeout(t *tout.Timeout, timeHandler func(t context.Context, tk, tp string, tl int)) {
	timer, cancel := context.WithCancel(context.Background())

	o.client.Subscribe(t.OID+t.Topic, 1, func(c mqtt.Client, m mqtt.Message) {
		cancel()
		go o.timeout(t, timeHandler)
	})

	go timeHandler(timer, t.OID, t.Topic, t.Tolerance)
}

func (o Orquestrator) Init() error {
	o.log.Register("Starting database")

	data.Init(o.log)

	o.client.Subscribe("Orquestrator/Register", 1, func(c mqtt.Client, m mqtt.Message) {
		go func(messagePayload []byte) {
			var clientID string = uuid.New().String()[:]
			worker := string(messagePayload)

			data.WorkersTable.Add(clientID, models.Worker{ID: clientID, KeepAliveDeadline: 1, Online: false})

			o.setMessageHandler(&clientID)

			o.log.Register("worker " + worker + " registed as " + clientID)

			o.client.Publish("Orquestrator/Register/Log", 1, false, worker+" "+clientID)
		}(m.Payload())
	})

	o.client.Subscribe("Orquestrator/Login", 1, func(c mqtt.Client, m mqtt.Message) {
		go func(messagePayload []byte) {
			token := string(messagePayload)

			worker, err := data.WorkersTable.Get(token)

			if err != nil {
				o.log.Register("Login error: unable to find worker with id " + token)
				return
			}

			worker.Online = true

			data.WorkersTable.Update(worker.ID, worker)

			o.log.Register("worker " + token + " loged")

			o.setMessageHandler(&token)

			(o.client.Publish(token+"/Login/Log", 1, false, "true")).Wait()

			go o.timeout(&tout.Timeout{
				OID:       token,
				Topic:     "/KeepAlive",
				RecLimit:  -1,
				Tolerance: o.tolerance,
			}, func(t context.Context, tk, tp string, tl int) {
				select {
				case <-t.Done():
					return
				case <-time.After(time.Second * time.Duration(tl)):
					o.client.Unsubscribe(tk + tp)
					o.log.Register("lost connection with worker " + tk)

					worker, err := data.WorkersTable.Get(tk)

					if err != nil {
						o.log.Register("Timeout error: unable to find worker with id " + tk)
						return
					}

					worker.Online = false

					data.WorkersTable.Update(worker.ID, worker)
				}
			})
		}(m.Payload())
	})

	return nil
}

func (o Orquestrator) End() {
	o.log.Register("disconnect mqtt.client")

	o.client.Disconnect(0)

	o.log.Register("shutdown")
}

func (o *Orquestrator) setMessageHandler(t *string) {
	o.client.Subscribe(*t+"/Experiments/Results", 2, func(c mqtt.Client, m mqtt.Message) {
		go func(payload []byte) {
			var output messages.ExperimentResult

			err := json.Unmarshal(payload, &output)

			if err != nil {
				log.Println(err.Error())
				return
			}

			o.response.m.Lock()
			o.response.items = append(o.response.items, output)
			o.response.m.Unlock()
			o.waitGroup.Done()
		}(m.Payload())
	})

	o.client.Subscribe(*t+"/Experiments/Status", 2, func(c mqtt.Client, m mqtt.Message) {
		func(payload []byte) {
			var exps messages.Status
			json.Unmarshal(payload, &exps)

			tokens := strings.Split(exps.Type, " ")
			expid, _ := strconv.Atoi(tokens[2])

			exp, err := data.ExperimentTable.Get(uint64(expid))

			if err != nil {
				return
			}

			switch exps.Status {
			case "start":
				return
			case "finish":
				return
			default:
				if exp.ID != 0 {
					o.waitGroup.Done()
				}
			}
		}(m.Payload())
	})
}

func (o *Orquestrator) StartExperiment(arg messages.Start) models.Experiment {
	expid := time.Now().Unix()

	o.response.m.Lock()
	o.response.items = nil
	o.response.m.Unlock()

	var cmd messages.Command
	var description messages.CommandExperiment
	var nwkrs int = 0
	var timeoutIsValid = true
	var validMutex sync.Mutex = sync.Mutex{}

	cmd.Name = "start"
	cmd.CommandType = "experiment command"

	description.Expid = expid
	description.Declaration = arg.Description

	err := description.Attach(&cmd)

	if err != nil {
		return models.Experiment{
			ID:     uint64(expid),
			Finish: true,
			Error:  err.Error(),
		}
	}

	msg, err := json.Marshal(cmd)

	if err != nil {
		return models.Experiment{
			ID:     uint64(expid),
			Finish: true,
			Error:  err.Error(),
		}
	}

	experiment := models.Experiment{
		ID:                    uint64(expid),
		Finish:                false,
		Broker:                arg.Description.Broker,
		Port:                  arg.Description.Port,
		MqttVersion:           arg.Description.MqttVersion,
		NumPublishers:         arg.Description.NumPublishers,
		NumSubscriber:         arg.Description.NumSubscriber,
		QosPublisher:          arg.Description.QosPublisher,
		QosSubscriber:         arg.Description.QosSubscriber,
		SharedSubscrition:     arg.Description.SharedSubscrition,
		Retain:                arg.Description.Retain,
		Topic:                 arg.Description.Topic,
		Payload:               arg.Description.Payload,
		NumMessages:           arg.Description.NumMessages,
		RampUp:                arg.Description.RampUp,
		RampDown:              arg.Description.RampDown,
		Interval:              arg.Description.Interval,
		SubscriberTimeout:     arg.Description.SubscriberTimeout,
		ExecTime:              arg.Description.ExecTime,
		LogLevel:              arg.Description.LogLevel,
		Ntp:                   arg.Description.Ntp,
		Output:                arg.Description.Output,
		User:                  arg.Description.User,
		Password:              arg.Description.Password,
		TlsTrustsore:          arg.Description.TlsTrustsore,
		TlsTruststorePassword: arg.Description.TlsTruststorePassword,
		TlsKeystore:           arg.Description.TlsKeystore,
		TlsKeystorePassword:   arg.Description.TlsKeystorePassword,
		WorkerIDs:             arg.Id,
	}

	data.ExperimentTable.Add(uint64(expid), experiment)

	if arg.Id[0] == "" {
		workers := data.WorkersTable.List()
		nwkrs = len(workers)
		workersIDs := make([]string, 0, nwkrs)

		if nwkrs < 1 {
			const errMessage = "Fail in experiment request. No workers"
			o.log.Register(errMessage)
			experiment.Error = errMessage
			return experiment
		}

		o.waitGroup.Add(nwkrs)

		for i := range workers {
			if !workers[i].Online {
				o.log.Register("Worker " + workers[i].ID + " is off, skipping")
				continue
			}

			workersIDs = append(workersIDs, workers[i].ID)

			(o.client.Publish(workers[i].ID+"/Command", 2, false, msg)).Wait()

			go o.expTimeount(&timeoutIsValid, &validMutex, uint64(expid), arg.Description.ExecTime*5, arg.Attempts)

			o.log.Register("Requesting experiment in worker " + workers[i].ID)
		}

		experiment.WorkerIDs = workersIDs
	} else {
		workers := make([]*models.Worker, 10)

		for _, i := range arg.Id {
			worker, err := data.WorkersTable.Get(i)

			if err != nil {
				errorMessage := fmt.Sprintf("Experiment Error: Unable to find worker %s, skipping", i)
				o.log.Register(errorMessage)
				experiment.Error = errorMessage
				return experiment
			}

			workers = append(workers, &worker)
		}

		nwkrs = len(workers)

		o.waitGroup.Add(nwkrs)

		for i := range workers {
			if !workers[i].Online {
				o.log.Register("Worker " + workers[i].ID + " is off, skipping")
			}

			(o.client.Publish(workers[i].ID+"/Command", 1, false, msg)).Wait()

			go o.expTimeount(&timeoutIsValid, &validMutex, uint64(expid), arg.Description.ExecTime*5, arg.Attempts)

			o.log.Register("Requesting experiment in worker " + workers[i].ID)
		}
	}

	o.waitGroup.Wait()

	validMutex.Lock()
	timeoutIsValid = false
	validMutex.Unlock()

	o.response.m.Lock()

	if err != nil {
		experiment.Error = "failed to run experiment, don't find experiment id database"
		return experiment
	}

	experiment.Finish = true
	experiment.Results = slices.Clone[[]messages.ExperimentResult](o.response.items)

	if len(o.response.items) < nwkrs {
		o.response.m.Unlock()
		experiment.Error = fmt.Sprintf("%d workers have failed to run the experiment", nwkrs-len(o.response.items))
		data.ExperimentTable.Update(experiment.ID, experiment)
		return experiment
	}

	o.response.m.Unlock()

	data.ExperimentTable.Update(experiment.ID, experiment)

	return experiment
}

func (o Orquestrator) CancelExperiment(expid int64) error {
	cmd := messages.Command{Name: "cancel", CommandType: "moderation command", Arguments: make(map[string]interface{})}

	cmd.Arguments["id"] = expid

	msg, err := json.Marshal(cmd)

	if err != nil {
		o.log.Register("Unable to build cancel message")
		return err
	}

	experiment, err := data.ExperimentTable.Get(uint64(expid))

	if err != nil {
		o.log.Register("Unable to find experimento to cancel")
		return err
	}

	for _, workerid := range experiment.WorkerIDs {
		token := o.client.Publish(workerid+"/Command", 2, false, msg)

		go func(t mqtt.Token, wid string) {
			<-t.Done()

			if t.Error() != nil {
				o.log.Register("Unable to send experiment cancel message to worker " + wid + ". " + t.Error().Error())
			}
		}(token, workerid)

		o.waitGroup.Done()
	}

	return nil
}

func (o Orquestrator) expTimeount(valid *bool, vm *sync.Mutex, id uint64, tolerance int, attemps uint) {
	<-time.After(time.Second * time.Duration(tolerance))

	vm.Lock()
	if *valid {
		vm.Unlock()
		if attemps == 0 {
			o.waitGroup.Done()
		} else {
			o.expTimeount(valid, vm, id, tolerance, attemps-1)
		}
	}
}
