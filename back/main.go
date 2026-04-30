package main

import (
	"fmt"
	"net/http"
)


func main() {
	http.HandleFunc("/ads", getAd)	
	http.HandleFunc("/ads/activate", activateAd)

	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}