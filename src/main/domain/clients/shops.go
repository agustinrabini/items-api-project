package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/agustinrabini/items-api-project/src/main/api/config"
	"github.com/agustinrabini/items-api-project/src/main/domain/models"
	"github.com/jopitnow/go-jopit-toolkit/goauth"
	"github.com/jopitnow/go-jopit-toolkit/rest"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

const (
	ShopsBaseEndpoint = "/shops"
)

type ShopClient interface {
	GetShopByUserID(ctx context.Context) (models.Shop, apierrors.ApiError)
}

type shopClient struct {
	Builder *rest.RequestBuilder
}

func NewShopClient() ShopClient {
	builder := &rest.RequestBuilder{
		BaseURL:        config.InternalBaseShopsClient,
		Timeout:        5 * time.Second,
		ContentType:    rest.JSON,
		EnableCache:    false,
		DisableTimeout: false,
		CustomPool:     &rest.CustomPool{MaxIdleConnsPerHost: 100},
		FollowRedirect: true,
		MetricsConfig:  rest.MetricsReportConfig{TargetId: "shops-api"},
	}

	return &shopClient{Builder: builder}
}

func (client shopClient) GetShopByUserID(ctx context.Context) (models.Shop, apierrors.ApiError) {
	var headers = http.Header{}
	var shop models.Shop

	headers.Add("Authorization", fmt.Sprint(ctx.Value(goauth.FirebaseAuthHeader)))
	headers.Add("X-Trace-ID", fmt.Sprint(ctx.Value("X-Trace-ID")))

	response := client.Builder.Get(ShopsBaseEndpoint, rest.Headers(headers))

	if response.Response == nil {
		return models.Shop{}, apierrors.NewInternalServerApiError(fmt.Sprint("unexpected error getting shop, url: "+ShopsBaseEndpoint), response.Err)
	}

	if response.StatusCode == http.StatusNotFound {
		return models.Shop{}, apierrors.NewNotFoundApiError("shop not found")
	}

	if response.StatusCode != http.StatusOK {
		return models.Shop{}, apierrors.NewInternalServerApiError(fmt.Sprintf("error getting shop with state %d, url: "+ShopsBaseEndpoint, response.StatusCode), response.Err)
	}

	if err := json.Unmarshal(response.Bytes(), &shop); err != nil {
		return models.Shop{}, apierrors.NewInternalServerApiError("unexpected error unmarshalling shop json response. value: "+string(response.Bytes()), err)
	}

	return shop, nil
}
