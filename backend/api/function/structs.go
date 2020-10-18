package main

type Location struct {
	LocationID   string `json:"locationID"`
	LocationName  string `json:"locationName"`
	Description  string `json:"description"`
	latitude float32 `json:"latitude"`
	longitude float32 `json:"longitude"`
}
