package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

var db *CountryMongoDB
var firstRegRoot = "https://restcountries.eu/rest/v2/name/"
var secondRegRoot = "?fields=region;name"

var firstBordRoot = "https://restcountries.eu/rest/v2/name/"
var secondBordRoot = "?fields=borders;name;alpha3Code"

var firstcRoot = "https://restcountries.eu/rest/v2/name/"
var secondcRoot = "?fields=currencies"

var firstRoot = "https://restcountries.eu/rest/v2/name/"
var popField = "?fields=population;name"

var firstcConvert = "https://free.currencyconverterapi.com/api/v6/convert?q="
var secondcConvert = "&compact=ultra"

type Region struct {
	Reg  string `json:"region"`
	Name string `json:"name"`
}

type Border struct {
	Borders []string `json:"borders"`
	Name    string   `json:"name"`
	Code    string   `json:"alpha3Code"`
}

type Population struct {
	Pop  int    `json:"population"`
	Name string `json:"name"`
}

type PopDb struct {
	Id   bson.ObjectId `bson:"_id,omitempty"`
	Pop  int           `json:"population"`
	Name string        `json:"name"`
}

/*func handlerCurrency(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 6 {
		http.Error(w, "Please Make sure you have written two countries and included '/' in the end.", http.StatusBadRequest)
		return
	}
	c1 := parts[3]
	c2 := parts[4]
	if c1 == c2 {
		http.Error(w, "Please choose two different countries.", http.StatusBadRequest)
		return
	}

	apiRoot := firstcConvert + c1 + "_" + c2 + secondcConvert

	response, err := http.Get(apiRoot)
	if err != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}
	//fmt.Print(string(body))
	var bord []float64
	fmt.Println(apiRoot)
	err = json.Unmarshal(body, &bord)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, bord)

}*/

func handlerPopulation(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) == 5 && parts[3] == "ranks" {
		if db.Length() == 0 {
			fmt.Fprint(w, "There are currently no countries added to the database.")
			return
		}
		temp := db.GetAll()
		var highestName string
		highestPop := 0
		for _, s := range temp {
			if s.Pop > highestPop {
				highestPop = s.Pop
				highestName = s.Name
			}
		}
		fmt.Fprint(w, "The country in the database with the highest population is "+highestName+" with a population of "+strconv.Itoa(highestPop))
		return

	}
	if len(parts) == 5 && parts[3] == "remove" {

		cName := parts[4]
		_, notOk := db.Get(cName)
		if !notOk {
			http.Error(w, "Country entered does not exist in the database.", http.StatusNotFound)
			return
		}
		db.Remove(cName)
		fmt.Fprint(w, "Removed "+cName+" from the database.")
		return
	}
	if len(parts) != 6 {
		http.Error(w, "Please make sure you have written two countries and included '/' in the end.", http.StatusBadRequest)
		return
	}
	if parts[5] != "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	p1 := parts[3]
	p2 := parts[4]
	if p1 == p2 {
		http.Error(w, "Please choose two different countries.", http.StatusBadRequest)
		return
	}

	apiRoot := firstRoot + p1 + popField
	apiRoot2 := firstRoot + p2 + popField

	response, err := http.Get(apiRoot)
	if err != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}
	var pop1 []Population
	fmt.Println(apiRoot)
	err = json.Unmarshal(body, &pop1)
	if err != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}

	response2, err2 := http.Get(apiRoot2)
	if err2 != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}
	body2, err2 := ioutil.ReadAll(response2.Body)
	if err2 != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}
	var pop2 []Population
	fmt.Println(apiRoot2)
	err2 = json.Unmarshal(body2, &pop2)
	if err2 != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}

	oId := bson.NewObjectId()
	//a temp to hold the info before put in globalDB
	tempPop := PopDb{oId, pop1[0].Pop, pop1[0].Name}

	//if the id is not taken
	_, notOk := db.Get(tempPop.Name)
	if !notOk {
		db.Add(tempPop)
		fmt.Fprintln(w, "Added "+pop1[0].Name+" to the db.")
		//http.Error(w, "This country has already been added", http.StatusBadRequest)
		//return
	}

	oId2 := bson.NewObjectId()
	//a temp to hold the info before put in globalDB
	tempPop2 := PopDb{oId2, pop2[0].Pop, pop2[0].Name}

	//if the id is not taken
	_, notOk2 := db.Get(tempPop2.Name)
	if !notOk2 {
		db.Add(tempPop2)
		fmt.Fprintln(w, "Added "+pop2[0].Name+" to the db.")
	}

	if pop1[0].Pop > pop2[0].Pop {
		fmt.Fprint(w, pop1[0].Name+" has "+strconv.Itoa(pop1[0].Pop-pop2[0].Pop)+" more inhabitants than "+pop2[0].Name)
	} else {
		fmt.Fprint(w, pop2[0].Name+" has "+strconv.Itoa(pop2[0].Pop-pop1[0].Pop)+" more inhabitants than "+pop1[0].Name)
	}

}

