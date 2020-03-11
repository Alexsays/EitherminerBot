package ethermine

import (
	"encoding/json"
	"net/http"
	"time"
)

type DashboardWorker struct {
	Worker string
	Time uint32
	LastSeen uint32
	ReportedHashrate float64
	CurrentHashrate float64
	ValidShares uint32
	InvalidShares uint32
	StaleShares uint32
}

type DashboardData struct {
	Workers []DashboardWorker
}

type DashboardResponse struct {
	Status string
	Data DashboardData
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func GetJson(url string, target interface{}) error {
    r, err := myClient.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()

    return json.NewDecoder(r.Body).Decode(target)
}

func GetDashboard() DashboardResponse {
	dashboardResp := DashboardResponse{}
	err := GetJson("https://api.ethermine.org/miner/0x746B683fD19526aA735d44352B3655e0439b516d/dashboard", &dashboardResp)
	if err != nil {
		panic(err)
	}

	return dashboardResp
}
