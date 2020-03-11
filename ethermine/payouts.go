package ethermine

type PayoutData struct {
	PaidOn uint32
	Start  float32
	End    float32
	Amount float64
	txHash string
}

type PayoutsResponse struct {
	Status string
	Data   []PayoutData
}

func GetPayouts() PayoutsResponse {
	resp := PayoutsResponse{}
	err := GetJson(BaseURL+MinerPATH+Wallet+"/payouts", &resp)
	if err != nil {
		panic(err)
	}

	return resp
}
