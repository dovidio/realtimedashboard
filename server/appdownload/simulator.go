package appdownload

import (
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

var Countries = [...]string{
	"Germany",
	"Italy",
	"USA",
	"Austria",
	"Switzerland",
	"Spain",
}

// GenerateData insert a random app periodically with the specified interval
func GenerateData(interval time.Duration, repository Repository) {
	for {
		time.Sleep(interval)
		var appDownload AppDownload
		appDownload.AppID = AppNames[rand.Int31n(int32(len(AppNames)))]
		appDownload.Latitude = rand.Float64()*20 + 40.0
		appDownload.Longitude = rand.Float64() * 35
		appDownload.DownloadedAt = time.Now().Unix()
		appDownload.Country = Countries[rand.Int31n(int32(len(Countries)))]

		repository.Add(appDownload)
	}
}
