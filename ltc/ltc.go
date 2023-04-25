package ltc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type LtcTransaction struct {
	Hash string `json:"hash"`
	Out  []struct {
		Value int64  `json:"value"`
		Addr  string `json:"addr"`
	} `json:"out"`
	Time int64 `json:"time"`
}

func GetLtcTransaction(txAddress string, rxAddress string, reqAmount int64, reqTimestamp int64) (bool, error) {
	// Make API request to get Litecoin transaction data
	resp, err := http.Get(fmt.Sprintf("https://chain.so/api/v2/get_tx_unspent/LTC/%s", txAddress))
	if err != nil {
		return false, fmt.Errorf("error making API request: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse API response
	var txList struct {
		Data struct {
			Txs []LtcTransaction `json:"txs"`
		} `json:"data"`
		Status string `json:"status"`
	}
	err = json.Unmarshal(body, &txList)
	if err != nil {
		return false, fmt.Errorf("error parsing JSON response: %v", err)
	}

	if txList.Status == "fail" {
		return false, fmt.Errorf("API error: %v", txList.Data)
	}

	// Search for matching transaction
	for _, tx := range txList.Data.Txs {
		if tx.Time >= reqTimestamp {
			for _, output := range tx.Out {
				if output.Addr == rxAddress && output.Value == reqAmount {
					log.Println("Transaction found: ", tx.Hash)
					return true, nil
				}
			}
		}
	}

	// No matching transaction found
	log.Println("Transaction not found.")
	return false, nil
}
