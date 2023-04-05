package orquestrator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/databases"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/databases/filters"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/databases/models"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/databases/services/experiments"
	seworkers "github.com/namelew/mqtt-bm-latency/internal/orquestrator/databases/services/workers"
	"github.com/namelew/mqtt-bm-latency/packages/housekeeper"
	"github.com/namelew/mqtt-bm-latency/packages/logs"
	"github.com/namelew/mqtt-bm-latency/packages/messages"
	local "github.com/namelew/mqtt-bm-latency/packages/network"
	tout "github.com/namelew/mqtt-bm-latency/packages/timeout"
	"github.com/namelew/mqtt-bm-latency/packages/waitgroup"
)

type queue struct {
	items []messages.ExperimentResult
	m     *sync.Mutex
}

type Orquestrator struct {
	log         *logs.Log
	workers     *seworkers.Workers
	experiments *experiments.Experiments
	client      *local.Client
	waitGroup   *waitgroup.WaitGroup
	hk          *housekeeper.Housekeeper
	response    *queue
	repress     *queue
	tolerance   int
}

func Build(c *local.Client, t int, hki int) *Orquestrator {
	return &Orquestrator{
		log:         c.Log,
		workers:     seworkers.Build(c.Log),
		experiments: experiments.Build(c.Log),
		waitGroup:   waitgroup.New(),
		response: &queue{
			items: []messages.ExperimentResult{},
			m:     &sync.Mutex{},
		},
		repress: &queue{
			items: []messages.ExperimentResult{},
			m:     &sync.Mutex{},
		},
		hk:        housekeeper.New(time.Hour*time.Duration(hki), c.Log),
		client:    c,
		tolerance: t,
	}
}

func (o Orquestrator) ListWorkers(filter *filters.Worker) []models.Worker {
	return o.workers.List(filter)
}

func (o Orquestrator) GetWorker(id int) *models.Worker {
	return o.workers.Get(id)
}

func (o Orquestrator) timeout(t *tout.Timeout, timeHandler func(t context.Context, tk, tp string, tl int)) {
	timer, cancel := context.WithCancel(context.Background())

	o.client.Register(t.OID+t.Topic, 1, true, func(c paho.Client, m paho.Message) {
		cancel()
		go o.timeout(t, timeHandler)
	})

	go timeHandler(timer, t.OID, t.Topic, t.Tolerance)
}

func (o Orquestrator) Init() error {
	o.client.Create()

	o.log.Register("Starting database")

	databases.Connect(o.log)

	o.client.Register("Orquestrator/Register", 1, false, func(c paho.Client, m paho.Message) {
		go func(messagePayload []byte) {
			var clientID string = ""
			worker := string(messagePayload)

			for i := 0; i < 10; i++ {
				seed := rand.NewSource(time.Now().UnixNano())
				random := rand.New(seed)
				clientID += fmt.Sprintf("%d", random.Int()%10)
			}

			o.workers.Add(models.Worker{Token: clientID, KeepAliveDeadline: 1, Online: false, Experiments: nil})

			o.setMessageHandler(&clientID)

			o.log.Register("worker " + worker + " registed as " + clientID)

			o.client.Send("Orquestrator/Register/Log", worker+"-"+clientID)
		}(m.Payload())
	})

	o.client.Register("Orquestrator/Login", 1, false, func(c paho.Client, m paho.Message) {
		go func(messagePayload []byte) {
			token := string(messagePayload)

			o.workers.ChangeStatus(&filters.Worker{Token: token, Online: true})

			o.log.Register("worker " + token + " loged")

			o.setMessageHandler(&token)

			o.client.Send(token+"/Login/Log", "true")

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
					o.client.Unregister(tk + tp)
					o.log.Register("lost connection with worker " + tk)
					o.workers.ChangeStatus(&filters.Worker{Token: tk, Online: false})
				}
			})
		}(m.Payload())
	})

	o.hk.Place(o.experiments)
	o.hk.Place(o.workers)
	go o.hk.Start()

	return nil
}

func (o Orquestrator) End() {
	o.log.Register("disconnect mqtt.client")

	o.client.Disconnect(0)

	o.log.Register("shutdown")
}

