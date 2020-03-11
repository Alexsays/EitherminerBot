package ethermine

type SettingsData struct {
	Email     string
	Monitor   uint8
	MinPayout float64
	Ip        string
}

type SettingsResponse struct {
	Status string
	Data   SettingsData
}

func GetSettings() SettingsResponse {
	resp := SettingsResponse{}
	err := GetJson(BaseURL+MinerPATH+Wallet+"/settings", &resp)
	if err != nil {
		panic(err)
	}

	return resp
}
