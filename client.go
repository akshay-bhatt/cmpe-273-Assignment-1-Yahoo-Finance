package main

import (
	"fmt"
	"log"
	"net/rpc"

	"bufio"
	"os"
)

type Args struct {
	Userstocksymbol string
	UserBudget      float64
}

type Response1 struct {
	Stocks         string
	UnvestedAmount float64
	TradeId        int
}

type GetId struct {
	TradeId int
}

type Response2 struct {
	Stocks         string
	UnvestedAmount float64
}

// Create Client
func main() {

	client, err := rpc.DialHTTP("tcp", "localhost:1331")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	/*

		client, err := net.Dial("tcp", "127.0.0.1:1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}
	*/
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("EEnter Stock Symbol along with the amount split percentage: ")
	Userstocksymbol, _ := reader.ReadString('\n')
	//	fmt.Println(Userstocksymbol)
	fmt.Print("Enter the total budget: ")
	var UserBudget float64
	fmt.Scan(&UserBudget)

	args := &Args{Userstocksymbol, UserBudget}

	var sendreference Response1
	err = client.Call("StockCalc.StockPrice", args, &sendreference)

	if err != nil {
		log.Fatal("OMG error:", err)
	}

	/*	// Synchronous call
		args := &Args{7, 8}
		var reply int
		c := jsonrpc.NewClient(client)
		err = c.Call("Calculator.Add", args, &reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		fmt.Printf("Result: %d+%d=%d\n", args.X, args.Y, reply)
	*/

	fmt.Println("Stock Price along with the values  :  ", sendreference.Stocks)
	fmt.Println("Unvested Amount: ", sendreference.UnvestedAmount)
	fmt.Println("Trade Request_Id1: ", sendreference.TradeId)

	fmt.Print("Enter TradeID to check the portfolio: ")
	var TradeId int
	fmt.Scan(&TradeId)

	id := GetId{TradeId}

	var sendreference2 Response2

	err = client.Call("StockCalc.UpdStockPrice", id, &sendreference2)

	if err != nil {
		log.Fatal("Update error:", err)
	}

	fmt.Println("Updated Stock Response: ", sendreference2.Stocks)
	fmt.Println("Unvested Amount: ", sendreference2.UnvestedAmount)

}
