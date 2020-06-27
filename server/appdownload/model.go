package appdownload

import "go.mongodb.org/mongo-driver/bson/primitive"

// AppID identifies the os and type of the application
type AppID string

const (
	// IosAlert does something
	IosAlert AppID = "IOS_ALERT"
	// IosMate does something
	IosMate = "IOS_MATE"
	// IosE4 does something
	IosE4 = "IOS_E4"
	// AndroidAlert does something
	AndroidAlert = "ANDROID_ALERT"
	// AndroidMate does something
	AndroidMate = "ANDROID_MATE"
	// AndroidE4 does something
	AndroidE4 = "ANDOID_E4"
)

// AppDownloads contains metadata about downloads
type AppDownload struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Longitude    float64            `bson:"longitude" json:"longitude"`
	Latitude     float64            `bson:"latitude" json:"latitude"`
	AppID        string             `bson:"app_id" json:"app_id"`
	DownloadedAt int64              `bson:"downloaded_at" json:"downloaded_at"`
}
