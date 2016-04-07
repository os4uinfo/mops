package server

import (
    "log"
    "fmt"
    "net"
    "time"
    "mops/msg"
    "mops/msg/pb"
)

type worker struct {
	Id                  string
	Path                string
	MachineStatProto    *pb.MachineStatProto
	ServerStatProtoMap  map[string]*pb.ServerStatProto
	waitingSendJobsChan chan *job
	conn                net.Conn
}

func workerListener(addr string) {
    logInfo("Serer listener start!", addr)
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatal(err)
    }
    for {
        conn, err := ln.Accept()
        if err != nil {
            logError("worker listener accept error!", err)
            continue
        }
        fmt.Println(conn)
        go workerHandler(conn)
    }
}

func workerHandler(c net.Conn) {
    workerId, err := msg.ReadString(c)
    if err != nil {
        logError("workerHandler读workerId出错！", err)
        c.Close()
        return
    }
    workerPath, err := msg.ReadString(c)
    if err != nil {
        logError("workerHandler读workerPath出错！", err)
        c.Close()
        return
    }

    waitingSendJobsChan := make(chan *job, 100)
    worker := &worker{
        Id:                  workerId,
        Path:                workerPath,
        MachineStatProto:    &pb.MachineStatProto{},
        ServerStatProtoMap:  make(map[string]*pb.ServerStatProto),
        waitingSendJobsChan: waitingSendJobsChan,
        conn:                c,
    }

    if !registerWorker(worker) {
        msg.WriteBool(c, false)
        c.Close()
        return
    }

    err = msg.WriteBool(c, true)
    if err != nil {
        logError("workerHandler写worker注册成功出错！", err)
        c.Close()
        return
    }

    defer func() {
        if err := recover(); err != nil {
            logError("workerHandler内部错误！", err)
            unregisterWorker(worker)
            c.Close()
        }
    }()
//    jobIdCounter := 0
//    waitingReplyJobsMap := make(map[int]chan []byte)

    jobsReplyChan := make(chan []byte, 10)
    workerCloseChan := make(chan int)

    heartBeatTimeOut := false


    go func() {
        defer func() {
            if err := recover(); err != nil {
                logError("workerHandler读数据内部错误！", err)
                workerCloseChan <- 1
            }
        }()

        for {
            data, err := msg.ReadBytes(c)
            if err != nil {
                workerCloseChan <- 1
                return
            }
            jobsReplyChan <- data

            heartBeatTimeOut = false
        }
    }()

    go func() {
        tickChan := time.Tick(time.Second * 30)
        for _ = range tickChan {
            if heartBeatTimeOut {
                logError("workerHandler heartBeatTimeOut:", workerId)
                workerCloseChan <- 1
            } else {
                heartBeatTimeOut = true
            }
        }
    }()
}
