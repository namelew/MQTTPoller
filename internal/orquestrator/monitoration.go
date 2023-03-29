package orquestrator

// func receiveControl(id int, timeout int) {
// 	var start int64 = time.Now().UnixMilli()
// 	exlog := workers[id].Historic.FindLarger()
// 	for utils.AbsInt(int(time.Now().UnixMilli()-start)) < (timeout * 1000) {
// 		if workers[id].ReceiveConfirmation || !workers[id].Status || exlog.Err {
// 			break
// 		}
// 	}

// 	workers[id].ReceiveConfirmation = false
// 	expWG.Done()

// 	if (timeout*1000) <= utils.AbsInt(int(time.Now().UnixMilli()-start)) || !workers[id].Status || exlog.Err {
// 		log.Printf("Error in worker %d: experiment don't return\n", id)
// 		exlog.Finished = true
// 		redoExperiment(id, exlog)
// 	}
// }

// func watcher(id int, tl int) {
// 	var start int64 = time.Now().UnixMilli()

// 	for utils.AbsInt(int(time.Now().UnixMilli()-start)) < (tl * 1000) {
// 		if workers[id].TestPing {
// 			return
// 		}
// 	}
// 	workers[id].Status = false
// 	workers[id].TestPing = true

// 	token := client.Unsubscribe(workers[id].Id + "/Experiments/Results")
// 	token.Wait()

// 	log.Printf("Worker %d is off\n", id)
// }

// func redoExperiment(worker int, experiment *messages.ExperimentLog) {
// 	exp := *experiment
// 	workers[worker].Historic.Remove(experiment.Id)

// 	if len(workers) <= 1 {
// 		return
// 	}

// 	if exp.Attempts > 0 {
// 		exp.Attempts--
// 		size := len(workers)
// 		var sample = make([]int, 0, size)
// 		var timeout int

// 		for i := 0; i < size; i++ {
// 			if i != worker && workers[i].Status {
// 				sample = append(sample, i)
// 			}
// 		}

// 		cmdExp := exp.Cmd.ToCommandExperiment()
// 		exp.Id = time.Now().Unix()
// 		cmdExp.Expid = exp.Id
// 		cmdExp.Attempts = exp.Attempts
// 		timeout = cmdExp.ExecTime * 5 * 2
// 		cmdExp.Attach(&exp.Cmd)

// 		nw := sample[rand.Intn(len(sample))]

// 		msg, err := json.Marshal(exp.Cmd)

// 		if err != nil {
// 			log.Fatal(err.Error())
// 		}

// 		workers[nw].Historic.Add(exp.Id, exp.Cmd, exp.Attempts)

// 		token := client.Publish(workers[nw].Id+"/Command", byte(1), false, msg)
// 		token.Wait()

// 		go receiveControl(nw, timeout)
// 	}
// }
