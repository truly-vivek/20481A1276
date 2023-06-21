package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Train struct {
	TrainName      string           `json:"trainName"`
	TrainNumber    string           `json:"trainNumber"`
	DepartureTime  Time             `json:"departureTime"`
	SeatsAvailable SeatAvailability `json:"seatsAvailable"`
	Price          TrainPrice       `json:"price"`
	DelayedBy      int              `json:"delayedBy"`
}

type Time struct {
	Hours   int `json:"Hours"`
	Minutes int `json:"Minutes"`
	Seconds int `json:"Seconds"`
}

type SeatAvailability struct {
	Sleeper int `json:"sleeper"`
	AC      int `json:"AC"`
}

type TrainPrice struct {
	Sleeper float64 `json:"sleeper"`
	AC      float64 `json:"AC"`
}

func main() {
	r := gin.Default()

	r.GET("/trains", func(c *gin.Context) {
		trainsList, err := getTrainSchedules()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, trainsList)
	})

	r.Run(":8000")
}

func getTrainSchedules() ([]Train, error) {
	apiURL := "http://104.211.219.98/train/trains"
	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODczMjgxMDAsImNvbXBhbnlOYW1lIjoiVml2ZWsgb2ZmaWNlIiwiY2xpZW50SUQiOiI5ZTExZTlhYy00NGRkLTRhYzEtYTEwNC0zZTIzYmRjOTA3NWMiLCJvd25lck5hbWUiOiIiLCJvd25lckVtYWlsIjoiIiwicm9sbE5vIjoiMjA0ODFBMTI3NiJ9.clKoUsq8GbWJsMfQz77NFDqQCnWg_YJ1H4D4Hk7wN84"
	client := &http.Client{}

	req, e := http.NewRequest("GET", apiURL, nil)
	if e != nil {
		return nil, e
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, e := client.Do(req)
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()

	var trains []Train
	e = json.NewDecoder(resp.Body).Decode(&trains)
	if e != nil {
		return nil, e
	}

	now := time.Now()
	endTime := now.Add(12 * time.Hour)
	filteredTrains := make([]Train, 0)

	for _, train := range trains {
		departureTime := time.Date(now.Year(), now.Month(), now.Day(), train.DepartureTime.Hours, train.DepartureTime.Minutes, train.DepartureTime.Seconds, 0, now.Location())
		if departureTime.After(now.Add(30*time.Minute)) && departureTime.Before(endTime) {
			filteredTrains = append(filteredTrains, train)
		}
	}

	return filteredTrains, nil
}
