package main

import (
	"encoding/json"
	"net/http"

	"github.com/bolshaaan/livecoin/coin"
)

func main() {
	lc := livecoin.NewLiveCoin("/Users/aleksandr/livecoin/secret", "/Users/aleksandr/livecoin/apiFile")

	handler := http.NewServeMux()
	handler.HandleFunc("/total", func(writer http.ResponseWriter, request *http.Request) {

		total := lc.GetTotalUSD()
		//writer.Write([]byte(strconv.FormatFloat(total, 'e', 0, 64)))

		enc := json.NewEncoder(writer)
		enc.Encode(&livecoin.TotalResult{USD: total})

		request.Header.Add("Content-Type", "application/json")

		//writer.Write([]byte(fmt.Sprintf("%.4f\n", total)))
		//fmt.Println("Total = ", total)

	})

	http.ListenAndServe("localhost:8081", handler)

}
