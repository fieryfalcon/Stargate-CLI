package models

type APODResponse struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Explanation string `json:"explanation"`
}
