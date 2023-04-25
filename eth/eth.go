package eth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
)

type EthTransaction struct {
	Hash  string `json:"hash"`
	Value string `json:"value"`
	Time  int64  `json:"timeStamp"`
	To    string `json:"to"`
	From  string `json:"from"`
}

func GetEthTransaction(txAddress string, rxAddress string, reqAmount int64, reqTimestamp int64) (bool, error) {
	// Make API request to get Ethereum transaction data
	resp, err := http.Get(fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=YourApiKeyToken", txAddress))
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// Parse API response
	var txList []EthTransaction
	err = json.Unmarshal(body, &txList)
	if err != nil {
		return false, err
	}

	// Search for matching transaction
	for _, tx := range txList {
		if tx.Time >= reqTimestamp {
			value := parseEthValue(tx.Value)
			if tx.To == rxAddress && value.Cmp(big.NewFloat(float64(reqAmount))) == 0 {
				fmt.Println("Transaction found: ", tx.Hash)
				return true, nil
			}
		}
	}

	// No matching transaction found
	fmt.Println("Transaction not found.")
	return false, nil
}

// Helper function to convert Ethereum value from wei to ether
func parseEthValue(value string) *big.Float {
	weiValue, _ := new(big.Int).SetString(value, 0)
	ethValue := new(big.Float).Quo(new(big.Float).SetInt(weiValue), big.NewFloat(1e18))
	return ethValue
}
