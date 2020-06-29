package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// DatabaseHelper abstracts (some) mongodb Database methods
type DatabaseHelper interface {
	Collection(name string) CollectionHelper
}

// CollectionHelper abstracts (some) mongodb Collection methods
type CollectionHelper interface {
	Find(context.Context, interface{}) (MultiResultHelper, error)
	InsertOne(context.Context, interface{}) (interface{}, error)
	Watch(context.Context, mongo.Pipeline, *options.ChangeStreamOptions) (MultiResultHelper, error)
}

// MultiResultHelper abstracts (some) mongodb Cursor methods
type MultiResultHelper interface {
	Next(context.Context) bool
	Decode(v interface{}) error
}

// ClientHelper abstracts (some) mongodb Client methods
type ClientHelper interface {
	Database(string) DatabaseHelper
	StartSession() (mongo.Session, error)
}

type mongoClient struct {
	cl *mongo.Client
}
type mongoDatabase struct {
	db *mongo.Database
}
type mongoCollection struct {
	coll *mongo.Collection
}

type mongoMultiResult struct {
	cur *mongo.Cursor
}

type mongoSession struct {
	mongo.Session
}

// NewClient performs the initial connection to mongo and tries to return a client
func NewClient() ClientHelper {
	dbHostname := os.Getenv("DB_HOSTNAME")
	if dbHostname == "" {
		dbHostname = "localhost"
	}
	dbPort, err := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 64)
	if err != nil {
		dbPort = 27017
	}

	dbConnectionString := fmt.Sprintf("mongodb://%s:%d/", dbHostname, dbPort)
	fmt.Println("Trying to connect to ", dbConnectionString)

	ctxWithDeadline, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctxWithDeadline, options.Client().ApplyURI(dbConnectionString).SetDirect(true))
	if err != nil {
		log.Fatalf("Could not get client: %v", err)
	}
	if err = client.Ping(ctxWithDeadline, readpref.Nearest()); err != nil {
		log.Fatalf("failed connecting to mongo: %v", err)
	}

	return &mongoClient{cl: client}
}

func (mc *mongoClient) Database(dbName string) DatabaseHelper {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) StartSession() (mongo.Session, error) {
	session, err := mc.cl.StartSession()
	return &mongoSession{session}, err
}

func (md *mongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}) (MultiResultHelper, error) {
	cursor, err := mc.coll.Find(ctx, filter)
	return &mongoMultiResult{cur: cursor}, err
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.coll.InsertOne(ctx, document)
	return id.InsertedID, err
}

func (mc *mongoCollection) Watch(ctx context.Context, pipeline mongo.Pipeline, options *options.ChangeStreamOptions) (MultiResultHelper, error) {
	return mc.coll.Watch(ctx, pipeline, options)
}

func (mmr *mongoMultiResult) Next(ctx context.Context) bool {
	return mmr.cur.Next(ctx)
}

func (mmr *mongoMultiResult) Decode(v interface{}) error {
	return mmr.cur.Decode(v)
}
