package main

import (
	"encoding/json"
	"net/http"

	"flag"

	"github.com/bolshaaan/livecoin/coin"
)

var secretFile = flag.String("secret_file", "/Users/aleksandr/livecoin/secret", "path to file with secret")
var apiFile = flag.String("api_file", "/Users/aleksandr/livecoin/apiFile", "path to file with api key")

func main() {

	flag.Parse()

	lc := livecoin.NewLiveCoin(*secretFile, *apiFile)

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
