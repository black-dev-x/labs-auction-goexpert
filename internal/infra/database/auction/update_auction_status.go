package auction

import (
	"auction-go/configuration/logger"
	"auction-go/internal/entity/auction_entity"
	"auction-go/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
)

func (ar *AuctionRepository) UpdateAuctionStatus(id string, status auction_entity.AuctionStatus) *internal_error.InternalError {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := ar.Collection.UpdateOne(nil, filter, update)
	if err != nil {
		logger.Error("Error trying to update auction status", err)
		return internal_error.NewInternalServerError("Error trying to update auction status")
	}

	return nil
}
