package prices

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/agustinrabini/items-api-project/src/main/domain/clients"
	"github.com/agustinrabini/items-api-project/src/tests/internal/domain/mocks"
	"github.com/jarcoal/httpmock"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetPriceByItemID_Success(t *testing.T) {
	var model = mocks.Price
	var endpoint = fmt.Sprintf("%s/item/%s", clients.PricesBaseEndpoint, mocks.ItemIdOne)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", endpoint,
		func(req *http.Request) (*http.Response, error) {
			if !httpmock.HeaderExists("X-Trace-Id").Check(req) {
				return nil, errors.New("missing X-Trace-ID header")
			}

			response, err := httpmock.NewJsonResponse(200, model)

			return response, err
		},
	)

	api := clients.NewPriceClient()

	body, err := api.GetPriceByItemID(context.Background(), mocks.ItemIdOne)

	assert.Nil(t, err)
	assert.Equal(t, model, body)

}

func TestClient_GetPriceByItemID_Response_Error(t *testing.T) {
	var endpoint = fmt.Sprintf("%s/item/%s", clients.PricesBaseEndpoint, mocks.ItemIdOne)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("response error")
		},
	)

	api := clients.NewPriceClient()

	_, err := api.GetPriceByItemID(context.Background(), mocks.ItemIdOne)

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, fmt.Sprintf("unexpected error getting price, url: %s", endpoint), err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestClient_GetPriceByItemID_Not_Found_Error(t *testing.T) {
	var endpoint = fmt.Sprintf("%s/item/%s", clients.PricesBaseEndpoint, mocks.ItemIdOne)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(404, nil)
		},
	)

	api := clients.NewPriceClient()

	_, err := api.GetPriceByItemID(context.Background(), mocks.ItemIdOne)

	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, "price not found", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestClient_GetPriceByItemID_Internal_Server_Error(t *testing.T) {
	var endpoint = fmt.Sprintf("%s/item/%s", clients.PricesBaseEndpoint, mocks.ItemIdOne)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(505, nil)
		},
	)

	api := clients.NewPriceClient()

	_, err := api.GetPriceByItemID(context.Background(), mocks.ItemIdOne)

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, fmt.Sprintf("error getting price with state 505, url: %s", endpoint), err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestClient_GetPriceByItemID_UnMarshall_Error(t *testing.T) {
	var endpoint = fmt.Sprintf("%s/item/%s", clients.PricesBaseEndpoint, mocks.ItemIdOne)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, "fake")
		},
	)

	api := clients.NewPriceClient()

	_, err := api.GetPriceByItemID(context.Background(), mocks.ItemIdOne)

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, "unexpected error unmarshalling price json response. value: \"fake\"", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestClient_GetItemsPrices_Success(t *testing.T) {
	var model = mocks.Prices
	var endpoint = fmt.Sprintf("%s%s", clients.PricesBaseEndpoint, clients.PricesItemsPrices)
	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", endpoint,
		func(req *http.Request) (*http.Response, error) {
			if !httpmock.HeaderExists("X-Trace-Id").Check(req) {
				return nil, errors.New("missing X-Trace-ID header")
			}

			response, err := httpmock.NewJsonResponse(200, model)

			return response, err
		},
	)

	api := clients.NewPriceClient()

	body, err := api.GetItemsPrices(context.Background(), mocks.ItemIds.Items)

	assert.Nil(t, err)
	assert.Equal(t, model, body)
}

func TestClient_GetItemsPrices_Response_Error(t *testing.T) {
	var endpoint = fmt.Sprintf("%s%s", clients.PricesBaseEndpoint, clients.PricesItemsPrices)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("response error")
		},
	)

	api := clients.NewPriceClient()

	_, err := api.GetItemsPrices(context.Background(), mocks.ItemIds.Items)

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestClient_GetItemsPrices_Not_Found_Error(t *testing.T) {
	var endpoint = fmt.Sprintf("%s%s", clients.PricesBaseEndpoint, clients.PricesItemsPrices)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(404, nil)
		},
	)

	api := clients.NewPriceClient()

	_, err := api.GetItemsPrices(context.Background(), mocks.ItemIds.Items)

	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestClient_GetItemsPrices_Internal_Server_Error(t *testing.T) {
	var endpoint = fmt.Sprintf("%s%s", clients.PricesBaseEndpoint, clients.PricesItemsPrices)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(505, nil)
		},
	)

	api := clients.NewPriceClient()

	_, err := api.GetItemsPrices(context.Background(), mocks.ItemIds.Items)

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestClient_GetItemsPrices_UnMarshall_Error(t *testing.T) {
	var endpoint = fmt.Sprintf("%s%s", clients.PricesBaseEndpoint, clients.PricesItemsPrices)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, "fake")
		},
	)

	api := clients.NewPriceClient()

	_, err := api.GetItemsPrices(context.Background(), mocks.ItemIds.Items)

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestClient_CreatePrice_Success(t *testing.T) {
	var model = mocks.Price

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", clients.PricesBaseEndpoint,
		func(req *http.Request) (*http.Response, error) {
			b, err := io.ReadAll(req.Body)

			body, err := json.Marshal(model)
			if err != nil {
				return nil, errors.New("error marshalling model")
			}

			if err != nil || !bytes.Contains(b, body) {
				return nil, errors.New("missing body")
			}

			if !httpmock.HeaderExists("X-Trace-Id").Check(req) {
				return nil, errors.New("missing X-Trace-ID header")
			}

			response, err := httpmock.NewJsonResponse(201, nil)

			return response, err
		},
	)

	api := clients.NewPriceClient()

	err := api.CreatePrice(context.Background(), &model)

	assert.Nil(t, err)
}

