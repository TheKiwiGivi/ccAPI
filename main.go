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

//database variable
var db *CountryMongoDB

//variables for various fields
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

//Region for a country's region
type Region struct {
	Reg  string `json:"region"`
	Name string `json:"name"`
}

//Border for a country's borders
type Border struct {
	Borders []string `json:"borders"`
	Name    string   `json:"name"`
	Code    string   `json:"alpha3Code"`
}

//Population for a country's population
type Population struct {
	Pop  int    `json:"population"`
	Name string `json:"name"`
}

//PopDb used for population database
type PopDb struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Pop  int           `json:"population"`
	Name string        `json:"name"`
}

func handlerPopulation(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	//if ranks is requested
	if len(parts) == 5 && parts[3] == "ranks" {
		//checks if empty
		if db.Length() == 0 {
			fmt.Fprint(w, "There are currently no countries added to the database.")
			return
		}

		temp := db.GetAll()
		var highestName string
		highestPop := 0
		//searching through all entires in the database and finds the one with highest population
		for _, s := range temp {
			if s.Pop > highestPop {
				highestPop = s.Pop
				highestName = s.Name
			}
		}
		//prints response
		fmt.Fprint(w, "The country in the database with the highest population is "+highestName+" with a population of "+strconv.Itoa(highestPop))
		return

	}
	//if remove is requested
	if len(parts) == 5 && parts[3] == "remove" {
		//checks if it actually exists in the db
		cName := parts[4]
		_, notOk := db.Get(cName)
		if !notOk {
			http.Error(w, "Country entered does not exist in the database.", http.StatusNotFound)
			return
		}
		//if so, removes it
		db.Remove(cName)
		fmt.Fprint(w, "Removed "+cName+" from the database.")
		return
	}
	//makes sure that the request has the right amount of parts. the after this point should be root/<SOMETHING>/country1/country2/
	if len(parts) != 6 {
		http.Error(w, "Please make sure you have written two countries and included '/' in the end.", http.StatusBadRequest)
		return
	}
	//if rubbish
	if parts[5] != "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//grabs the country names
	p1 := parts[3]
	p2 := parts[4]
	//makes sure the countries are not the same
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
	//makes sure the countries are not the same, this time not only input but the actual country name.
	if pop1[0].Name == pop2[0].Name {
		http.Error(w, "Please choose two different countries.", http.StatusBadRequest)
		return
	}

	oID := bson.NewObjectId()
	//a temp to hold the info before put in db
	tempPop := PopDb{oID, pop1[0].Pop, pop1[0].Name}

	//if the id is not taken
	_, notOk := db.Get(tempPop.Name)
	//adds to the database
	if !notOk {
		db.Add(tempPop)
		fmt.Fprintln(w, "Added "+pop1[0].Name+" to the db.")

	}

	oID2 := bson.NewObjectId()
	//a temp to hold the info before put in db
	tempPop2 := PopDb{oID2, pop2[0].Pop, pop2[0].Name}

	//if the id is not taken
	_, notOk2 := db.Get(tempPop2.Name)
	if !notOk2 {
		db.Add(tempPop2)
		fmt.Fprintln(w, "Added "+pop2[0].Name+" to the db.")
	}
	//checks which country has the highest population and prints out
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
	//used for getting the two countries
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
	//goes through all the borders, and if it finds a match, prints out the relevant info.
	for _, s := range bord[0].Borders {
		if s == bord2[0].Code {
			fmt.Fprint(w, "The countries "+bord[0].Name+" and "+bord2[0].Name+" share borders.")
			return
		}
	}
	//if no match
	fmt.Fprint(w, bord[0].Name+" and "+bord2[0].Name+" do not share borders.")

}

func handlerRegion(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 6 {
		http.Error(w, "Please make sure you have written two countries and included '/' in the end.", http.StatusBadRequest)
		return
	}
	//if last is not blanc, then it is rubbish
	if parts[5] != "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//grabbing countries
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
	//checks if the regions are the same
	if re[0].Reg == r2[0].Reg {
		fmt.Fprintln(w, "The countries "+re[0].Name+" and "+r2[0].Name+" are both in "+re[0].Reg)
	} else {
		fmt.Fprintln(w, "The countries "+re[0].Name+" and "+r2[0].Name+" are in different regions. ("+re[0].Reg+" and "+r2[0].Reg+")")
	}
}

func main() {

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

	http.ListenAndServe(address, nil)

}
