package appdownload

import "go.mongodb.org/mongo-driver/bson/primitive"

// AppID identifies the os and type of the application
type AppID string

// AppNames hold a list of possible app names
var AppNames = [...]string{
	"IOS_ALERT",
	"IOS_MATE",
	"IOS_E4",
	"ANDROID_ALERT",
	"ANDROID_MATE",
	"ANDOID_E4",
}

// AppDownload contains metadata about downloads
type AppDownload struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Longitude    float64            `bson:"longitude" json:"longitude"`
	Latitude     float64            `bson:"latitude" json:"latitude"`
	AppID        string             `bson:"app_id" json:"app_id"`
	DownloadedAt int64              `bson:"downloaded_at" json:"downloaded_at"`
}
