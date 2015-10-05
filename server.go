package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"

	"encoding/json"

	"net/rpc"
	"strconv"
	"strings"
)

type Call struct {
	List flist `json:"list"`
}

type flist struct {
	Meta      fmeta        `json:"-"`
	Resources []fresources `json:"resources"`
}

type fmeta struct {
	Type  string `json:"-"`
	Start int32  `json:"-"`
	Count int32  `json:"-"`
}

/*

type MyResponse struct {
	List struct {
		Meta struct {
			Type  string `json:"type"`
			Start int    `json:"start"`
			Count int    `json:"count"`
		} `json:"meta"`
		Resources []struct {
			Resource struct {
				Classname string `json:"classname"`
				Fields    struct {
					Name    string `json:"name"`
					Price   string `json:"price"`
					Symbol  string `json:"symbol"`
					Ts      string `json:"ts"`
					Type    string `json:"type"`
					Utctime string `json:"utctime"`
					Volume  string `json:"volume"`
				} `json:"fields"`
			} `json:"resource"`
		} `json:"resources"`
	} `json:"list"`
}


*/
type fresources struct {
	Resource fresource `json:"resource"`
}

type fresource struct {
	Classname string  `json:"classname"`
	Fields    ffields `json:"fields"`
}

type ffields struct {
	Price  string `json:"price"`
	Symbol string `json:"symbol"`
}

type Args struct {
	Userstocksymbol string
	UserBudget      float64
}

type UpdatedResponse1 struct {
	Stocks         string  `json:"stocksymbol"`
	UnvestedAmount float64 `json:"stockprice"`
	TradeId        int     `json:"id"`
}

type SendId struct {
	TradeId int `json:"id"`
}

type UpdatedResponse2 struct {
	Stocks         string  `json:"stocksymbol"`
	UnvestedAmount float64 `json:"stockprice"`
}

type StockCalc int

var M map[int]UpdatedResponse1

func (t *StockCalc) StockPrice(args *Args, quote *UpdatedResponse1) error {

	str1 := string(args.Userstocksymbol[:])

	str1 = strings.Replace(str1, ":", ",", -1)
	str1 = strings.Replace(str1, "%", ",", -1)
	str1 = strings.Replace(str1, ",,", ",", -1)
	//
	//type url_gen struct{}

	//func (*url_gen) concate_url() error {

	//	return nil
	//
	str1 = strings.Trim(str1, " ")
	str1 = strings.Replace(str1, "\"", "", -1)
	str1 = strings.TrimSpace(str1)
	//}
	str1 = strings.TrimSuffix(str1, ",")
	finalstr1 := strings.Split(str1, ",")

	var ReqUrl string

	for i := 0; i < len(finalstr1); i++ {
		i = i + 1

		temp, _ := strconv.ParseFloat(finalstr1[i], 64)
		temp = (temp * args.UserBudget * 0.01)

		ReqUrl = ReqUrl + (finalstr1[i-1] + ",")

	}
	ReqUrl = strings.TrimSuffix(ReqUrl, ",")

	UrlStr := "http://finance.yahoo.com/webservice/v1/symbols/" + ReqUrl + "/quote?format=json"

	client := &http.Client{}

	resp, _ := client.Get(UrlStr)
	req, _ := http.NewRequest("GET", UrlStr, nil)

	req.Header.Add("If-None-Match", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// make request
	resp, _ = client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var C Call
		xx, _ := ioutil.ReadAll(resp.Body)

		err := json.Unmarshal(xx, &C)

		n := len(finalstr1)

		res := make([]float64, n, n)

		for i := 0; i < n; i++ {
			i = i + 1
			TempFloat, _ := strconv.ParseFloat(finalstr1[i], 64)
			res[i] = (TempFloat * args.UserBudget * 0.01)

		}

		var buffRead bytes.Buffer
		addCount := 0
		for _, Sample := range C.List.Resources {

			temp1 := Sample.Resource.Fields.Symbol
			temp2, _ := strconv.ParseFloat(Sample.Resource.Fields.Price, 64)
			temp3 := (int)(res[addCount+1] / temp2)

			/*

				if capacity > price {
					bought, _ := math.Modf(capacity / price)
					unvested = unvested + (capacity - (bought * price))
					fmt.Println("bought:", bought, "unves


			*/
			temp4 := math.Mod(res[addCount+1], temp2)
			addCount = addCount + 2

			quote.Stocks = fmt.Sprintf("%s:%g:%d", temp1, temp2, temp3)
			quote.UnvestedAmount = quote.UnvestedAmount + temp4
			buffRead.WriteString(quote.Stocks)
			buffRead.WriteString(",")
		}
		quote.TradeId = quote.TradeId + 1
		quote.Stocks = (buffRead.String())
		quote.Stocks = strings.TrimSuffix(quote.Stocks, ",")

		M = map[int]UpdatedResponse1{
			quote.TradeId: {quote.Stocks, quote.UnvestedAmount, quote.TradeId},
		}

		if err == nil {
			fmt.Println("")
		}
	} else {
		fmt.Println(resp.Status)

	}
	return nil
}

