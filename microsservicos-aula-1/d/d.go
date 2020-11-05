package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CcCvv struct {
	Code string
}

type CcCvvs struct {
	CcCvv []CcCvv
}

func (c CcCvvs) Check(code string) string {
	for _, item := range c.CcCvv {
		if code == item.Code {
			return "valid"
		}
	}
	return "invalid"
}

type Result struct {
	Status string
}

var ccCvvs CcCvvs

func main() {
	ccCvv := CcCvv{
		Code: "1010",
	}

	ccCvvs.CcCvv = append(ccCvvs.CcCvv, ccCvv)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9093", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Println(r.FormValue("ccCvv"))
	ccCvv := r.PostFormValue("ccCvv")
	valid := ccCvvs.Check(ccCvv)

	result := Result{Status: valid}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))

}
