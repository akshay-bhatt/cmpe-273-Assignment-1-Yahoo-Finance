package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
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

type Args1 struct {
	Server_side_symbol string
	Server_side_budget float64
}

type Updated_Response1 struct {
	Stocks         string  `json:"stocksymbol"`
	UnvestedAmount float64 `json:"stockprice"`
	TradeId        int     `json:"id"`
}

type Updated_Id struct {
	TradeId int `json:"id"`
}

type Updated_Res2 struct {
	Stocks         string  `json:"stocksymbol"`
	UnvestedAmount float64 `json:"stockprice"`
}

type Calc_stock int

var StoreResp map[int]Updated_Response1

func (t *Calc_stock) Stock_price(args_ser *Args1, updated_Res1 *Updated_Response1) error {
	fmt.Println("Hello World")
	str_temp := string(args_ser.Server_side_symbol[:])
	fmt.Println(str_temp)
	str_temp = strings.Replace(str_temp, ":", ",", -1)
	str_temp = strings.Replace(str_temp, "%", ",", -1)
	str_temp = strings.Replace(str_temp, ",,", ",", -1)
	str_temp = strings.Trim(str_temp, " ")
	str_temp = strings.Replace(str_temp, "\"", "", -1)
	str_temp = strings.TrimSpace(str_temp)
	str_temp = strings.TrimSuffix(str_temp, ",")
	final_val := strings.Split(str_temp, ",")

	var Url_send string

	for i := 0; i < len(final_val); i++ {
		i = i + 1

		temp, _ := strconv.ParseFloat(final_val[i], 64)
		temp = (temp * args_ser.Server_side_budget * 0.01)

		Url_send = Url_send + (final_val[i-1] + ",")

	}
	Url_send = strings.TrimSuffix(Url_send, ",")

	Str_url_final := "http://finance.yahoo.com/webservice/v1/symbols/" + Url_send + "/Updated_Response1?format=json"
	fmt.Println(Str_url_final)
	client_call1 := &http.Client{}

	toy, _ := client_call1.Get(Str_url_final)
	req, _ := http.NewRequest("GET", Str_url_final, nil)

	req.Header.Add("If-None-Match", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// make request
	toy, _ = client_call1.Do(req)
	if toy.StatusCode >= 200 && toy.StatusCode < 300 {
		var hi Call
		x, _ := ioutil.ReadAll(toy.Body)

		err := json.Unmarshal(x, &hi)

		n := len(final_val)

		res := make([]float64, n, n)

		for i := 0; i < n; i++ {
			i = i + 1
			temp_store, _ := strconv.ParseFloat(final_val[i], 64)
			res[i] = (temp_store * args_ser.Server_side_budget * 0.01)

		}

		var buff_obj bytes.Buffer
		add1 := 0
		for _, Sample := range hi.List.Resources {

			temp_var_1 := Sample.Resource.Fields.Symbol
			temp_var_2, _ := strconv.ParseFloat(Sample.Resource.Fields.Price, 64)
			temp_var_3 := (int)(res[add1+1] / temp_var_2)
			temp4 := math.Mod(res[add1+1], temp_var_2)
			add1 = add1 + 2

			updated_Res1.Stocks = fmt.Sprintf("%s:%g:%d", temp_var_1, temp_var_2, temp_var_3)
			updated_Res1.UnvestedAmount = updated_Res1.UnvestedAmount + temp4
			buff_obj.WriteString(updated_Res1.Stocks)
			buff_obj.WriteString(",")
		}
		updated_Res1.TradeId = updated_Res1.TradeId + 1
		updated_Res1.Stocks = (buff_obj.String())
		updated_Res1.Stocks = strings.TrimSuffix(updated_Res1.Stocks, ",")

		StoreResp = map[int]Updated_Response1{
			updated_Res1.TradeId: {updated_Res1.Stocks, updated_Res1.UnvestedAmount, updated_Res1.TradeId},
		}

		if err == nil {
			fmt.Println("Completed")
		}
	} else {
		fmt.Println(toy.Status)

	}
	return nil
}

func (t *Calc_stock) Updated_price_func(id *Updated_Id, var_update *Updated_Res2) error {

	var x1, y1 string
	var fin_val_store_split []string

	var tmp = StoreResp[id.TradeId]
	y1 = string(tmp.Stocks[:])
	y1 = strings.Replace(y1, ",", ":", -1)
	y1 = strings.Trim(y1, " ")
	y1 = strings.TrimSpace(y1)

	fin_val_store_split = strings.Split(y1, ":")

	for i := 0; i < len(fin_val_store_split); i++ {
		x1 = x1 + "," + fin_val_store_split[i]
		i = i + 2

	}

	x1 = strings.TrimLeft(x1, ",")

	Str_url_final := "http://finance.yahoo.com/webservice/v1/symbols/" + x1 + "/Updated_Response1?format=json"

	client_call1 := &http.Client{}

	toy, _ := client_call1.Get(Str_url_final)
	req, _ := http.NewRequest("GET", Str_url_final, nil)

	req.Header.Add("If-None-Match", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// make request
	toy, _ = client_call1.Do(req)
	if toy.StatusCode >= 200 && toy.StatusCode < 300 {
		var hi Call

		x, _ := ioutil.ReadAll(toy.Body)

		_ = json.Unmarshal(x, &hi)

		var buff_val2 bytes.Buffer

		k := 1
		for _, Sample := range hi.List.Resources {

			temp_var_1 := Sample.Resource.Fields.Symbol
			temp_var_2, _ := strconv.ParseFloat(Sample.Resource.Fields.Price, 64)
			temp_var_3, _ := strconv.ParseFloat(fin_val_store_split[k], 64)

			if temp_var_3 > temp_var_2 {
				var_update.Stocks = fmt.Sprintf("%s:%s%v:%v", temp_var_1, "-", temp_var_2, fin_val_store_split[k+1])
				buff_val2.WriteString(var_update.Stocks)
				buff_val2.WriteString(",")

			} else if temp_var_3 < temp_var_2 {

				var_update.Stocks = fmt.Sprintf("%s:%s%v:%v", temp_var_1, "+", temp_var_2, fin_val_store_split[k+1])
				buff_val2.WriteString(var_update.Stocks)
				buff_val2.WriteString(",")

			} else {

				var_update.Stocks = fmt.Sprintf("%s:%v:%v", temp_var_1, temp_var_2, fin_val_store_split[k+1])
				buff_val2.WriteString(var_update.Stocks)
				buff_val2.WriteString(",")
			}

			var_update.Stocks = (buff_val2.String())
			var_update.Stocks = strings.TrimSuffix(var_update.Stocks, ",")
			var_update.UnvestedAmount = tmp.UnvestedAmount
			k = k + 3
		}
	}
	return nil
}

func main() {
	stk := new(Calc_stock)
	rpc.Register(stk)
	rpc.HandleHTTP()

	err := http.ListenAndServe(":6611", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
