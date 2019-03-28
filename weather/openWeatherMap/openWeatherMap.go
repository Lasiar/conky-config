package openWeatherMap

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type clouds struct {
	All int `json:"all"`
}

type Weather struct {
	Weather []weather `json:"weather"`
	Base    string    `json:"base"`
	Main    struct {
		Temp     float64 `json:"temp"`
		Pressure float64 `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  int     `json:"temp_min"`
		TempMax  int     `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed int `json:"speed"`
		Deg   int `json:"deg"`
	} `json:"wind"`
	Clouds clouds `json:"clouds"`
	Dt     int    `json:"dt"`
	Sys    struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}

type Forecast struct {
	Cod     string  `json:"cod"`
	Message float64 `json:"message"`
	Cnt     int     `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  float64 `json:"pressure"`
			SeaLevel  float64 `json:"sea_level"`
			GrndLevel float64 `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []weather `json:"weather"`
		Clouds  clouds    `json:"clouds"`
		Wind    struct {
			Speed float64 `json:"speed"`
			Deg   float64 `json:"deg"`
		} `json:"wind"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
}

func New(apiKey string, cityID int, versionAPI string) API {
	api := API{u: struct {
		url.URL
		url.Values
	}{URL: url.URL{Scheme: "https", Host: "api.openweathermap.org", Path: "data/" + versionAPI + "/"}}}

	query := api.u.Query()

	query.Set("APPID", apiKey)
	query.Set("id", strconv.Itoa(cityID))
	query.Set("units", "metric")

	api.u.Values = query

	return api
}

type API struct {
	u struct {
		url.URL
		url.Values
	}
	Debug bool
}

func (a API) GetWeather() (Weather, error) {
	a.u.Path += "weather"

	w := new(Weather)

	r, err := a.exec()
	if err != nil {
		return Weather{}, err
	}

	if err := json.NewDecoder(r).Decode(&w); err != nil {
		return Weather{}, err
	}
	return *w, nil
}

func (a API) GetForecast() (Forecast, error) {
	a.u.Path += "forecast"

	f := new(Forecast)

	r, err := a.exec()
	if err != nil {
		return Forecast{}, err
	}

	if err := json.NewDecoder(r).Decode(&f); err != nil {
		return Forecast{}, err
	}

	return *f, nil
}

func (a API) exec() (io.Reader, error) {
	a.u.URL.RawQuery = a.u.Values.Encode()
	resp, err := http.Get(a.u.String())
	if err != nil {
		return nil, err
	}

	if a.Debug {
		buf := bytes.Buffer{}
		if _, err := buf.ReadFrom(resp.Body); err != nil {
			return nil, err
		}

		log.Println(buf.String())
		return &buf, nil
	}

	return resp.Body, nil
}
