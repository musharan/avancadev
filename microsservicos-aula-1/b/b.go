package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Result struct {
	Status string
}

func main() {

	http.HandleFunc("/", home)
	http.ListenAndServe(":9091", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Println(r.FormValue("coupon"))
	log.Println(r.FormValue("ccNumber"))
	log.Println(r.FormValue("ccCvv"))

	coupon := r.PostFormValue("coupon")
	ccNumber := r.PostFormValue("ccNumber")
	ccCvv := r.PostFormValue("ccCvv")

	resultCoupon := makeHttpCall("http://localhost:9092", coupon, ccCvv)

	result := Result{Status: "declined"}

	if ccNumber == "1" {
		result.Status = "approved"
	}

	if resultCoupon.Status == "invalid" {
		result.Status = "invalid coupon"
	}

	if resultCoupon.Status == "cvvError" {
		result.Status = "invalid cvv"
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error processing json")
	}

	fmt.Fprintf(w, string(jsonData))
}

func makeHttpCall(urlMicroservice string, coupon string, ccCvv string) Result {

	values := url.Values{}
	values.Add("coupon", coupon)
	values.Add("ccCvv", ccCvv)

	res, err := http.PostForm(urlMicroservice, values)
	if err != nil {
		result := Result{Status: "Servidor fora do ar!"}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result

}