func (t *StockCalc) UpdStockPrice(id *SendId, upRes2 *UpdatedResponse2) error {

	var TempStr1, TempStr2 string
	var strArr []string

	var tmp = M[id.TradeId]
	TempStr2 = string(tmp.Stocks[:])
	TempStr2 = strings.Replace(TempStr2, ",", ":", -1)
	TempStr2 = strings.Trim(TempStr2, " ")

	/*
		length := s.Query.Count
		sliceLength := make([]string, length, length*2)

		for i := 0; i < s.Query.Count; i++ {
			sliceLength[i] = s.Query.Results.Quote[i].LastTradePriceOnly
			fmt.Printf(sliceLength[i])
		}
	*/
	TempStr2 = strings.TrimSpace(TempStr2)

	strArr = strings.Split(TempStr2, ":")

	for i := 0; i < len(strArr); i++ {
		TempStr1 = TempStr1 + "," + strArr[i]
		i = i + 2

	}

	TempStr1 = strings.TrimLeft(TempStr1, ",")

	UrlStr := "http://finance.yahoo.com/webservice/v1/symbols/" + TempStr1 + "/quote?format=json"

	client := &http.Client{}

	resp, _ := client.Get(UrlStr)
	req, _ := http.NewRequest("GET", UrlStr, nil)

	req.Header.Add("If-None-Match", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// make request
	resp, _ = client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var C Call

		xx, _ := ioutil.ReadAll(resp.Body)

		_ = json.Unmarshal(xx, &C)

		var buf bytes.Buffer

		k := 1
		for _, Sample := range C.List.Resources {

			temp1 := Sample.Resource.Fields.Symbol
			temp2, _ := strconv.ParseFloat(Sample.Resource.Fields.Price, 64)
			temp3, _ := strconv.ParseFloat(strArr[k], 64)

			if temp3 > temp2 {
				upRes2.Stocks = fmt.Sprintf("%s:%s%v:%v", temp1, "-", temp2, strArr[k+1])
				buf.WriteString(upRes2.Stocks)
				buf.WriteString(",")

			} else if temp3 < temp2 {

				upRes2.Stocks = fmt.Sprintf("%s:%s%v:%v", temp1, "+", temp2, strArr[k+1])
				buf.WriteString(upRes2.Stocks)
				buf.WriteString(",")

			} else {

				upRes2.Stocks = fmt.Sprintf("%s:%v:%v", temp1, temp2, strArr[k+1])
				buf.WriteString(upRes2.Stocks)
				buf.WriteString(",")
			}

			upRes2.Stocks = (buf.String())
			upRes2.Stocks = strings.TrimSuffix(upRes2.Stocks, ",")
			upRes2.UnvestedAmount = tmp.UnvestedAmount
			k = k + 3
		}
	}
	return nil
}

func main() {
	stockcalc := new(StockCalc)
	rpc.Register(stockcalc)
	rpc.HandleHTTP()

	err := http.ListenAndServe(":1331", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
