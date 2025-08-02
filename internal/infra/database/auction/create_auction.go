package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {

	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	go ar.scheduleAuctionClose(auctionEntity.Id, auctionEntity.Timestamp)

	return nil
}

func getAuctionDuration() time.Duration {
	secondsStr := os.Getenv("AUCTION_DURATION_SECONDS")
	seconds, err := strconv.Atoi(secondsStr)
	if err != nil || seconds <= 0 {
		return 30 * time.Second // valor padrão
	}
	return time.Duration(seconds) * time.Second
}

func (ar *AuctionRepository) scheduleAuctionClose(auctionID string, startTime time.Time) {
	duration := getAuctionDuration()
	timeToClose := startTime.Add(duration).Sub(time.Now())
	if timeToClose <= 0 {
		timeToClose = 0
	}
	timer := time.NewTimer(timeToClose)

	<-timer.C

	err := ar.closeAuctionIfExpired(auctionID)
	if err != nil {
		logger.Error("Erro ao fechar leilão automaticamente", err)
	}
}

func (ar *AuctionRepository) closeAuctionIfExpired(auctionID string) error {
	var auction AuctionEntityMongo
	err := ar.Collection.FindOne(context.Background(), bson.M{"_id": auctionID}).Decode(&auction)
	if err != nil {
		return err
	}

	if auction.Status == auction_entity.Completed {
		return nil
	}

	now := time.Now().Unix()
	duration := getAuctionDuration()
	expirationTime := auction.Timestamp + int64(duration.Seconds())

	if now >= expirationTime {
		_, err := ar.Collection.UpdateOne(
			context.Background(),
			bson.M{"_id": auctionID},
			bson.M{"$set": bson.M{"status": auction_entity.Completed}},
		)
		return err
	}
	return nil
}
