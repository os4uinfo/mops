package server

import (
    "encoding/json"
    "os"
    "mops/msg/pb"
    "log"
    "time"
    "io/ioutil"
)

type Cfg struct {
	PlatformInfos  []PlatformInfo  `json:"platformInfos"`
	HistoryWorkers []HistoryWorker `json:"historyWorkers"`
}

type PlatformInfo struct {
    Id     int    `json:"id"`
    Name   string `json:"name"`
    Domain string `json:"domain"`
    Cname  string `json:"cname"`
}

type HistoryWorker struct {
    Id           string `json:"id"`
    FirstRegTime string `json:"firstRegTime"`
}

var (
    fileCfg string
)

var (
	platformInfosMap  map[int]PlatformInfo
	historyWorkersMap map[string]HistoryWorker
	platformDomainMap = new(pb.DomainMap)
)

func loadCfg() error {
    platformInfosMap = make(map[int]PlatformInfo)
    historyWorkersMap = make(map[string]HistoryWorker)

    data, err := ioutil.ReadFile(fileCfg)
    if err != nil {
        if os.IsNotExist(err) {
            saveCfg()
            return nil
        }
    return err
    }

    var cfg Cfg
    err = json.Unmarshal(data, &cfg)
    if err != nil {
        return err
    }
    for _, pi := range cfg.PlatformInfos {
        if pi.Id <= 0 {
            log.Fatal("paltid must large then 0!", pi.Id)
        }
        platformInfosMap[pi.Id] = pi

        platformDomainMap.Domain = append(platformDomainMap.Domain, pi.Domain)
        platformDomainMap.OperatorId = append(platformDomainMap.OperatorId, int32(pi.Id))
    }

    for _, hw := range cfg.HistoryWorkers {
        historyWorkersMap[hw.Id] = hw
    }
    return nil
}

func saveCfg() error {
    platformInfos := []PlatformInfo{}
    historyWorkers := []HistoryWorker{}

    for _, pi := range platformInfosMap {
        platformInfos = append(platformInfos, pi)
    }

    for _, hw := range historyWorkersMap {
        historyWorkers = append(historyWorkers, hw)
    }

    cfg := Cfg{
        PlatformInfos:  platformInfos,
        HistoryWorkers: historyWorkers,
    }

    data, err := json.MarshalIndent(cfg, "", "    ")
    if err != nil {
        return err
    }

    logInfo("Generate new cfg.ini \n", string(data))

    return ioutil.WriteFile(fileCfg, data, os.ModePerm)
}

func onWorkerConnected(workerId string) {
	_, ok := historyWorkersMap[workerId]
	if ok {
		return
	}
	historyWorkersMap[workerId] = HistoryWorker{
		Id:           workerId,
		FirstRegTime: time.Now().Format("2006-01-02 03:04:05"),
	}

	err := saveCfg()
	if err != nil {
		logError("master.onWorkerConnected:", err)
	}
}
