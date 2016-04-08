package server

import (
    "flag"
    "mops/server/front"
    "log"
)

var (
    flags              = flag.NewFlagSet("Server", flag.ExitOnError)
    portClientListener string
    portHttp           string
)

func init() {
    flags.StringVar(&portClientListener, "clientlistener", "20001", "server listen port")
    flags.StringVar(&portHttp, "ph", "127.0.0.1:20000", "server http port ")
    flags.StringVar(&fileCfg, "c", "cfg.ini", "master platform config")
}

func Start(args []string) {
    flags.Parse(args)

    err := loadCfg()

    if err != nil {
        log.Fatal(err)
    }

    go workerManager("0.0.0.0:" + portClientListener)
    front.HttpServer(portHttp)
}

