package handlers

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/alisilver78/goWebAPI/config"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	c        *mongo.Client
	db       *mongo.Database
	col      *mongo.Collection
	usersCol *mongo.Collection
	cfg      config.Properties
	h        ProductHandler
	uh       UsersHandler
)

func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Configuration cannot be read: %v", err)
	}
	connectURI := fmt.Sprintf("mongodb://%s:%s", cfg.DBHost, cfg.DPPort)
	var err error
	c, err = mongo.Connect(context.Background(), options.Client().ApplyURI(connectURI))
	if err != nil {
		log.Fatalf("Unable connect to database: %v", err)
	}

	db = c.Database(cfg.DBName)
	col = db.Collection(cfg.ProductsCollection)
	usersCol = db.Collection(cfg.UsersCollection)
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	testCode := m.Run()
	//after test
	col.Drop(ctx)
	usersCol.Drop(ctx)
	db.Drop(ctx)
	os.Exit(testCode)
}
