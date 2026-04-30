package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// Ad represents a single advertisement with targeting and tracking fields
type Ad struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Type string `json:"type"`
	Page string `json:"page"`
	Device string `json:"device"`
	Views int `json:"views"`
	Active bool `json:"active"`
}

// loadAds reads ads from ads.json and returns them as a slice of Ad
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

// selectedAds finds an active ad matching page and device and returns it with its index
func selectedAds(ads []Ad, page string, device string) (Ad, int, error) {
	for i, ad := range ads {
		if ad.Page == page && ad.Device == device && ad.Active {
			return ad, i, nil
		}
	}
	return Ad{}, -1, fmt.Errorf("ad not found for page: %s and device: %s", page, device)
}

// getAd handles HTTP requests to fetch an ad based on page and device parameters
// It also increments the view counter and saves the updated data
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

// saveAds writes the updated ads slice back to ads.json
func saveAds(ads []Ad) error {
	data, err := json.MarshalIndent(ads, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}
	return os.WriteFile("ads.json", data, 0644)
}

// turnActivatorAd activates an ad by ID if it is currently inactive
func turnActivatorAd (ads[] Ad, id int) (int, bool, error) {
		
		for i, ad := range ads {
		if ad.ID == id {
			if !ad.Active {
			ads[i].Active = true		
			 
			} else {
				ads[i].Active = false
				}
		err:= saveAds(ads)
			if err != nil {
				return -1, ads[i].Active, fmt.Errorf("error saving ads: %w", err)
				}
		return i, ads[i].Active, nil
	}
}
return -1, false, fmt.Errorf("ad with id %d not found or already active", id)
}

// activateAd handles HTTP requests to activate an ad by its ID
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

	_, isActive, err := turnActivatorAd(ads, idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response:= ""
	if isActive==true {
		response = "activated"		
	} else {
		response = "deactivated"	
	}
	

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Ad with id %d %s successfully", idInt, response)})
}

	

