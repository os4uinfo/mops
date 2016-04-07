package server

import (
     _  "fmt"
)

const (
    chanBufSize = 1000
)

var (
    workerRegChan   = make(chan *workerReg, chanBufSize)
    workerUnRegChan = make(chan string, chanBufSize)
)

func workerManager(listenAddr string) {
    go workerListener(listenAddr)
    workerMap := make(map[string]*worker)
    for {
        select {
        case reg := <-workerRegChan:
            processWorkerReg(workerMap, reg)
        }
    }
}


func registerWorker(w *worker) bool {
    okChan := make(chan bool, 1)
    workerRegChan <- &workerReg{w, okChan}
    return <-okChan
}

func unregisterWorker(w *worker) {
    workerUnRegChan <- w.Id
}

func processWorkerReg(workerMap map[string]*worker, wr *workerReg) {
	worker := wr.worker
	if _, ok := workerMap[worker.Id]; ok {
		logError("worker已经存在，但又收到了注册申请!", worker)
		wr.reply <- false
		return
	}
	workerMap[worker.Id] = worker
	wr.reply <- true
	onWorkerConnected(worker.Id)
	logInfo("worker注册成功!", worker.Id)
}
