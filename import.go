package main

import (
	"io/ioutil"
	"log"
	"runtime"
	"strings"
	"time"

	"hlc/models"

	"github.com/mailru/easyjson"
)

//easyjson:json
type Users struct {
	Users []*models.User `json:"users"`
}

//easyjson:json
type Locations struct {
	Locations []*models.Location `json:"locations"`
}

//easyjson:json
type Visits struct {
	Visits []*models.Visit `json:"visits"`
}

func importDataFromPath(storage *models.Storage, dataPath string) error {
	defer func() func() {
		s := time.Now()
		log.Printf("[import] started from %s", dataPath)

		return func() {
			log.Printf("[import] finished in %s", time.Since(s))
		}
	}()()

	dir, err := ioutil.ReadDir(dataPath)
	if err != nil {
		return err
	}

	var u Users
	for _, f := range dir {
		if strings.HasPrefix(f.Name(), "users_") {
			if err := ReadAndUnmarshalFile(dataPath+"/"+f.Name(), &u); err != nil {
				return err
			}

			storage.SetUser(u.Users...)
		}
	}
	u = Users{}
	runtime.GC()

	var l Locations
	for _, f := range dir {
		if strings.HasPrefix(f.Name(), "locations_") {
			if err := ReadAndUnmarshalFile(dataPath+"/"+f.Name(), &l); err != nil {
				return err
			}

			storage.SetLocation(l.Locations...)
		}
	}
	l = Locations{}
	runtime.GC()

	var v Visits
	for _, f := range dir {
		if strings.HasPrefix(f.Name(), "visits_") {
			if err := ReadAndUnmarshalFile(dataPath+"/"+f.Name(), &v); err != nil {
				return err
			}

			storage.SetVisit(v.Visits...)
		}
	}
	v = Visits{}
	runtime.GC()

	return nil
}

func ReadAndUnmarshalFile(filename string, v easyjson.Unmarshaler) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return easyjson.Unmarshal(data, v)
}