func (o *Orquestrator) setMessageHandler(t *string) {
	o.client.Register(*t+"/Experiments/Results", 1, false, func(c paho.Client, m paho.Message) {
		var output messages.ExperimentResult

		err := json.Unmarshal(m.Payload(), &output)

		if err != nil {
			log.Fatal(err.Error())
		}

		exp := o.experiments.Get(output.Meta.ID)

		if exp.Finish {
			o.repress.m.Lock()
			o.repress.items = append(o.repress.items, output)
			o.repress.m.Unlock()
		} else {
			o.response.m.Lock()
			o.response.items = append(o.response.items, output)
			o.response.m.Unlock()
			o.waitGroup.Done()
		}
	})

	o.client.Register(*t+"/Experiments/Status", 1, false, func(c paho.Client, m paho.Message) {
		var exps messages.Status
		json.Unmarshal(m.Payload(), &exps)

		tokens := strings.Split(exps.Type, " ")
		expid, _ := strconv.Atoi(tokens[2])

		exp := o.experiments.Get(uint64(expid))

		switch exps.Status {
		case "start":
			return
		case "finish":
			return
		default:
			o.waitGroup.Done()
			if exp.ID != 0 {
				go o.experiments.Update(uint64(expid), models.Experiment{Finish: true, Error: exps.Status})
				//redoExperiment(id, exp)
			}
		}
	})
}

func (o *Orquestrator) StartExperiment(arg messages.Start) ([]messages.ExperimentResult, error) {
	expid := time.Now().Unix()

	o.response.m.Lock()
	o.response.items = nil
	o.response.m.Unlock()

	var cmd messages.Command
	var experiment messages.CommandExperiment
	var nwkrs int = 0

	cmd.Name = "start"
	cmd.CommandType = "experiment command"

	experiment.Expid = expid
	experiment.Declaration = arg.Description

	err := experiment.Attach(&cmd)

	if err != nil {
		return o.response.items, err
	}

	msg, err := json.Marshal(cmd)

	if err != nil {
		return o.response.items, err
	}

	o.experiments.Add(
		models.Experiment{
			ID:     uint64(expid),
			Finish: false,
		},
		models.ExperimentDeclaration{
			Tool:                  arg.Description.Tool,
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
		},
		arg.Id...,
	)

	if arg.Id[0] == -1 {
		workers := o.workers.List(nil)

		nwkrs = len(workers)

		o.waitGroup.Add(nwkrs)

		for i := range workers {
			if !workers[i].Online {
				o.log.Register("Worker " + workers[i].Token + " is off, skipping")
				continue
			}

			o.client.Send(workers[i].Token+"/Command", msg)

			go o.expTimeount(uint64(expid), arg.Description.ExecTime*5, arg.Attempts)

			o.log.Register("Requesting experiment in worker " + workers[i].Token)
		}
	} else {
		workers := make([]*models.Worker, 10)

		for _, i := range arg.Id {
			workers = append(workers, o.workers.Get(i))
		}

		nwkrs = len(workers)

		o.waitGroup.Add(nwkrs)

		for i := range workers {
			if !workers[i].Online {
				o.log.Register("Worker " + workers[i].Token + " is off, skipping")
			}

			o.client.Send(workers[i].Token+"/Command", msg)

			go func(id uint64, tolerance int) {
				<-time.After(time.Second * time.Duration(tolerance))
				exp := o.experiments.Get(id)

				if !exp.Finish {
					o.waitGroup.Done()
				}
			}(uint64(expid), arg.Description.ExecTime*5)

			o.log.Register("Requesting experiment in worker " + workers[i].Token)
		}
	}

	o.waitGroup.Wait()

	o.response.m.Lock()

	if len(o.response.items) < nwkrs {
		o.response.m.Unlock()
		go o.experiments.Update(uint64(expid), models.Experiment{Finish: true, Error: fmt.Sprintf("%d workers have failed to run the experiment", nwkrs-len(o.response.items))})
		return o.response.items, errors.New("failed to run experiment")
	}

	o.repress.m.Lock()
	o.response.items = append(o.response.items, o.repress.items...)
	o.response.m.Unlock()

	o.repress.items = nil
	o.repress.m.Unlock()

	go o.experiments.Update(uint64(expid), models.Experiment{Finish: true})

	return o.response.items, nil
}

func (o Orquestrator) CancelExperiment(id int, expid int64) error {
	cmd := messages.Command{Name: "cancel", CommandType: "moderation command", Arguments: make(map[string]interface{})}

	cmd.Arguments["id"] = expid

	msg, err := json.Marshal(cmd)

	if err != nil {
		return err
	}

	worker := o.workers.Get(id)

	o.client.Send(worker.Token+"/Command", msg)

	o.waitGroup.Done()

	return nil
}

func (o Orquestrator) expTimeount(id uint64, tolerance int, attemps uint) {
	<-time.After(time.Second * time.Duration(tolerance))
	exp := o.experiments.Get(id)

	if !exp.Finish {
		if attemps == 0 {
			o.waitGroup.Done()
		} else {
			o.expTimeount(id, tolerance, attemps-1)
		}
	}
}
