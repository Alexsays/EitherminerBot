package ethermine

type HistoryData struct {
	ReportedHashrate float64
	AverageHashrate  float64
	CurrentHashrate  float64
	ValidShares      uint32
	InvalidShares    uint32
	StaleShares      uint32
	ActiveWorkers    uint32
}

type HistoryResponse struct {
	Status string
	Data   []HistoryData
}

func GetHistory() HistoryResponse {
	resp := HistoryResponse{}
	err := GetJson(BaseURL+MinerPATH+Wallet+"/history", &resp)
	if err != nil {
		panic(err)
	}

	return resp
}
