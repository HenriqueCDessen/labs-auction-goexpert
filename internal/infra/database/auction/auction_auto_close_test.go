package auction_test

import (
	"context"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestMongo(t *testing.T) *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://admin:admin@localhost:27017/auctions?authSource=admin")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		t.Fatalf("failed to ping MongoDB: %v", err)
	}

	return client.Database("auctions")
}

func TestAuctionAutoClose(t *testing.T) {
	os.Setenv("AUCTION_DURATION_SECONDS", "2")

	db := setupTestMongo(t)
	repo := auction.NewAuctionRepository(db)

	auctionId := uuid.New().String()
	now := time.Now().UTC()

	auctionEntity := &auction_entity.Auction{
		Id:          auctionId,
		ProductName: "Teste",
		Category:    "Eletrônicos",
		Description: "Teste fechamento",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   now,
	}

	err := repo.CreateAuction(context.TODO(), auctionEntity)
	assert.Nil(t, err)

	time.Sleep(3 * time.Second)

	var result auction.AuctionEntityMongo
	errMongo := repo.Collection.FindOne(context.TODO(), map[string]interface{}{"_id": auctionId}).Decode(&result)
	assert.Nil(t, errMongo)

	assert.Equal(t, auction_entity.Completed, result.Status)
}
