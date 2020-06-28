package appdownload

import "go.mongodb.org/mongo-driver/bson/primitive"

// AppID identifies the os and type of the application
type AppID string

// AppDownload contains metadata about downloads
type AppDownload struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" mapstructure:"_id"`
	Longitude    float64            `bson:"longitude" json:"longitude"`
	Latitude     float64            `bson:"latitude" json:"latitude"`
	AppID        string             `bson:"app_id" json:"app_id" mapstructure:"app_id"`
	DownloadedAt int64              `bson:"downloaded_at" json:"downloaded_at" mapstructure:"downloaded_at"`
	Country      string             `bson:"country" json:"country" mapstructure:"country"`
}
