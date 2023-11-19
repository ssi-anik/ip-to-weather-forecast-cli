package types

import "time"

type IpGeoLocation struct {
	Ip             string  `json:"ip"`
	Org            string  `json:"org"`
	Hostname       string  `json:"hostname"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	CountryCode    string  `json:"country_code"`
	CountryName    string  `json:"country_name"`
	TimezoneName   string  `json:"timezone_name"`
	ConnectionType string  `json:"connection_type"`
	CurrencyCode   string  `json:"currency_code"`
	CurrencyName   string  `json:"currency_name"`
}

func NewIpGeoLocation() *IpGeoLocation {
	return new(IpGeoLocation)
}

type Forecast struct {
	Weathers []struct {
		Timestamp     time.Time `json:"timestamp"`
		Precipitation float64   `json:"precipitation"`
		Temperature   float64   `json:"temperature"`
		WindDirection int       `json:"wind_direction"`
		WindSpeed     float64   `json:"wind_speed"`
		Condition     string    `json:"condition"`
	} `json:"weather"`
}

func NewForecast() *Forecast {
	return new(Forecast)
}
