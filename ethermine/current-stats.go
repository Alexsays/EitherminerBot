package ethermine

type CurrentStatsData struct {
	Time             uint32
	LastSeen         uint32
	ReportedHashrate float64
	AverageHashrate  float64
	CurrentHashrate  float64
	ValidShares      uint32
	InvalidShares    uint32
	StaleShares      uint32
	ActiveWorkers    uint16
	Unpaid           float64
	Unconfirmed      float64
	CoinsPerMin      float64
	UsdPerMin        float64
	BtcPerMin        float64
}

type CurrentStatsResponse struct {
	Status string
	Data   CurrentStatsData
}

func GetCurrentStats() CurrentStatsResponse {
	resp := CurrentStatsResponse{}
	err := GetJson(BaseURL+MinerPATH+Wallet+"/currentStats", &resp)
	if err != nil {
		panic(err)
	}

	return resp
}
