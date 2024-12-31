package clients

import (
	"context"

	"github.com/agustinrabini/items-api-project/src/main/domain/models"
	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

type ShopClientMock struct {
	HandleGetShopByUserID func(ctx context.Context) (models.Shop, apierrors.ApiError)
}

func NewShopClientMock() ShopClientMock {
	return ShopClientMock{}
}

func (mock ShopClientMock) GetShopByUserID(ctx context.Context) (models.Shop, apierrors.ApiError) {
	if mock.HandleGetShopByUserID != nil {
		return mock.HandleGetShopByUserID(ctx)
	}
	return models.Shop{}, nil
}
