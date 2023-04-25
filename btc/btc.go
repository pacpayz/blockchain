package btc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type BtcTransaction struct {
	Hash string `json:"hash"`
	Out  []struct {
		Value int64  `json:"value"`
		Addr  string `json:"addr"`
	} `json:"out"`
	Time int64 `json:"time"`
}

func GetBtcTransaction(txAddress string, rxAddress string, reqAmount int64, reqTimestamp int64) (bool, error) {
	// Make API request to get Bitcoin transaction data
	resp, err := http.Get(fmt.Sprintf("https://blockchain.info/rawaddr/%s", txAddress))
	if err != nil {
		return false, fmt.Errorf("Error making request to blockchain.info: %s", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("Error reading response body: %s", err)
	}

	// Parse API response
	var txList []BtcTransaction
	err = json.Unmarshal(body, &txList)
	if err != nil {
		return false, fmt.Errorf("Error unmarshalling JSON response: %s", err)
	}

	// Search for matching transaction
	for _, tx := range txList {
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
	return false, nil
}
