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
	databaseName   = "mbws-db"
	collectionName = "users"
)

// Dao provides access to database
type Dao interface {
	Connect() error
	Insert(interface{}) error
	FindUser(string) (*User, error)
	Find(bson.M) (interface{}, error)
	UpdateUser(*User, ...string) error
	UpdateUserStock(*User) error
	Update(bson.M, bson.D, ...*options.UpdateOptions) (interface{}, error)
	Delete(bson.M)
}

// UserDao is a Dao for user db
type UserDao struct {
	URL            string
	DBName         string
	CollectionName string
	Client         *mongo.Client
}

func checkConnect(c *mongo.Client) error {
	if c == nil {
		return errors.New("Not connected")
	}
	return nil
}

func getUserCollection(
	c *mongo.Client, dbname, colname string) (
	*mongo.Collection, error) {
	if err := checkConnect(c); err != nil {
		return nil, err
	}
	return c.Database(dbname).Collection(colname), nil
}

// Connect to database
func (ud *UserDao) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(ud.URL))
	ud.Client = client
	return err
}

// Insert a user
func (ud *UserDao) Insert(u interface{}) error {
	collection, err := getUserCollection(ud.Client, ud.DBName, ud.CollectionName)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, u)
	id := res.InsertedID
	log.Println(id)
	return err
}

// FindUser finds a user from db by email
func (ud *UserDao) FindUser(email string) (*User, error) {
	filter := bson.M{"email": email}
	user, err := ud.Find(filter)
	if err != nil {
		return nil, err
	}
	return user.(*User), err
}

// Find document from collection
func (ud *UserDao) Find(filter bson.M) (interface{}, error) {
	collection, err := getUserCollection(ud.Client, ud.DBName, ud.CollectionName)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates user info
func (ud *UserDao) UpdateUser(u *User, updateKeys ...string) error {
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
		_, err := ud.Update(filter, update)
		return err
	}

	return nil
}

// UpdateUserStock add user's new stock info
func (ud *UserDao) UpdateUserStock(u *User) error {
	filter := bson.M{
		"email": u.Email,
	}

	update := bson.D{{Key: "$set", Value: bson.M{
		"stocks": u.Stocks}}}

	_, err := ud.Update(filter, update)
	return err
}

// Update document from collection of db
func (ud *UserDao) Update(filter bson.M, update bson.D,
	opts ...*options.UpdateOptions) (interface{}, error) {
	collection, err := getUserCollection(ud.Client, ud.DBName, ud.CollectionName)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return collection.UpdateOne(ctx, filter, update, opts...)
}

// Delete document from collection of db
func (ud *UserDao) Delete(filter bson.M) {

}
