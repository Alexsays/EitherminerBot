package ethermine

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

const BaseURL = "https://api.ethermine.org/"
const MinerPATH = "miner/"

var Wallet = os.Getenv("WALLET_TOKEN")
var myClient = &http.Client{Timeout: 10 * time.Second}

func GetJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
