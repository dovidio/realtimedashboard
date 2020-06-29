package appdownload_test

import (
	"fmt"
	. "realtimedashboard/appdownload"
	mocks "realtimedashboard/appdownload/mock_appdownload"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func TestRegisterUnregisterObserver(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	dbHelper := mocks.NewMockDatabaseHelper(ctrl)
	observer := mocks.NewMockObserver(ctrl)
	observers := make(map[uuid.UUID]Observer, 0)
	target := NewMongoWatchHandler(dbHelper, observers)

	// when
	uuid := target.RegisterObserver(observer)

	// then
	if len(observers) != 1 {
		t.Error("Should have one observer registered")
	}

	// when
	target.UnregisterObserver(uuid)
	if len(observers) != 0 {
		t.Error("Should have unregistered observer")
	}
}

func TestNotifiesObserversOnDatabaseChange(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	dbHelper := mocks.NewMockDatabaseHelper(ctrl)
	collectionHelper := mocks.NewMockCollectionHelper(ctrl)
	resultHelper := mocks.NewMockMultiResultHelper(ctrl)

	dbHelper.EXPECT().Collection("appdownloads").Return(collectionHelper).Times(1)
	collectionHelper.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).Return(resultHelper, nil).Times(1)
	firstCall := resultHelper.EXPECT().Next(gomock.Any()).Return(true)
	resultHelper.EXPECT().Next(gomock.Any()).After(firstCall).Return(false)
	resultHelper.EXPECT().Decode(gomock.Any()).Do(func(_arg *bson.M) {
		fmt.Println(_arg)

		*_arg = bson.M{"fullDocument": bson.M{}}

		fmt.Println(_arg)
	})

	observer1 := mocks.NewMockObserver(ctrl)
	observer2 := mocks.NewMockObserver(ctrl)
	observers := make(map[uuid.UUID]Observer, 0)
	observer1.EXPECT().OnNewAppDownload(gomock.Any())
	observer2.EXPECT().OnNewAppDownload(gomock.Any())

	target := NewMongoWatchHandler(dbHelper, observers)

	target.RegisterObserver(observer1)
	target.RegisterObserver(observer2)

	// when
	target.(*MongoDbWatchHandler).WatchAppDownloadsImpl()

}
