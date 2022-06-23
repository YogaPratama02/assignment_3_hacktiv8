package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path"
	"text/template"
)

type DataFile struct {
	Water       int    `json:"water"`
	WaterStatus string `json:"water_status"`
	Wind        int    `json:"wind"`
	WindStatus  string `json:"wind_status"`
}

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	c := make(chan DataFile)
	go AutoUpdate(c)
	data := <-c

	switch {
	case data.Water <= 5:
		data.WaterStatus = "Aman"
	case data.Water >= 6 && data.Water <= 8:
		data.WaterStatus = "Siaga"
	case data.Water > 8:
		data.WaterStatus = "Bahaya"
	default:
		break
	}

	switch {
	case data.Wind <= 6:
		data.WindStatus = "Aman"
	case data.Wind >= 7 && data.Wind <= 15:
		data.WindStatus = "Siaga"
	case data.Wind > 15:
		data.WindStatus = "Bahaya"
	default:
		break
	}

	var filePath = path.Join("main.html")
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		log.Printf("Error open file with err: %s", err)
	}
	newJson, _ := json.Marshal(data)
	ioutil.WriteFile("file.json", newJson, 0644)
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AutoUpdate(c chan DataFile) {
	for {
		data := DataFile{}

		randomWater := rand.Intn(15)
		randomWind := rand.Intn(20)
		data.Water = randomWater
		data.Wind = randomWind
		c <- data
	}
}
