package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type Args struct {
	User_Symbol_String string
	User_Budget        float64
}

type Response1 struct {
	Stocks         string
	UnvestedAmount float64
	TradeId        int
}

type Request_Id1 struct {
	TradeId int
}

type Response2 struct {
	Stocks         string
	UnvestedAmount float64
}

// Create Client
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "server")
		os.Exit(1)
	}
	serverAddress := os.Args[1]

	client, err := rpc.DialHTTP("tcp", serverAddress+":6611")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Stock Symbol along with the amount split percentage: ")
	User_Symbol_String, _ := reader.ReadString('\n')
	//	fmt.Println(User_Symbol_String)
	fmt.Print("Enter the total budget: ")
	var User_Budget float64
	fmt.Scan(&User_Budget)

	args := Args{User_Symbol_String, User_Budget}

	var obj_res Response1
	err = client.Call("Calc_stock.Stock_price", args, &obj_res)

	if err != nil {
		log.Fatal("arith error:", err)
	}

	fmt.Println("Stock Price along with the values  :   ", obj_res.Stocks)
	fmt.Println("Unvested Amount  :   ", obj_res.UnvestedAmount)
	fmt.Println("Trade Request_Id1: ", obj_res.TradeId)

	fmt.Print("Enter TradeID to check the portfolio: ")
	var Id2 int
	fmt.Scan(&Id2)

	Request_Id1 := Request_Id1{Id2}

	var Res2 Response2

	err = client.Call("Calc_stock.Updated_price_func", Request_Id1, &Res2)

	if err != nil {
		log.Fatal("Update error:", err)
	}

	fmt.Println("Updated Stock Response: ", Res2.Stocks)
	fmt.Println("Unvested Amount: ", Res2.UnvestedAmount)

}
