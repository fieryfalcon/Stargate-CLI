package models

type RocketLaunch struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Provider          Provider
	Vehicle           Vehicle
	Pad               Pad
	Missions          []Mission
	LaunchDescription string `json:"launch_description"`
	WinOpen           string `json:"win_open"`
	T0                string `json:"t0"`
	DateStr           string `json:"date_str"`
	Slug              string `json:"slug"`
}

type Provider struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Vehicle struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Pad struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location Location
}

type Location struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Country string `json:"country"`
}

type Mission struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RocketLaunchResponse struct {
	Count  int            `json:"count"`
	Result []RocketLaunch `json:"result"`
}
