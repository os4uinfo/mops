package client

import (
    "flag"
    "log"
    "sync/atomic"
    "mops/msg"
    "net"
    "fmt"
    "time"
    "strconv"
)

var (
    flags       = flag.NewFlagSet("worker", flag.ExitOnError)
    baseDir     string
    workerId    string
    masterAddr  string
    httpPort    string
)

var (
    conn net.Conn
)
var (
    doingJobCounter int32
)

func init() {
    flags.StringVar(&baseDir, "dir", "/tmp/tmp", "worker base dir")
    flags.StringVar(&workerId, "id", "yy-168.168.0.123", "worker id")
    flags.StringVar(&masterAddr, "addr", "168.168.60.15:20001", "master address")
    flags.StringVar(&httpPort, "port", "20002", "http port")
}

func Start(args []string) {
    flags.Parse(args)

//    err := delInvalidServers(baseDir, nginxCfgDir)
//    if err != nil {
//        log.Fatal(err)
//    }

    if isIdIllegal(workerId) {
        log.Fatal("workerId非法，合法字符: a-z A-Z 0-9 - _ .")
    }

    go httpService("127.0.0.1:" + httpPort)
//    go tickResClear()

    conn = connectUntilSuccess()
//    go tickUpdateWorkerStat(conn)

//    atomic.StoreInt32(&doingJobCounter, 0)

    for {
        data, err := msg.ReadBytes(conn)
        logInfo(data)
        if err != nil {
            conn.Close()
            logError("worker msg.ReadBytes error, reconnect!", err)
            conn = connectUntilSuccess()
//            go tickUpdateWorkerStat(conn)
            continue
        }
        fmt.Println("recv job...")
//        atomic.AddInt32(&doingJobCounter, 1)
//        go process(conn, data)
    }
}

func isIdIllegal(id string) bool {
    for _, c := range id {
        if c != 45 && c != 46 && c != 95 && // - . _
            (c < 48 || c > 57) && // 0-9
            (c < 65 || c > 90) && // A-Z
            (c < 97 || c > 122) { // a-z
            return true
        }
    }
    return false
}

func connectUntilSuccess() net.Conn {
    for {
        c, err := connect()
        if err == nil {
            return c
        }
        time.Sleep(time.Second)
    }
}

func connect() (net.Conn, error) {
    logInfo("try connect master:", masterAddr)

    c, err := net.Dial("tcp", masterAddr)
    if err != nil {
        return nil, err
    }

    err = msg.WriteString(c, workerId)
    if err != nil {
        return nil, err
    }

    err = msg.WriteString(c, baseDir)
    if err != nil {
        return nil, err
    }

    ok, err := msg.ReadBool(c)
    if err != nil {
        return nil, err
    }

    if !ok {
        log.Fatal("worker注册失败!")
    }

    logInfo("worker连接成功!", workerId)

    return c, nil
}

//func process(c net.Conn, data []byte) {
//    defer func() {
//        atomic.AddInt32(&doingJobCounter, -1)
//    }()
//
//    jobInfo, jobId := msg.Split2(data)
//
//    logInfo("process job:", jobId)
//
//    switch jobId {
//    case msg.JOB_ID_UNION_SERVER:
//        processUnionServer(c, jobInfo)
//    default:
//        jobData, jobType := msg.Split(jobInfo)
//
//        switch jobType {
//        case msg.JOB_TYPE_CREATE_SERVER:
//            processCreateServer(c, jobId, jobData)
//        case msg.JOB_TYPE_START_SERVER:
//            processStartServer(c, jobId, jobData)
//        case msg.JOB_TYPE_STOP_SERVER:
//            processStopServer(c, jobId, jobData)
//        case msg.JOB_TYPE_UPDATE_SERVER:
//            processUpdateServer(c, jobId, jobData)
//        default:
//            msg.WriteBytes(c, msg.Assemble2([]byte("illegal jobType:"+strconv.Itoa(jobType)), jobId))
//            logError("illegal jobType:", jobType)
//        }
//    }
//}
func close() {
    if conn != nil {
        conn.Close()
    }

    for i := 0; i < 60; i++ {
        if v := atomic.LoadInt32(&doingJobCounter); v == 0 {
            break
        } else if v < 0 {
            log.Fatal("worker doingJobCounter < 0!" + strconv.Itoa(int(v)))
        } else {
            logInfo("worker wait close:", v)
        }
        time.Sleep(time.Second)
    }
    log.Fatal("worker close!")
}
