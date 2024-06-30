package app 

import (
	"fmt"
	"time"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type App struct {
	client *resty.Client
}

const (
	ApiEndpoint = "https://toncenter.com/api/v3"
	SmartContractEndpoint = "https://api.ton.cat/v2/contracts"
)

func New() *App {
	a := &App{}
	a.client = resty.New()

	return a
}

func (a *App) GetTransactionsAfterBlockHeight(address string, h string) error {
	fmt.Println(fmt.Sprintf("finding transactions after block height: %s....", h))

	blockHeight, err := strconv.Atoi(h)
	if err != nil {
		return fmt.Errorf("error converting height string to number: %s", err.Error())
	}


	queryParams := map[string]string{
		"account": address,
		"limit": "256",
	}
	offset := 0

	var results Transactions


	
outer:
	for {
		var res Transactions 
		queryParams["offset"] = fmt.Sprintf("%d", offset)

		_, err := a.client.R(). 
			SetResult(&res). 
			SetQueryParams(queryParams). 
			Get(ApiEndpoint + "/transactions")

		if err != nil {
			return fmt.Errorf("error: %s", err.Error())
		}

		for _, v := range res.Transactions {
			if v.McBlockSeqno > blockHeight {
				results.Transactions = append(results.Transactions, v)
			} else {
				break outer
			}
		}

		if (len(res.Transactions) < 255) {
			break
		}

		time.Sleep(1 * time.Second)
		offset += len(res.Transactions)
	}


	fmt.Println(results.Transactions)
	return nil
}

func (a *App) GetBalanceAtHeight(address string, h string) error {
	fmt.Println(fmt.Sprintf("finding balance at block height: %s", h))

	blockHeight, err := strconv.Atoi(h)
	if err != nil {
		return fmt.Errorf("error converting height string to number: %s", err.Error())
	}

	queryParams := map[string]string{
		"account": address,
		"limit": "256",
	}
	offset := 0


	var balance string

outer:
	for {
		var res Transactions 
		queryParams["offset"] = fmt.Sprintf("%d", offset)

		_, err := a.client.R(). 
			SetResult(&res). 
			SetQueryParams(queryParams). 
			Get(ApiEndpoint + "/transactions")

		if err != nil {
			return fmt.Errorf("error: %s", err.Error())
		}

		for _, v := range res.Transactions {
			if v.McBlockSeqno <= blockHeight {
				balance = v.AccountStateAfter.Balance 
				break outer
			}
		}

		if (len(res.Transactions) == 0) {
			break
		}

		time.Sleep(1 * time.Second)
		offset += len(res.Transactions)
	}

	balFl, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		return fmt.Errorf("error getting balance: %s", err.Error())
	}


	bal := fmt.Sprintf("%.6f TON", balFl/1000000000)
	fmt.Println(bal)

	return nil
}

func (a *App) ListNominators(address string) error {
	fmt.Println(fmt.Sprintf("fetching nominators and pool data for: %s", address))

	res, err := a.client.R(). 
			Get(SmartContractEndpoint + "/nominator_pool/" + address)

	if err != nil {
		return fmt.Errorf("err: %s", err.Error())
	}

	fmt.Println(res)
	return nil
}