package listing

// Service provides beer or review listing operations
type Service interface {
	GetAppDownloads() []beers.Beer
}
