package main

import (
	"fmt"
	"log"
	"time"

	"pacpayz-blockchain/btc"
	"pacpayz-blockchain/eth"
	"pacpayz-blockchain/ltc"
)

var ticker = ""
var txAddress = ""
var rxAddress = ""
var req_amount = int64(0)
var timestamp = int64(0)

var err_starter = "error: "

func main() {
	// test the checker on a btc transaction
	ticker = "eth"
	txAddress = "randomAddress"
	rxAddress = "randomAddress"
	req_amount = int64(5000000000)
	timestamp = int64(1231006505)

	var success, err = Checker(txAddress, rxAddress, req_amount)

	if success && (err == nil) {
		log.Println("Payment Successful")
	} else {
		log.Println(err_starter, err.Error())
	}
}

// check for a transaction meeting the specifications every 5 seconds, stop checking after 5 minutes

func Checker(txAddress string, rxAddress string, req_amount int64) (bool, error) {
	// Define variables to store the start time and the elapsed time
	startTime := time.Now()
	elapsedTime := time.Since(startTime)

	// Continue checking for the transaction until it is found or 5 minutes have elapsed

	/*
		1. user enters their address, it is stored as tx_address
		3. user clicks 'submit payment'
		4. the processor starts checking if there is a transaction recieved by rx_address that meets specifications
		3. the user sends crypto to rx_address
		4. If there exists a transactions from tx_address to rx_address that meets specification, then the payment is successful
	*/

	for elapsedTime < (10 * time.Minute) {
		transaction := Tracker(txAddress, rxAddress, req_amount, timestamp)
		if transaction {
			// Transaction found, return true
			return true, nil
		}

		// Wait for 5 seconds before checking again
		time.Sleep(5 * time.Second)
		elapsedTime = time.Since(startTime)
	}

	// Transaction not found within 5 minutes, return false with an error
	return false, fmt.Errorf("Transaction not found within 5 minutes")
}

// Logic to track transactions on blockchains

func Tracker(txAddress string, rxAddress string, reqAmount int64, reqTimestamp int64) bool {

	var transaction bool

	log.Println("Tracking transaction...")
	// if the ticker is 'btc', then run the btc function
	if ticker == "btc" {
		var btc_bool, err = btc.GetBtcTransaction(txAddress, rxAddress, reqAmount, reqTimestamp)
		if btc_bool && (err == nil) {
			transaction = true
		} else {
			transaction = false
			log.Println(err_starter, err.Error())
		}
	} else if ticker == "eth" {
		var eth_bool, err = eth.GetEthTransaction(txAddress, rxAddress, reqAmount, reqTimestamp)
		if eth_bool && (err == nil) {
			transaction = true
		} else {
			transaction = false
			log.Println(err_starter, err.Error())
		}
	} else if ticker == "ltc" {
		var ltc_bool, err = ltc.GetLtcTransaction(txAddress, rxAddress, reqAmount, reqTimestamp)
		if ltc_bool && (err == nil) {
			transaction = true
		} else {
			transaction = false
			log.Println(err_starter, err.Error())
		}
	} else {
		// if the ticker is neither 'btc' nor 'eth', then print an error
		log.Println(err_starter, "Invalid Ticker")
		transaction = false
	}

	return transaction
}
