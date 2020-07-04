package main

import (
	"os"
	"testing"
)

var ud *UserDao

func TestMain(m *testing.M) {
	ud = &UserDao{
		URL: "mongodb+srv://" +
			os.Getenv("ATLAS_USER") + ":" +
			os.Getenv("ATLAS_PASS") + "@" +
			os.Getenv("ATLAS_URI"),
		DBName:         databaseName,
		CollectionName: collectionName,
		Client:         nil,
	}
	ud.Connect()
	m.Run()
}

func TestUserDao_Insert(t *testing.T) {
	type args struct {
		u interface{}
	}
	tests := []struct {
		name    string
		ud      *UserDao
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Insert test",
			ud:   ud,
			args: args{
				User{
					Name:      "test name",
					Email:     "test email",
					AvatarURL: "test url",
					Stocks: map[string]Stock{
						"KRX 123456": {
							Name: "test stock",
							Buyed: []Transaction{
								{Price: 12000, Quantity: 5},
								{Price: 14000, Quantity: 2},
							},
						},
						"KRX 111111": {
							Name: "test stock2",
							Buyed: []Transaction{
								{Price: 2000, Quantity: 20},
								{Price: 3000, Quantity: 15},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ud.Insert(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("UserDao.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserDao_UpdateUserStock(t *testing.T) {
	type args struct {
		u *User
	}
	tests := []struct {
		name    string
		ud      *UserDao
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Update Stock test",
			ud:   ud,
			args: args{
				&User{
					Name:      "Wooseok Son",
					Email:     "aiceru@gmail.com",
					AvatarURL: "https://lh3.googleusercontent.com/a-/AOh14GgW6EQt1TEkfEVQnOr66MqEXoeBxK2mi5hecvmIenU",
					Stocks: map[string]Stock{
						"KRX 123456": {
							Name: "test stock",
							Buyed: []Transaction{
								{Price: 12000, Quantity: 5},
							},
						},
						"KRX 111111": {
							Name: "test stock2",
							Buyed: []Transaction{
								{Price: 2000, Quantity: 20},
								{Price: 3000, Quantity: 15},
								{Price: 1000, Quantity: 15},
							},
						},
						"KRX 000000": {
							Name: "test stock3",
							Buyed: []Transaction{
								{Price: 2000, Quantity: 20},
								{Price: 3000, Quantity: 15},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ud.UpdateUserStock(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("UserDao.UpdateUserStock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
