package ethermine

type RoundData struct {
	Block  uint32
	Amount float64
}

type RoundsResponse struct {
	Status string
	Data   []RoundData
}

func GetRounds() RoundsResponse {
	resp := RoundsResponse{}
	err := GetJson(BaseURL+MinerPATH+Wallet+"/rounds", &resp)
	if err != nil {
		panic(err)
	}

	return resp
}
