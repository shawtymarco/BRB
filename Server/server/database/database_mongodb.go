package database

import (
	"context"
	"errors"
	"fmt"
	"server/server/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBDatabase(uri string) (*MongoDBDatabase, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err = client.Database("brb").CreateCollection(context.TODO(), "player_data"); err != nil {
		return nil, err
	}
	return &MongoDBDatabase{client: client, cache: NewLocalDatabase()}, nil
}

type MongoDBDatabase struct {
	client *mongo.Client
	cache  *LocalDatabase
}

func (*MongoDBDatabase) String() string {
	return "MongoDB"
}

func (d *MongoDBDatabase) playerCollection() *mongo.Collection {
	return d.client.Database("brb").Collection("player_data")
}

func (d *MongoDBDatabase) CreatePlayer(data *PlayerData) error {
	_, err := d.playerCollection().InsertOne(context.TODO(), data)
	if err != nil {
		return err
	}
	return d.cache.CreatePlayer(data)
}

func (d *MongoDBDatabase) SavePlayer(data *PlayerData) error {
	res := d.playerCollection().FindOneAndReplace(context.TODO(), bson.D{{"uuid", data.Uuid}}, data)
	if res.Err() != nil {
		return res.Err()
	}
	return d.cache.SavePlayer(data)
}

func (d *MongoDBDatabase) DeletePlayerByName(playerName string, opts *PlayerNameSearchOpts) error {
	player, err := d.FindPlayerByName(playerName, opts)
	if err != nil {
		return err
	}
	_, err = d.playerCollection().DeleteOne(context.TODO(), bson.D{{"uuid", player.Uuid}})
	if err != nil {
		return err
	}
	return d.cache.DeletePlayerByName(playerName, opts)
}

func (d *MongoDBDatabase) FindPlayer(uuid uuid.UUID) (*PlayerData, error) {
	if data, err := d.cache.FindPlayer(uuid); err == nil {
		return data, nil
	}
	return d.findPlayerFromQuery(bson.D{{"uuid", uuid}}, uuid.String())
}

func (d *MongoDBDatabase) FindPlayerByDiscordID(id string) (*PlayerData, error) {
	if data, err := d.cache.FindPlayerByDiscordID(id); err == nil {
		return data, nil
	}
	return d.findPlayerFromQuery(bson.D{{"userid", id}}, id)
}

func (d *MongoDBDatabase) FindPlayerByName(playerName string, opts *PlayerNameSearchOpts) (*PlayerData, error) {
	if data, err := d.cache.FindPlayerByName(playerName, opts); err == nil {
		return data, nil
	}
	if opts == nil {
		opts = &PlayerNameSearchOpts{}
	}
	pattern := playerName
	if !opts.PartialMatch {
		pattern = fmt.Sprintf("^%v$", pattern)
	}
	query := bson.D{{"username", primitive.Regex{Pattern: pattern, Options: regexOptions(opts)}}}
	return d.findPlayerFromQuery(query, playerName)
}

func (d *MongoDBDatabase) findPlayerFromQuery(query bson.D, identifier string) (*PlayerData, error) {
	result := d.playerCollection().FindOne(context.TODO(), query)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, utils.PlayerDataNotFoundError{Identifier: identifier}
		}
		return nil, result.Err()
	}
	var data *PlayerData
	if err := result.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (d *MongoDBDatabase) SaveAll() map[string]error {
	errs := make(map[string]error)

	d.cache.mu.RLock()
	defer d.cache.mu.RUnlock()
	for _, data := range d.cache.data {
		err := d.SavePlayer(data)
		if err != nil {
			errs[data.Username] = err
		}
	}
	return errs
}

func regexOptions(o *PlayerNameSearchOpts) (str string) {
	if !o.PartialMatch {
		str += "m"
	}
	if o.CaseInsensitive {
		str += "i"
	}
	return
}
