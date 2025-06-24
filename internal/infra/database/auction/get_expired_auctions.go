package auction

import (
	"auction-go/configuration/logger"
	"auction-go/internal/entity/auction_entity"
	"auction-go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (ar *AuctionRepository) GetNextExpiredAuctions() ([]auction_entity.Auction, *internal_error.InternalError) {
	duration := os.Getenv("AUCTION_DURATION")
	if duration == "" {
		duration = "24h"
	}
	durationParsed, _ := time.ParseDuration(duration)
	filter := bson.M{"timestamp": bson.M{"$lt": time.Now().Add(-durationParsed)}}
	filter["status"] = auction_entity.Active

	cursor, err := ar.Collection.Find(nil, filter)
	if err != nil {
		logger.Error("Error finding expired auctions", err)
		return nil, internal_error.NewInternalServerError("Error finding expired auctions")
	}
	defer cursor.Close(nil)

	var auctionsMongo []AuctionEntityMongo

	if err := cursor.All(nil, &auctionsMongo); err != nil {
		logger.Error("Error decoding expired auctions", err)
		return nil, internal_error.NewInternalServerError("Error decoding expired auctions")
	}

	var auctionsEntity []auction_entity.Auction
	for _, auction := range auctionsMongo {
		auctionsEntity = append(auctionsEntity, auction_entity.Auction{
			Id:          auction.Id,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Status:      auction.Status,
			Description: auction.Description,
			Condition:   auction.Condition,
			Timestamp:   time.Unix(auction.Timestamp, 0),
		})
	}

	return auctionsEntity, nil
}
