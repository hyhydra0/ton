package main 

import (
  //"fmt"
  "os"
  "log"
  "ton/app"
)



func mainCore() error {
	a := app.New()

	args := os.Args 
	if len(args) == 1 {
		log.Fatal("provide arguments")
	}

	if args[1] == "transactions" || args[1] == "balance" {
		if len(args) < 3 {
			log.Fatal("please provide wallet address and block height")
		}

		if args[1] == "transactions" {
			return a.GetTransactionsAfterBlockHeight(args[2], args[3])
		} 

		if args[1] == "balance" {
			return a.GetBalanceAtHeight(args[2], args[3])
		}
	} else if args[1] == "contract" {
		if len(args) < 2 {
			log.Fatal("please provide smart contract address")
		}

		return a.ListNominators(args[2])
	}

	return nil
}

func main() {
	if err := mainCore(); err != nil {
		log.Fatal(err)
	}
}