func TestClient_CreatePrice_Response_Error(t *testing.T) {
	var model = mocks.Price

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", clients.PricesBaseEndpoint,
		func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("response error")
		},
	)

	api := clients.NewPriceClient()

	err := api.CreatePrice(context.Background(), &model)

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, fmt.Sprintf("unexpected error creating price, url: %s", clients.PricesBaseEndpoint), err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestClient_CreatePrice_Internal_Server_Error(t *testing.T) {
	var model = mocks.Price

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", clients.PricesBaseEndpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(505, nil)
		},
	)

	api := clients.NewPriceClient()

	err := api.CreatePrice(context.Background(), &model)

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, fmt.Sprintf("error creating price with state 505, url: %s", clients.PricesBaseEndpoint), err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestClient_UpdatePrice_Success(t *testing.T) {
	var model = mocks.Price
	var endpoint = fmt.Sprintf("%s/%s", clients.PricesBaseEndpoint, model.ID)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", endpoint,
		func(req *http.Request) (*http.Response, error) {
			b, err := io.ReadAll(req.Body)

			body, err := json.Marshal(model)
			if err != nil {
				return nil, errors.New("error marshalling model")
			}

			if err != nil || !bytes.Contains(b, body) {
				return nil, errors.New("missing body")
			}

			if !httpmock.HeaderExists("X-Trace-Id").Check(req) {
				return nil, errors.New("missing X-Trace-ID header")
			}

			response, err := httpmock.NewJsonResponse(200, nil)

			return response, err
		},
	)

	api := clients.NewPriceClient()

	err := api.UpdatePrice(context.Background(), &model)

	assert.Nil(t, err)
}

func TestClient_UpdatePrice_Response_Error(t *testing.T) {
	var model = mocks.Price
	var endpoint = fmt.Sprintf("%s/%s", clients.PricesBaseEndpoint, model.ID)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("response error")
		},
	)

	api := clients.NewPriceClient()

	err := api.UpdatePrice(context.Background(), &model)

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, fmt.Sprintf("unexpected error updating price, url: %s", endpoint), err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestClient_UpdatePrice_Not_Found_Error(t *testing.T) {
	var model = mocks.Price
	var endpoint = fmt.Sprintf("%s/%s", clients.PricesBaseEndpoint, model.ID)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(404, nil)
		},
	)

	api := clients.NewPriceClient()

	err := api.UpdatePrice(context.Background(), &model)

	assert.EqualValues(t, "not_found", err.Code())
	assert.EqualValues(t, "price not found", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestClient_UpdatePrice_Internal_Server_Error(t *testing.T) {
	var model = mocks.Price
	var endpoint = fmt.Sprintf("%s/%s", clients.PricesBaseEndpoint, model.ID)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(505, nil)
		},
	)

	api := clients.NewPriceClient()

	err := api.UpdatePrice(context.Background(), &model)

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, fmt.Sprintf("error updating price with state 505, url: %s", endpoint), err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestClient_DeletePrice_Success(t *testing.T) {
	var model = mocks.Price
	var endpoint = fmt.Sprintf("%s/item/%s", clients.PricesBaseEndpoint, model.ID)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", endpoint,
		func(req *http.Request) (*http.Response, error) {

			if !httpmock.HeaderExists("X-Trace-Id").Check(req) {
				return nil, errors.New("missing X-Trace-ID header")
			}

			response, err := httpmock.NewJsonResponse(204, nil)

			return response, err
		},
	)

	api := clients.NewPriceClient()

	err := api.DeletePrice(context.Background(), model.ID)

	assert.Nil(t, err)
}

func TestClient_DeletePrice_Response_Error(t *testing.T) {
	var model = mocks.Price
	var endpoint = fmt.Sprintf("%s/item/%s", clients.PricesBaseEndpoint, model.ID)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("response error")
		},
	)

	api := clients.NewPriceClient()

	err := api.DeletePrice(context.Background(), model.ID)

	assert.EqualValues(t, "internal_server_error", err.Code())
	assert.EqualValues(t, fmt.Sprintf("unexpected error updating price, url: %s", endpoint), err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
}

func TestClient_DeletePrice_Not_Found_Error(t *testing.T) {
	var model = mocks.Price
	var endpoint = fmt.Sprintf("%s/item/%s", clients.PricesBaseEndpoint, model.ID)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(404, nil)
		},
	)

	api := clients.NewPriceClient()

	err := api.DeletePrice(context.Background(), model.ID)

	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestClient_DeletePrice_Bad_Request_Invalid_Hex_Error(t *testing.T) {
	var model = mocks.Price
	var endpoint = fmt.Sprintf("%s/item/%s", clients.PricesBaseEndpoint, model.ID)

	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(400, nil)
		},
	)

	api := clients.NewPriceClient()

	err := api.DeletePrice(context.Background(), model.ID)

	assert.EqualValues(t, http.StatusBadRequest, err.Status())
}
