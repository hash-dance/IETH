package db

import (
	"context"
	"fmt"
	"reflect"

	"github.com/guowenshuai/ieth/types"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var localCli *mongo.Client

func Connect(ctx context.Context, config *types.Config) (*mongo.Database, error) {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://" + config.Mongodb.Server)
	if !config.Mongodb.NoAuth {
		clientOptions.SetAuth(options.Credential{
			Username:    config.Mongodb.Username,
			Password:    config.Mongodb.Password,
			PasswordSet: false,
		})
	}
	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	localCli = client
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")
	database := client.Database(config.Mongodb.Database)
	createIndex(database)
	return database, nil
}

func Disconnect() error {
	if localCli == nil {
		return nil
	}
	return localCli.Disconnect(context.Background())
}

const (
	IpfsCollectionName  = "filecids"
	DealsCollectionName = "offlinedeals"
)

type idx struct {
	colName string
	mod     mongo.IndexModel
}

func createIndex(database *mongo.Database) error {
	todo := []idx{
		{
			IpfsCollectionName,
			mongo.IndexModel{
				Keys: bson.M{
					"payloadcid": 1,
				},
				Options: options.Index().SetUnique(true).SetName("cids"),
			},
		}, {
			DealsCollectionName,
			mongo.IndexModel{
				Keys: bson.M{
					"dealcid": 1,
				},
				Options: options.Index().SetUnique(true).SetName("deals"),
			},
		},
	}
	for _, i := range todo {
		initIndex(database, i)
	}
	return nil
}

func initIndex(database *mongo.Database, info idx) error {
	database.CreateCollection(context.Background(), info.colName)
	collection := database.Collection(info.colName)
	ctx := context.Background()

	ind, err := collection.Indexes().CreateOne(ctx, info.mod)
	// Check if the CreateOne() method returned any errors
	if err != nil {
		logrus.Errorf("Indexes().CreateOne() ERROR: %s", err)
		return err
	} else {
		// API call returns string of the index name
		logrus.Println("CreateOne() index:", ind)
		logrus.Printf("CreateOne() type: %s", reflect.TypeOf(ind))
	}
	return nil
}
