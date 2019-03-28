package main

import (
	"bytes"
	"conky/weather/config"
	"conky/weather/openWeatherMap"
	"fmt"
	"log"
	"sync"
	"time"
)

var layaot = "2006-01-02 15:04:05"

func main1() {
	apiWeather := openWeatherMap.New(config.GetConfig().ApiKey, config.GetConfig().CityID, "2.5")
	apiWeather.Debug = false
	wg := sync.WaitGroup{}

	weather := new(openWeatherMap.Weather)

	forecast := new(openWeatherMap.Forecast)

	wg.Add(2)
	go func(weather *openWeatherMap.Weather) {
		defer wg.Done()

		var err error
		*weather, err = apiWeather.GetWeather()
		if err != nil {
			log.Println(err)
		}
	}(weather)

	go func(forecast *openWeatherMap.Forecast) {
		defer wg.Done()
		var err error
		*forecast, err = apiWeather.GetForecast()
		if err != nil {
			log.Println(err)
		}
	}(forecast)

	wg.Wait()

	//dayHour := strconv.Itoa(time.Now().Truncate(24 * time.Hour).Add(9 * time.Hour).UTC().Hour()) + ":00:00"

	days := [3]struct {
		icon      string
		condition string
		high      float64
		low       float64
		time      time.Time
	}{}

	for i := range days {
		days[i].time = time.Now().UTC().Truncate(24 * time.Hour).AddDate(0, 0, i+1)

		days[i].high = -50
		days[i].low = 50
	}

	{
		dayID := 0

		for _, row := range forecast.List {

			currentDate, err := time.Parse(layaot, row.DtTxt)
			if err != nil {
				log.Printf("Error parse date: %v", err)
			}

			if currentDate.Before(days[0].time) {
				continue
			}

			if currentDate.Format("2006-01-02") != days[dayID].time.Format("2006-01-02") {
				dayID++
			}

			if dayID > len(days)-1 {
				break
			}

			if currentDate.After(days[dayID].time.Add(6 * time.Hour)) {
				if days[dayID].low > row.Main.TempMin {
					days[dayID].icon = row.Weather[0].Icon
					days[dayID].condition = row.Weather[0].Main
					days[dayID].low = row.Main.TempMin
				}
			} else {
				if days[dayID].high < row.Main.TempMax {
					days[dayID].icon = row.Weather[0].Icon
					days[dayID].condition = row.Weather[0].Main
					days[dayID].high = row.Main.TempMax
				}

			}

		}

	}

	tmplWeather := "${color0}Weather %v ${hr 2}\n"
	tmplSkyCurrent := "${goto 20}${color0}Sky ${goto 130}${color gray60}%v\n"
	tmplTemperatureCurrent := "${goto 20}${color0}Temperature ${goto 130}${color gray60}%v\n"
	tmplHumidityCurrent := "${goto 20}${color0}Humidity ${goto 130}${color gray60}%v\n"
	tmplVisibilityCurrent := "${goto 20}${color0}Visibility ${goto 130}${color gray60}%v\n"
	tmplDays := "${goto 25}${color gray60}%v${goto 130}%v${goto 235}%v\n\n"
	tmplPicDays := "${image $HOME/.conky/weather_icons/%v.png -p %v,167 -s 75x75}${image $HOME/.conky/weather_icons/%v.png -p %v,167 -s 75x75}${image $HOME/.conky/weather_icons/%v.png -p %v,167 -s 75x75} \n\n\n"
	tmplDaysTemp := "${goto 20}${color gray60}%.1f/%.1f°C${goto 135}%.1f/%.1f°C${goto 240}%.1f/%.1f°C"
	buf := bytes.Buffer{}

	buf.WriteString(fmt.Sprintf(tmplWeather, weather.Name))
	buf.WriteString(fmt.Sprintf(tmplSkyCurrent, weather.Weather[0].Main))
	buf.WriteString(fmt.Sprintf(tmplTemperatureCurrent, weather.Main.Temp))
	buf.WriteString(fmt.Sprintf(tmplHumidityCurrent, weather.Main.Humidity))
	buf.WriteString(fmt.Sprintf(tmplVisibilityCurrent, weather.Visibility/1000))
	buf.WriteString(fmt.Sprintf(tmplDays, days[0].time.Weekday().String(), days[1].time.Weekday().String(), days[2].time.Weekday().String()))
	buf.WriteString(fmt.Sprintf(tmplPicDays, days[0].icon, 10, days[1].icon, 115, days[2].icon, 220))
	buf.WriteString(fmt.Sprintf(tmplDaysTemp, days[0].high, days[0].low, days[1].high, days[1].low, days[2].high, days[2].low))

	fmt.Println(buf.String())

}
