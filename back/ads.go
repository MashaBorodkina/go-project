package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Ad struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Type string `json:"type"`
	Page string `json:"page"`
	Device string `json:"device"`
	Views int `json:"views"`
	Active bool `json:"active"`
}

func loadAds() ([]Ad, error) {
	data, err := os.ReadFile("ads.json")
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	var ads []Ad
	err = json.Unmarshal(data, &ads)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}
	return ads, nil
}

func selectedAds(ads []Ad, page string, device string) (Ad, int, error) {
	for i, ad := range ads {
		if ad.Page == page && ad.Device == device && ad.Active {
			return ad, i, nil
		}
	}
	return Ad{}, -1, fmt.Errorf("ad not found for page: %s and device: %s", page, device)
}

func getAd(w http.ResponseWriter, r *http.Request) {
	device:= r.URL.Query().Get("device")
	page := r.URL.Query().Get("page")

	if device == "" || page == "" {
		http.Error(w, "Missing 'device' or 'page' query parameter", http.StatusBadRequest)
		return
	}
	ads, err := loadAds()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, index, err := selectedAds(ads, page, device)
	if err != nil {
		http.Error(w, "Ad not found", http.StatusNotFound)
		return
	}
	ads[index].Views++

	err = saveAds(ads)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ads[index])	
	
}

func saveAds(ads []Ad) error {
	data, err := json.Marshal(ads)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}
	return os.WriteFile("ads.json", data, 0644)
}

func turnActivatorAd (ads[] Ad, id int) (int,error) {
		
		for i, ad := range ads {
		if ad.ID == id && !ad.Active {
			ads[i].Active = true
			
			err:= saveAds(ads)
			if err != nil {
				return -1, fmt.Errorf("error saving ads: %w", err)
		}
		return i, nil
	}
}
return -1, fmt.Errorf("ad with id %d not found or already active", id)
}

func activateAd(w http.ResponseWriter, r *http.Request) {
	id:= r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return		
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid 'id' query parameter", http.StatusBadRequest)
		return
	}

	ads, err := loadAds()
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
		}

	_, err = turnActivatorAd(ads, idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)	
}

	

