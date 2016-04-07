package front
//package main

import (
    "fmt"
    "github.com/ziutek/mymysql/mysql"
    _ "github.com/ziutek/mymysql/native" // Native engine
    "time"
)

type ServersInfo struct {
    GameName      string  `json:"gameName"`
    MergedGame    string  `json:"mergedGame"`
    ServerIP      string  `json:"serverIP"`
    GameDomain    string  `json:"gameDomain"`
    GameDBIP      string  `json:"gameDBIP"`
    GameDBPort    string  `json:"gameDBPort"`
    GameDBName    string  `json:"gameDBName"`
    WorldID       string  `json:"worldID"`
    LogServerIP   string  `json:"logServerIP"`
    LogListenPort string  `json:"logListenPort"`
    GamePort      string  `json:"gamePort"`
    PlatId        string  `json:"platId"`
    PlatName      string  `json:"platName"`
    StartTime     string  `json:"startTime"`
    MergedTime    string  `json:"mergedTime"`
    ServerState   int     `json:"serverState"`
}

type logInfo struct {
    Id       int     `json:"id"`
    Lid      string `json:"lid"`
    Logname  string `json:logname"`
    Game     string `json:"game"`
    Thread   string `json:"thread"`
    Loglevel string `json:"loglevel"`
    Logtime  string `json:"logtime"`
    Loginfo  string `json:"loginfo"`
}

const (
    defaultDbUrl = "172.17.0.1:3306"
    defaultDbUser = "gameops"
    defaultDbPwd = "work@Mokylin"
)

//func main() {
//   r := getServerInfo()
//   fmt.Println(r)
//}

func getServerInfo() []ServersInfo{
    db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
    db.Register("set names utf8")
    err := db.Connect()
    if err != nil {
        fmt.Print("connect db 失败: \n", err)
    }
    defer db.Close()
    stmt, err := db.Prepare("select * from gameops.gameinfo")
    if err != nil {
    }
    servers := []ServersInfo{}
    rows, _, err := stmt.Exec()
    for _, row := range rows {
        if row.Str(14) != "0" {
            server := ServersInfo{
                GameName     : row.Str(0),
                MergedGame   : row.Str(1),
                ServerIP     : row.Str(2),
                GameDomain   : row.Str(3),
                GameDBIP     : row.Str(4),
                GameDBPort   : row.Str(5),
                GameDBName   : row.Str(6),
                WorldID      : row.Str(7),
                LogServerIP  : row.Str(8),
                LogListenPort: row.Str(9),
                GamePort     : row.Str(10),
                PlatId       : row.Str(11),
                PlatName     : row.Str(12),
                StartTime    : time.Unix(row.Int64(13), 0).Format("2006-01-02 15:04:05"),
                MergedTime   : time.Unix(row.Int64(14), 0).Format("2006-01-02 15:04:05"),
                ServerState  : row.Int(15),
            }
            servers = append(servers, server)
        } else {
            server := ServersInfo{
                GameName     : row.Str(0),
                MergedGame   : row.Str(1),
                ServerIP     : row.Str(2),
                GameDomain   : row.Str(3),
                GameDBIP     : row.Str(4),
                GameDBPort   : row.Str(5),
                GameDBName   : row.Str(6),
                WorldID      : row.Str(7),
                LogServerIP  : row.Str(8),
                LogListenPort: row.Str(9),
                GamePort     : row.Str(10),
                PlatId       : row.Str(11),
                PlatName     : row.Str(12),
                StartTime    : time.Unix(row.Int64(13), 0).Format("2006-01-02 15:04:05"),
                MergedTime   : " ",
                ServerState  : row.Int(15),
            }
            servers = append(servers, server)
        }
    }
    return servers
}

func getLogInfo() []logInfo{
//    db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "gameops")
    db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
    db.Register("set names utf8")
    err := db.Connect()
    if err != nil {
//        log.Printf("connect db 失败: %s\n", err)
        fmt.Print("connect db 失败: \n", err)
//        Println(err)
    }
    defer db.Close()
    stmt, err := db.Prepare("select * from mysite.polls_datasave limit 10")
    if err != nil {
        fmt.Println(err)
    }
    info := []logInfo{}
    rows, _, err := stmt.Exec()
    for _, row := range rows {
        inf := logInfo{
           Id:          row.Int(0),
           Lid:         row.Str(1),
           Logname:     row.Str(2),
           Game:        row.Str(3),
           Thread:      row.Str(4),
           Loglevel:    row.Str(5),
           Logtime:     row.Str(6),
           Loginfo:     row.Str(7),
        }

        info = append(info, inf)
    }
    return info
}

