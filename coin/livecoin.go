package livecoin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const APIUrl string = "https://api.livecoin.net"

type CurPairResult struct {
	Cur     string  `json:"cur"`
	Symbol  string  `json:"symbol"`
	Last    float64 `json:"last"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Volume  float64 `json:"volume"`
	Vwap    float64 `json:"vwap"`
	MaxBid  float64 `json:"max_bid"`
	MinAsk  float64 `json:"min_ask"`
	BestBid float64 `json:"best_bid"`
	BestAsk float64 `json:"best_ask"`
}

type BalCur struct {
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}
type BalancesRes []BalCur

type LCInterface interface {
	Balances(cur string) BalancesRes
	APIKey() string
	GetTotalUSD() float64
}

func (lc *LiveCoin) APIKey() string {
	return lc.apiKey
}

func (lc *LiveCoin) Balances(cur string) BalancesRes {
	vals := &url.Values{}
	if cur != "" {
		vals.Add("currency", cur)
	}

	resp, err := lc.doRequest(APIUrl+"/payment/balances", vals)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		expResp := make(BalancesRes, 1)
		dec := json.NewDecoder(resp.Body)

		if err := dec.Decode(&expResp); err != nil {
			panic(err)
		}

		return expResp
	} else {
		if true {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
		}

		panic(fmt.Errorf("Status code = %d", resp.StatusCode))
	}

	return nil
}

func (lc *LiveCoin) GetTotalUSD() (total float64) {
	btcUSD := lc.CurrencyBTCUSD()
	balances := lc.Balances("")

	for _, b := range balances {
		if b.Value > 0 && b.Type == "available" {
			var usd float64

			if b.Currency != "BTC" {
				pair := lc.CurrencyPair(b.Currency + "/BTC")
				usd = btcUSD.Last * pair.Last * b.Value
			} else {
				usd = btcUSD.Last * b.Value
			}

			total += usd
			fmt.Printf("%s = %f = %f \n", b.Currency, b.Value, usd)
		}
	}
	return
}

func (lc *LiveCoin) CurrencyBTCUSD() *CurPairResult {
	return lc.CurrencyPair("BTC/USD")
}

func (lc *LiveCoin) CurrencyPair(pair string) *CurPairResult {
	resp, err := http.Get(APIUrl + "/exchange/ticker?currencyPair=" + pair)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		raw := &CurPairResult{}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(raw); err != nil {
			panic(err)
		}

		return raw
	} else {
		panic(fmt.Errorf("Status code = %d", resp.StatusCode))
	}

	return nil

}

func (lc *LiveCoin) doRequest(cmdUrl string, values *url.Values) (*http.Response, error) {
	if values.Encode() != "" {
		cmdUrl += "?" + values.Encode()
	}

	req, err := http.NewRequest("GET", cmdUrl, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Api-key", lc.APIKey())

	mac := hmac.New(sha256.New, []byte(lc.secret))
	mac.Write([]byte(values.Encode()))

	sign := hex.EncodeToString(mac.Sum(nil))

	req.Header.Add("Sign", strings.ToUpper(sign))

	client := &http.Client{}
	return client.Do(req)
}

type LiveCoin struct {
	LCInterface
	secret string
	apiKey string
}

func NewLiveCoin(secFile, apiFile string) (lc LCInterface) {
	secret, err := ioutil.ReadFile(secFile)
	if err != nil {
		panic(err)
	}

	apiKey, err := ioutil.ReadFile(apiFile)

	lc = &LiveCoin{
		secret: strings.TrimSpace(string(secret)),
		apiKey: strings.TrimSpace(string(apiKey)),
	}

	return
}

type TotalResult struct {
	USD float64 `json:"usd"`
}
