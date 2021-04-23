package coinbase

type PriceData struct {
	Base     string
	Currency string
	Amount   string
}

type PriceResponse struct {
	Data PriceData
}

func GetPrice() PriceResponse {
	resp := PriceResponse{}
	err := GetJson(BaseURL+"prices/ETH-EUR/spot", &resp)
	if err != nil {
		panic(err)
	}

	return resp
}
