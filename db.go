package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//CountryMongoDB is a struct
type CountryMongoDB struct {
	URL                   string
	Name                  string
	CountryCollectionName string
}

func setupsDB() *CountryMongoDB {

	db := CountryMongoDB{
		"mongodb://kiwi:Butikker1@ds143754.mlab.com:43754/country",
		"country",
		"countryNumbers",
	}

	session, err := mgo.Dial(db.URL)
	fmt.Print("ddd")
	if err != nil {
		//fmt.Print(err.Error())
	}
	fmt.Print("@@@")
	defer session.Close()
	fmt.Print("@@@")
	return &db
}

func tearDownsDB(db *CountryMongoDB) {
	session, err := mgo.Dial(db.URL)
	defer session.Close()
	if err != nil {
		fmt.Print("error in teardown")
	}
	err = session.DB(db.Name).DropDatabase()
	if err != nil {

	}
}

//Init initializes
func (db *CountryMongoDB) Init() {
	fmt.Print("kake")
	session, err := mgo.Dial(db.URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = session.DB(db.Name).C(db.CountryCollectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

}

//Get gets
func (db *CountryMongoDB) Get(keyID string) (PopDb, bool) {
	session, err := mgo.Dial(db.URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	tr := PopDb{}
	allWasGood := true

	err = session.DB(db.Name).C(db.CountryCollectionName).Find(bson.M{"name": keyID}).One(&tr)
	if err != nil {
		allWasGood = false
	}

	fmt.Printf(tr.Name)
	return tr, allWasGood
}

//Add adds
func (db CountryMongoDB) Add(p PopDb) {
	session, err := mgo.Dial(db.URL)
	if err != nil {

		panic(err)
	}
	defer session.Close()

	err = session.DB(db.Name).C(db.CountryCollectionName).Insert(p)
	if err != nil {
		fmt.Printf("error in insert %v", err.Error())
	}
}

func (db *CountryMongoDB) GetAll() []PopDb {
	session, err := mgo.Dial(db.URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	var all []PopDb
	session.DB(db.Name).C(db.CountryCollectionName).Find(bson.M{}).All(&all)
	return all

}

func (db CountryMongoDB) Length() int {
	session, err := mgo.Dial(db.URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count, err := session.DB(db.Name).C(db.CountryCollectionName).Count()
	if err != nil {
		fmt.Printf("error in count()")
		return -1
	}
	return count

}
func (db CountryMongoDB) Remove(n string) {
	session, err := mgo.Dial(db.URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.DB(db.Name).C(db.CountryCollectionName).Remove(bson.M{"name": n})

}
