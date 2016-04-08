package main

import "fmt"
import "mops/server"
import "mops/client"
import "flag"
import "runtime"
import "strings"


func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    flag.Parse()

    action := flag.Arg(0)
    if strings.HasPrefix(action, "start") {
       switch action {
       case "startserver":
           server.Start(flag.Args()[1:])
       case "startclient":
           client.Start(flag.Args()[1:])
       default:
           fmt.Println("Usage!")
       }
    }
}