func handlerBorder(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 6 {
		http.Error(w, "Please make sure you have written two countries and included '/' in the end.", http.StatusBadRequest)
		return
	}
	if parts[5] != "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	b1 := parts[3]
	b2 := parts[4]
	if b1 == b2 {
		http.Error(w, "Please choose two different countries.", http.StatusBadRequest)
		return
	}

	apiRoot := firstBordRoot + b1 + secondBordRoot
	apiRoot2 := firstBordRoot + b2 + secondBordRoot
	response, err := http.Get(apiRoot)
	if err != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}
	var bord []Border
	fmt.Println(apiRoot)
	err = json.Unmarshal(body, &bord)
	if err != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}

	response2, err2 := http.Get(apiRoot2)
	if err2 != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}
	body2, err2 := ioutil.ReadAll(response2.Body)
	if err2 != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}
	var bord2 []Border
	fmt.Println(apiRoot2)
	err2 = json.Unmarshal(body2, &bord2)
	if err2 != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}

	for _, s := range bord[0].Borders {
		if s == bord2[0].Code {
			fmt.Fprint(w, "The countries "+bord[0].Name+" and "+bord2[0].Name+" share borders.")
			return
		}
	}
	fmt.Fprint(w, bord[0].Name+" and "+bord2[0].Name+" do not share borders.")

}

func handlerRegion(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 6 {
		http.Error(w, "Please Make sure you have written two countries and included '/' in the end.", http.StatusBadRequest)
		return
	}
	if parts[5] != "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	c1 := parts[3]
	c2 := parts[4]
	if c1 == c2 {
		http.Error(w, "Please choose two different countries.", http.StatusBadRequest)
		return
	}

	apiRoot := firstRegRoot + c1 + secondRegRoot
	apiRoot2 := firstRegRoot + c2 + secondRegRoot
	response, err := http.Get(apiRoot)
	if err != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	var re []Region
	fmt.Println(apiRoot)
	err = json.Unmarshal(body, &re)
	if err != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}

	response2, err2 := http.Get(apiRoot2)
	if err2 != nil {
		return
	}

	body2, err2 := ioutil.ReadAll(response2.Body)
	if err2 != nil {
		return
	}

	var r2 []Region
	fmt.Println(apiRoot2)
	err2 = json.Unmarshal(body2, &r2)
	if err2 != nil {
		http.Error(w, "Please choose a valid country.", http.StatusBadRequest)
		return
	}

	if re[0].Reg == r2[0].Reg {
		fmt.Fprintln(w, "The countries "+re[0].Name+" and "+r2[0].Name+" are both in "+re[0].Reg)
	} else {
		fmt.Fprintln(w, "The countries "+re[0].Name+" and "+r2[0].Name+" are in different regions. ("+re[0].Reg+" and "+r2[0].Reg+")")
	}
}

func main() {

	//data, _ := ioutil.ReadAll(response.Body)
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("Port was not specified.")
		return
	}

	address := ":" + port

	db = setupsDB()
	defer tearDownsDB(db)
	db.Init()

	http.HandleFunc("/country/region/", handlerRegion)
	http.HandleFunc("/country/border/", handlerBorder)
	http.HandleFunc("/country/population/", handlerPopulation)
	//http.HandleFunc("/country/currency/", handlerCurrency)
	http.ListenAndServe(address, nil)

}
