package auction_usecase

import (
	"auction-go/internal/entity/auction_entity"
	"auction-go/internal/internal_error"
	"context"
	"fmt"
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
	return m.CreateAuctionFn(ctx, auction)
}

func (m *AuctionRepositoryMock) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	return m.FindAuctionByIdFn(ctx, id)
}

func (m *AuctionRepositoryMock) FindAuctions(ctx context.Context, status auction_entity.AuctionStatus, category, name string) ([]auction_entity.Auction, *internal_error.InternalError) {
	return m.FindAuctionsFn(ctx, status, category, name)
}

func (m *AuctionRepositoryMock) GetNextExpiredAuctions() ([]auction_entity.Auction, *internal_error.InternalError) {
	fmt.Println("Mock: GetNextExpiredAuctions called")
	return m.GetNextExpiredAuctionsFn()
}

func (m *AuctionRepositoryMock) UpdateAuctionStatus(id string, status auction_entity.AuctionStatus) *internal_error.InternalError {
	return m.UpdateAuctionStatusFn(id, status)
}

func NewAuctionRepositoryMock() *AuctionRepositoryMock {
	return &AuctionRepositoryMock{}
}

func TestCheckForExpiredAuctions_BackgroundJob(t *testing.T) {
	called := make(chan string, 2)
	auctionRepo := NewAuctionRepositoryMock()

	auctionRepo.GetNextExpiredAuctionsFn = func() ([]auction_entity.Auction, *internal_error.InternalError) {
		expired25Hours := time.Now().Add(-25 * time.Hour)
		return []auction_entity.Auction{
			{Id: "auction1", Status: auction_entity.Active, Timestamp: expired25Hours},
			{Id: "auction2", Status: auction_entity.Active, Timestamp: expired25Hours},
		}, nil
	}

	auctionRepo.UpdateAuctionStatusFn = func(id string, status auction_entity.AuctionStatus) *internal_error.InternalError {
		called <- id
		return nil
	}

	auction := NewAuctionUseCase(auctionRepo, nil)
	if auction == nil {
		t.Fatal("Failed to create AuctionUseCase")
	}

	ids := map[string]bool{}
	timeout := time.After(10 * time.Second)

	for i := 0; i < 2; i++ {
		select {
		case id := <-called:
			ids[id] = true
		case <-timeout:
			t.Fatal("timeout waiting for UpdateAuctionStatusFn to be called")
		}
	}

}
