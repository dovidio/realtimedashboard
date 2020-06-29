package appdownload_test

import (
	"realtimedashboard/appdownload"
	. "realtimedashboard/appdownload"
	mocks "realtimedashboard/appdownload/mock_appdownload"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetAll(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)

	dbHelper := mocks.NewMockDatabaseHelper(ctrl)
	collectionHelper := mocks.NewMockCollectionHelper(ctrl)
	resultHelper := mocks.NewMockMultiResultHelper(ctrl)

	dbHelper.EXPECT().Collection("appdownloads").Return(collectionHelper)
	collectionHelper.EXPECT().Find(gomock.Any(), gomock.Any()).Return(resultHelper, nil)
	firstCall := resultHelper.EXPECT().Next(gomock.Any()).Return(true)
	resultHelper.EXPECT().Next(gomock.Any()).After(firstCall).Return(false)
	resultHelper.EXPECT().Decode(gomock.Any())

	repository := NewMongoRepository(dbHelper)

	// when
	downloads := repository.GetAll()

	// then
	if len(downloads) != 1 {
		t.Errorf("expected size of downloads to be 1, found %d", len(downloads))
	}
	//no exception is thrown
}

func TestAdd(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)

	download := appdownload.AppDownload{
		AppID:        "some_id",
		Country:      "some_country",
		DownloadedAt: 0,
		Latitude:     0,
		Longitude:    0,
	}

	dbHelper := mocks.NewMockDatabaseHelper(ctrl)
	collectionHelper := mocks.NewMockCollectionHelper(ctrl)
	dbHelper.EXPECT().Collection("appdownloads").Return(collectionHelper)
	collectionHelper.EXPECT().InsertOne(gomock.Any(), gomock.Eq(download)).Return(nil, nil)

	repository := NewMongoRepository(dbHelper)

	// when
	err := repository.Add(download)

	// then
	if err != nil {
		t.Error("Unexpected error found", err)
	}
}
