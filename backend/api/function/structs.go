package main

type UploadedData struct {
	CameraID     string  `json:"cameraID"`
	Base64Image string  `json:"base64Image"`
}

type ReturnedData struct {
	CameraID     string  `json:"cameraID"`
	Date string  `json:"date"`
	NoFaceCovering string  `json:"noFaceCovering"`
	ImproperFaceCovering string  `json:"improperFaceCovering"`
	Total string  `json:"totalSeen"`
}
