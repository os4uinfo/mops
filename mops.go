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
           var v []string
           v = []string{"a", "b"}
           server.Start(v)
       case "startclient":
           var v []string
           v = []string{"a", "b"}
           client.Start(v)
       default:
           fmt.Println("Usage!")
       }
    }
}
