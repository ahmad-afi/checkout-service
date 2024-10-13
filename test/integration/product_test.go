package integration_test

import (
	"checkout-service/internal/helper"
	"checkout-service/internal/usecase/productu"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetProduct(t *testing.T) {
	request := httptest.NewRequest("GET", "/v1/product", nil)
	response, err := app.Test(request, int(600*time.Second))
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	respBytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	fmt.Println("respBytes", string(respBytes))

	var resp helper.Response
	err = json.Unmarshal(respBytes, &resp)
	assert.Nil(t, err)

	// Cast resp.Data to []interface{}
	dataInterface, ok := resp.Data.([]interface{})
	assert.True(t, ok)

	var products []productu.GetListProduct
	for _, item := range dataInterface {
		// Convert each item to map and then marshal it to json
		itemBytes, err := json.Marshal(item)
		assert.Nil(t, err)

		// Unmarshal it into the target type
		var product productu.GetListProduct
		err = json.Unmarshal(itemBytes, &product)
		assert.Nil(t, err)

		// Append to the products slice
		products = append(products, product)
	}

	// Now products is of type []productu.GetListProduct
	fmt.Println("products", products)
}
