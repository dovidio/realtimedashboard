package appdownload

import (
	"math"
	"math/rand"
	"time"
)

// AppNames hold a list of possible app names
var AppNames = [...]string{
	"IOS_ALERT",
	"IOS_MATE",
	"IOS_E4",
	"ANDROID_ALERT",
	"ANDROID_MATE",
	"ANDOID_E4",
}

var countries = [...]string{
	"Germany",
	"Italy",
	"USA",
	"Austria",
	"Switzerland",
	"Spain",
	"Sudan",
	"Russia",
	"Australia",
	"Brasil",
	"Argentina",
}

// GenerateData insert a random app periodically with the specified interval
func GenerateData(interval time.Duration, repository Repository, quit chan struct{}) {
	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(interval)
			var appDownload AppDownload
			appDownload.AppID = AppNames[rand.Int31n(int32(len(AppNames)))]
			latitude, longitude := getRandomCoordinates()
			appDownload.Latitude = latitude
			appDownload.Longitude = longitude
			appDownload.DownloadedAt = time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
			appDownload.Country = countries[rand.Int31n(int32(len(countries)))]

			repository.Add(appDownload)
		}
	}
}

// returns random latitude and longitude
func getRandomCoordinates() (float64, float64) {

	var radius float64
	var circleX float64
	var circleY float64

	switch rand.Intn(6) {
	case 0: // Europea
		radius = 5.0
		circleX = 47.0
		circleY = 8.0
	case 1: // Asia
		radius = 18.0
		circleX = 56.0
		circleY = 86.0
	case 2: // Africa
		radius = 12.0
		circleX = 17.0
		circleY = 7.0
	case 3: // USA
		radius = 13.0
		circleX = 43.0
		circleY = -102.0
	case 4: // Australia
		radius = 8.0
		circleX = -23
		circleY = 123
	case 5: // South America
		radius = 10.0
		circleX = -11
		circleY = -58
	}

	alpha := 2 * math.Pi * rand.Float64()
	r := radius * math.Sqrt(rand.Float64())

	x := r*math.Cos(alpha) + circleX
	y := r*math.Sin(alpha) + circleY

	return x, y
}
