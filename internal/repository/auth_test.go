package repository

import (
	"context"
	"log"
	"testing"

	_ "acsp/internal/dto"
	"acsp/internal/model"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection

func TestAuthMongo_CreateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientOptions(
		options.
			Client().ApplyURI("mongodb+srv://yeldosmanap:1245emer@cluster0.ax3nacu.mongodb.net/?retryWrites=true&w=majority").
			SetAuth(options.Credential{
				Username: "yeldosmanap",
				Password: "1245emer",
			})).
		DatabaseName("gorest").
		CollectionName("users").
		// DatabaseName("test").
		// CollectionName("users").
		ClientType(mtest.Mock))

	defer mt.Close()

	log.Println(mt.DB.Name())
	log.Println(mt.Coll.Name())

	collection := mt.Coll
	// db := NewAuthPostgres(mt.DB)
	ctx := context.Background()

	tests := []struct {
		name    string
		mock    func()
		input   model.User
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				user := model.User{
					ID:       "",
					Name:     "Test",
					Email:    "test",
					Password: "password",
				}
				collection.FindOne(ctx, user)
				/*mock.ExpectQuery("INSERT INTO users").
				WithArgs("Test", "test", "password").WillReturnRows(rows)*/
			},
			input: model.User{
				Email:    "Test",
				Name:     "test",
				Password: "password",
			},
			want: 1,
		},
		{
			name: "Empty Fields",
			mock: func() {
				user := model.User{
					ID:       "",
					Name:     "",
					Email:    "",
					Password: "",
				}
				collection.FindOne(ctx, user)
			},
			input: model.User{
				Email:    "",
				Name:     "",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		mt.Run(testCase.name, func(t *mtest.T) {
			testCase.mock()

			got, err := db.CreateUser(ctx, testCase.input)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
		})
	}
}
