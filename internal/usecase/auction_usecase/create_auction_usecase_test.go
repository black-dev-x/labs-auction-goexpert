package auction_usecase

import (
	"auction-go/internal/entity/auction_entity"
	"auction-go/internal/internal_error"
	"context"
	"testing"
	"time"
)

type AuctionRepositoryMock struct {
	CreateAuctionFn          func(ctx context.Context, auction *auction_entity.Auction) *internal_error.InternalError
	FindAuctionByIdFn        func(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError)
	FindAuctionsFn           func(ctx context.Context, status auction_entity.AuctionStatus, category, name string) ([]auction_entity.Auction, *internal_error.InternalError)
	GetNextExpiredAuctionsFn func() ([]auction_entity.Auction, *internal_error.InternalError)
	UpdateAuctionStatusFn    func(id string, status auction_entity.AuctionStatus) *internal_error.InternalError
}

func (m *AuctionRepositoryMock) CreateAuction(ctx context.Context, auction *auction_entity.Auction) *internal_error.InternalError {
	m.CreateAuctionFn(ctx, auction)
}

func (m *AuctionRepositoryMock) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	return m.FindAuctionByIdFn(ctx, id)
}

func (m *AuctionRepositoryMock) FindAuctions(ctx context.Context, status auction_entity.AuctionStatus, category, name string) ([]auction_entity.Auction, *internal_error.InternalError) {
	return m.FindAuctionsFn(ctx, status, category, name)
}

func (m *AuctionRepositoryMock) GetNextExpiredAuctions() ([]auction_entity.Auction, *internal_error.InternalError) {
	return m.GetNextExpiredAuctionsFn()
}

func (m *AuctionRepositoryMock) UpdateAuctionStatus(id string, status auction_entity.AuctionStatus) *internal_error.InternalError {
	return m.UpdateAuctionStatusFn(id, status)
}

func NewAuctionRepositoryMock() *AuctionRepositoryMock {
	return &AuctionRepositoryMock{}
}

func TestCheckForExpiredAuctions(t *testing.T) {
	expired25Hours := time.Now().Add(-25 * time.Hour)
	expiredAuctions := []auction_entity.Auction{
		{Id: "auction1", Status: auction_entity.Active, Timestamp: expired25Hours},
		{Id: "auction2", Status: auction_entity.Active, Timestamp: expired25Hours},
	}
	auctionRepo := NewAuctionRepositoryMock()
	auctionUseCase := NewAuctionUseCase(auctionRepo, nil)
}
