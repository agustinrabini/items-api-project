package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/agustinrabini/items-api-project/src/main/api/config"
	"github.com/agustinrabini/items-api-project/src/main/domain/models"
	"github.com/agustinrabini/items-api-project/src/main/domain/models/dto"
	"github.com/jopitnow/go-jopit-toolkit/goauth"
	"github.com/jopitnow/go-jopit-toolkit/rest"
	"github.com/jopitnow/go-jopit-toolkit/tracing"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

const (
	PricesBaseEndpoint = "/prices"
	PricesItemsPrices  = "/items"
)

var errorPricingService error = fmt.Errorf("error at Prices Services. URL: ")
var errorPricingServiceUnmarshal error = fmt.Errorf("Unexpected error unmarshalling price json response. Body: ")

type PriceClient interface {
	GetPriceByItemID(ctx context.Context, itemID string) (models.Price, apierrors.ApiError)
	GetItemsPrices(ctx context.Context, itemsIDs []string) (models.Prices, apierrors.ApiError)
	CreatePrice(ctx context.Context, price *models.Price) apierrors.ApiError
	UpdatePrice(ctx context.Context, price *models.Price) apierrors.ApiError
	DeletePrice(ctx context.Context, priceID string) apierrors.ApiError
}

type priceClient struct {
	Builder *rest.RequestBuilder
}

func NewPriceClient() PriceClient {
	builder := &rest.RequestBuilder{
		BaseURL:        config.InternalBasePricesClient,
		Timeout:        5 * time.Second,
		ContentType:    rest.JSON,
		EnableCache:    false,
		DisableTimeout: false,
		CustomPool:     &rest.CustomPool{MaxIdleConnsPerHost: 100},
		FollowRedirect: true,
		MetricsConfig:  rest.MetricsReportConfig{TargetId: "prices-api"},
	}

	return &priceClient{Builder: builder}
}

func (client priceClient) GetPriceByItemID(ctx context.Context, itemID string) (models.Price, apierrors.ApiError) {
	var headers = http.Header{}
	var price models.Price

	headers.Add("X-Trace-Id", fmt.Sprint(ctx.Value(tracing.XtraceHeaderKey)))

	endpoint := fmt.Sprintf("%s/item/%s", PricesBaseEndpoint, itemID)
	response := client.Builder.Get(endpoint, rest.Context(ctx), rest.Headers(headers))

	if response.Response == nil {
		return models.Price{}, apierrors.NewInternalServerApiError(fmt.Sprintf("unexpected error getting price, url: %s", endpoint), response.Err)
	}

	if response.StatusCode == http.StatusNotFound {
		return models.Price{}, apierrors.NewNotFoundApiError("price not found")
	}

	if response.Response == nil || response.StatusCode != http.StatusOK {
		return models.Price{}, apierrors.NewInternalServerApiError(fmt.Sprintf("error getting price with state %d, url: %s", response.StatusCode, endpoint), response.Err)
	}

	if err := json.Unmarshal(response.Bytes(), &price); err != nil {
		return models.Price{}, apierrors.NewInternalServerApiError("unexpected error unmarshalling price json response. value: "+string(response.Bytes()), err)
	}

	return price, nil
}

func (client priceClient) GetItemsPrices(ctx context.Context, ids []string) (models.Prices, apierrors.ApiError) {

	var rawBody []byte
	prices := models.Prices{}

	/* 	authHeader := http.Header{}
	   	authHeader.Add("Authorization", fmt.Sprint(ctx.Value(goauth.FirebaseAuthHeader)))
	*/
	traceHeader := http.Header{}
	xid := ctx.Value(tracing.XtraceHeaderKey)
	traceHeader.Add("X-Trace-ID", fmt.Sprint(xid))

	req := dto.RequestItemsList{
		ItemsIDs: ids,
	}

	endpoint := fmt.Sprintf("%s%s", PricesBaseEndpoint, PricesItemsPrices)

	response := client.Builder.Post(endpoint, req, rest.Context(ctx), rest.Headers(traceHeader)) //, rest.Headers(authHeader)

	if response.Response == nil || (response.StatusCode != http.StatusNotFound && response.StatusCode != http.StatusOK) {
		return models.Prices{}, apierrors.NewInternalServerApiError(fmt.Sprint(errorPricingService, endpoint), response.Err)
	}

	if response.StatusCode == http.StatusNotFound {
		return models.Prices{}, apierrors.NewNotFoundApiError("no prtice found")
	}

	rawBody = response.Bytes()
	if rawBody != nil {

		if unmarshallError := json.Unmarshal(rawBody, &prices); unmarshallError != nil {
			return models.Prices{}, apierrors.NewInternalServerApiError(errorPricingServiceUnmarshal.Error()+string(rawBody), unmarshallError)
		}

		return prices, nil
	}

	return models.Prices{}, apierrors.NewNotFoundApiError("price not found")
}

