package server

// peer提交给manager的
type workerReg struct {
	worker *worker
	reply  chan bool
}

// manager提交给peer的
type job struct {
	data      []byte
	replyChan chan []byte
}

