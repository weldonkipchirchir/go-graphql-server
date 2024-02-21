package repo

import (
	"context"
	"os"
	"time"

	"github.com/weldonkipchirchir/go-graphql-server/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VideoRepo interface {
	Save(video *model.Video)
	FindAll() []*model.Video
}

type Database struct {
	Client *mongo.Client
}

func New() VideoRepo {

	MONGODB := os.Getenv("MONGODB")

	clientOptions := options.Client().ApplyURI(MONGODB)
	// clientOptions := options.Client().ApplyURI("mongodb+srv://weldon:Fh7WhrXt8Tvu6Iky@cluster0.5n6wuao.mongodb.net/go-graphql?retryWrites=true&w=majority")

	clientOptions = clientOptions.SetMaxPoolSize(50)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	dbclient, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		panic(err)
	}
	err = dbclient.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	return &Database{
		Client: dbclient,
	}
}

const (
	databaseName   = "graphql"
	collectionName = "videos"
)

func (d *Database) Save(video *model.Video) {
	collection := d.Client.Database(databaseName).Collection(collectionName)
	_, err := collection.InsertOne(context.TODO(), video)
	if err != nil {
		panic(err)
	}
}

func (d *Database) FindAll() []*model.Video {
	collection := d.Client.Database(databaseName).Collection(collectionName)
	var videos []*model.Video
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var video model.Video
		err := cursor.Decode(&video)
		if err != nil {
			panic(err)
		}
		videos = append(videos, &video)
	}
	return videos
}
