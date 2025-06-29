package auction_usecase

import (
	"auction-go/configuration/logger"
	"auction-go/internal/entity/auction_entity"
	"auction-go/internal/entity/bid_entity"
	"auction-go/internal/internal_error"
	"auction-go/internal/usecase/bid_usecase"
	"context"
	"fmt"
	"sync"
	"time"
)

type AuctionInputDTO struct {
	ProductName string           `json:"product_name" binding:"required,min=1"`
	Category    string           `json:"category" binding:"required,min=2"`
	Description string           `json:"description" binding:"required,min=10,max=200"`
	Condition   ProductCondition `json:"condition" binding:"oneof=0 1 2"`
}

type AuctionOutputDTO struct {
	Id          string           `json:"id"`
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
	Status      AuctionStatus    `json:"status"`
	Timestamp   time.Time        `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type WinningInfoOutputDTO struct {
	Auction AuctionOutputDTO          `json:"auction"`
	Bid     *bid_usecase.BidOutputDTO `json:"bid,omitempty"`
}

func NewAuctionUseCase(
	auctionRepositoryInterface auction_entity.AuctionRepositoryInterface,
	bidRepositoryInterface bid_entity.BidEntityRepository) AuctionUseCaseInterface {

	logger.Info("Creating AuctionUseCase...")

	auctionUseCase := AuctionUseCase{
		auctionRepositoryInterface: auctionRepositoryInterface,
		bidRepositoryInterface:     bidRepositoryInterface,
	}
	go auctionUseCase.CheckForExpiredAuctions()
	return &auctionUseCase
}

type AuctionUseCaseInterface interface {
	CreateAuction(
		ctx context.Context,
		auctionInput AuctionInputDTO) *internal_error.InternalError

	FindAuctionById(
		ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError)

	FindAuctions(
		ctx context.Context,
		status AuctionStatus,
		category, productName string) ([]AuctionOutputDTO, *internal_error.InternalError)

	FindWinningBidByAuctionId(
		ctx context.Context,
		auctionId string) (*WinningInfoOutputDTO, *internal_error.InternalError)
}

type ProductCondition int64
type AuctionStatus int64

type AuctionUseCase struct {
	auctionRepositoryInterface auction_entity.AuctionRepositoryInterface
	bidRepositoryInterface     bid_entity.BidEntityRepository
}

func (au *AuctionUseCase) updateAuctionToCompleted(id string) {
	logger.Info("Found expired auction: " + id)
	if err := au.auctionRepositoryInterface.UpdateAuctionStatus(id, auction_entity.Completed); err != nil {
		logger.Error("Error updating auction status", err)
	}
}

func (au *AuctionUseCase) CheckForExpiredAuctions() {
	fmt.Println("Starting background job to check for expired auctions...")
	for {
		auctions, error := au.auctionRepositoryInterface.GetNextExpiredAuctions()
		if error != nil {
			logger.Error("Error getting expired auctions", error)
			return
		}
		var waitGroup sync.WaitGroup
		waitGroup.Add(len(auctions))
		for _, auction := range auctions {
			go func(id string) {
				au.updateAuctionToCompleted(id)
				waitGroup.Done()
			}(auction.Id)
		}
		waitGroup.Wait()
		<-time.After(5 * time.Second)
	}
}

func (au *AuctionUseCase) CreateAuction(
	ctx context.Context,
	auctionInput AuctionInputDTO) *internal_error.InternalError {
	auction, err := auction_entity.CreateAuction(
		auctionInput.ProductName,
		auctionInput.Category,
		auctionInput.Description,
		auction_entity.ProductCondition(auctionInput.Condition))
	if err != nil {
		return err
	}

	if err := au.auctionRepositoryInterface.CreateAuction(
		ctx, auction); err != nil {
		return err
	}

	return nil
}