func (client priceClient) CreatePrice(ctx context.Context, price *models.Price) apierrors.ApiError {
	var headers = http.Header{}
	var response *rest.Response

	headers.Add("X-Trace-Id", fmt.Sprint(ctx.Value(tracing.XtraceHeaderKey)))

	response = client.Builder.Post(PricesBaseEndpoint, price, rest.Context(ctx), rest.Headers(headers))

	if response.Response == nil {
		return apierrors.NewInternalServerApiError(fmt.Sprintf("unexpected error creating price, url: %s", PricesBaseEndpoint), response.Err)
	}

	if response.StatusCode != http.StatusCreated {
		return apierrors.NewInternalServerApiError(fmt.Sprintf("error creating price with state %d, url: %s", response.StatusCode, PricesBaseEndpoint), response.Err)
	}

	return nil
}

func (client priceClient) UpdatePrice(ctx context.Context, price *models.Price) apierrors.ApiError {
	var headers = http.Header{}
	var response *rest.Response

	headers.Add("X-Trace-Id", fmt.Sprint(ctx.Value(tracing.XtraceHeaderKey)))
	headers.Add("Authorization", fmt.Sprint(ctx.Value(goauth.FirebaseAuthHeader)))

	endpoint := fmt.Sprintf("%s/%s", PricesBaseEndpoint, price.ID)
	response = client.Builder.Put(endpoint, price, rest.Context(ctx), rest.Headers(headers))

	if response.Response == nil {
		return apierrors.NewInternalServerApiError(fmt.Sprintf("unexpected error updating price, url: %s", endpoint), response.Err)
	}

	if response.StatusCode == http.StatusNotFound {
		return apierrors.NewNotFoundApiError("price not found")
	}

	if response.StatusCode != http.StatusOK {
		return apierrors.NewInternalServerApiError(fmt.Sprintf("error updating price with state %d, url: %s", response.StatusCode, endpoint), response.Err)
	}

	return nil
}

func (client priceClient) DeletePrice(ctx context.Context, itemID string) apierrors.ApiError {

	var response *rest.Response

	headers := http.Header{}
	xid := ctx.Value(tracing.XtraceHeaderKey)
	headers.Add("X-Trace-ID", fmt.Sprint(xid))
	headers.Add("Authorization", fmt.Sprint(ctx.Value(goauth.FirebaseAuthHeader)))

	endpoint := fmt.Sprintf("%s/item/%s", PricesBaseEndpoint, itemID)
	response = client.Builder.Delete(endpoint, rest.Context(ctx), rest.Headers(headers))

	if response.Response == nil {
		return apierrors.NewInternalServerApiError(fmt.Sprintf("unexpected error updating price, url: %s", endpoint), response.Err)
	}

	if response.StatusCode == http.StatusNotFound {
		return apierrors.NewNotFoundApiError("price not found")
	}

	if response.StatusCode != http.StatusNoContent {
		return apierrors.NewApiError(fmt.Sprintf("error updating price with state %d, url: %s", response.StatusCode, endpoint), "error hitting price api", response.StatusCode, apierrors.CauseList{})
	}

	return nil
}

func getQueryParams(itemsIDs []string) string {
	var query = "?"
	var param = "id"

	for i := 0; i < len(itemsIDs); i++ {
		if i > 0 {
			query += "&"
		}

		query += param + "=" + itemsIDs[i]
	}

	return query
}
