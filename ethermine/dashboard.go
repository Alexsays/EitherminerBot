package ethermine

type DashboardCurrentStatistic struct {
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
}

type DashboardWorker struct {
	Worker           string
	Time             uint32
	LastSeen         uint32
	ReportedHashrate float64
	CurrentHashrate  float64
	ValidShares      uint32
	InvalidShares    uint32
	StaleShares      uint32
}

type DashboardSettings struct {
	Email     string
	Monitor   uint8
	MinPayout float64
}

type DashboardStatistic struct {
	Time             uint32
	ReportedHashrate float64
	CurrentHashrate  float64
	ValidShares      uint32
	InvalidShares    uint32
	StaleShares      uint32
	ActiveWorkers    uint16
}

type DashboardData struct {
	CurrentStatistics DashboardCurrentStatistic
	Settings          DashboardSettings
	Statistics        []DashboardStatistic
	Workers           []DashboardWorker
}

type DashboardResponse struct {
	Status string
	Data   DashboardData
	Error  string
}

func GetDashboard() DashboardResponse {
	resp := DashboardResponse{}
	err := GetJson(BaseURL+MinerPATH+Wallet+"/dashboard", &resp)
	if err != nil {
		panic(err)
	}

	return resp
}
