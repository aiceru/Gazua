package main

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName = "mbws-db"
	userCollName = "users"
	corpCollName = "corps"
)

// UserDao provides access to user database
type UserDao interface {
	Connect() error
	InsertUser(*User) error
	FindUser(string) (*User, error)
	UpdateUser(*User, ...string) error
	UpdateUserStock(*User) error
}

// CorpDao provides access to corp database
type CorpDao interface {
	Connect() error
	InsertCorp(*Corp) error
	FindCorp(string) (*Corp, error)
	ReplaceCorp(*Corp) error
}

// MongoDao is a Dao for user db
type MongoDao struct {
	URL    string
	DBName string
	Client *mongo.Client
}

func checkConnect(c *mongo.Client) error {
	if c == nil {
		return errors.New("Not connected")
	}
	return nil
}

// Connect to database
func (md *MongoDao) Connect() error {
	if checkConnect(md.Client) != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(md.URL))
		md.Client = client
		return err
	}
	return nil
}

func getCollection(
	c *mongo.Client, dbname, colname string) (
	*mongo.Collection, error) {
	if err := checkConnect(c); err != nil {
		return nil, err
	}
	return c.Database(dbname).Collection(colname), nil
}

// InsertUser inserts a user
func (md *MongoDao) InsertUser(u *User) error {
	return md.Insert(userCollName, u)
}

// InsertCorp inserts a corp
func (md *MongoDao) InsertCorp(c *Corp) error {
	return md.Insert(corpCollName, c)
}

// Insert insert document to mongodb
func (md *MongoDao) Insert(colname string, d interface{}) error {
	collection, err := getCollection(md.Client, md.DBName, colname)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, d)
	id := res.InsertedID
	log.Println(id)
	return err
}

// FindUser finds a user from db by email
func (md *MongoDao) FindUser(email string) (*User, error) {
	filter := bson.M{"email": email}
	user := new(User)
	f, err := md.Find(filter, userCollName)
	if err == nil {
		f.Decode(user)
	}
	return user, err
}

// FindCorp finds a corp from db by name
func (md *MongoDao) FindCorp(name string) (*Corp, error) {
	filter := bson.M{"name": name}
	corp := new(Corp)
	f, err := md.Find(filter, corpCollName)
	if err == nil {
		f.Decode(corp)
	}
	return corp, err
}

// Find document from collection
func (md *MongoDao) Find(filter bson.M, collname string) (*mongo.SingleResult, error) {
	collection, err := getCollection(md.Client, md.DBName, collname)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return collection.FindOne(ctx, filter), nil
}

// UpdateUser updates user info
func (md *MongoDao) UpdateUser(u *User, updateKeys ...string) error {
	filter := bson.M{"email": u.Email}

	data := bson.D{}
	for _, key := range updateKeys {
		switch key {
		case "name":
			data = append(data, bson.E{Key: key, Value: u.Name})
		case "avatar_url":
			data = append(data, bson.E{Key: key, Value: u.AvatarURL})
		default:
			log.Println("Document has no key with " + key)
		}
	}

	if len(data) > 0 {
		update := bson.D{{Key: "$set", Value: data}}
		_, err := md.Update(filter, update, userCollName)
		return err
	}

	return nil
}

// UpdateUserStock add user's new stock info
func (md *MongoDao) UpdateUserStock(u *User) error {
	filter := bson.M{
		"email": u.Email,
	}

	update := bson.D{{Key: "$set", Value: bson.M{
		"stocks": u.Stocks}}}

	_, err := md.Update(filter, update, userCollName)
	return err
}

// Update document from collection of db
func (md *MongoDao) Update(filter bson.M, update bson.D, collname string,
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	collection, err := getCollection(md.Client, md.DBName, collname)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return collection.UpdateOne(ctx, filter, update, opts...)
}

// ReplaceCorp replaces corp with corp name + refresh TTL
func (md *MongoDao) ReplaceCorp(c *Corp) error {
	filter := bson.M{
		"name": c.Name,
	}

	opt := options.Replace().SetUpsert(true)

	_, err := md.ReplaceOne(filter, c, corpCollName, opt)
	return err
}

// ReplaceOne replaces one document from collection of db
func (md *MongoDao) ReplaceOne(filter bson.M, replace interface{}, collname string,
	opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	collection, err := getCollection(md.Client, md.DBName, collname)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return collection.ReplaceOne(ctx, filter, replace, opts...)
}
