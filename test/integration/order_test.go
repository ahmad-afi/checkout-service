package integration_test

import (
	"bytes"
	"checkout-service/internal/helper"
	"checkout-service/internal/usecase/orderu"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrderFlow(t *testing.T) {
	t.Run("CheckOrder", CheckOrder)
	t.Run("CreateOrder", CreateOrder)
	t.Run("GetOrderHistory", GetOrderHistory)
	t.Run("GetOrderDetail", GetOrderDetail)

}
func CheckOrder(t *testing.T) {
	t.Helper()
	data := orderu.OrderCreateReq{
		Data: []orderu.OrderData{
			{
				ProductID: "01HKBSM317D1K9JPBKSAT9QVY9",
				Qty:       2,
			},
		},
	}

	dataByte, err := json.Marshal(data)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/v1/order/check", bytes.NewReader(dataByte))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request, int(600*time.Second))
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	respBytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	fmt.Println("respBytes", string(respBytes))

	var resp helper.Response
	err = json.Unmarshal(respBytes, &resp)
	assert.Nil(t, err)
}
func CreateOrder(t *testing.T) {
	t.Helper()

	data := orderu.OrderCreateReq{
		Data: []orderu.OrderData{
			{
				ProductID: "01HKBSM317D1K9JPBKSAT9QVY9",
				Qty:       2,
			},
		},
	}
	dataByte, err := json.Marshal(data)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/v1/order/confirm", bytes.NewReader(dataByte))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request, int(600*time.Second))
	assert.Nil(t, err)
	assert.Equal(t, 201, response.StatusCode)

	respBytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	fmt.Println("respBytes", string(respBytes))

	var resp helper.Response
	err = json.Unmarshal(respBytes, &resp)
	assert.Nil(t, err)
}

func GetOrderHistory(t *testing.T) {
	t.Helper()
	request := httptest.NewRequest("GET", "/v1/order", nil)
	response, err := app.Test(request, int(600*time.Second))
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	respBytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	fmt.Println("respBytes", string(respBytes))
}

func GetOrderDetail(t *testing.T) {
	t.Helper()

	t.Log("Get order first")
	request := httptest.NewRequest("GET", "/v1/order", nil)
	response, err := app.Test(request, int(600*time.Second))
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	respOrderBytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	fmt.Println("respOrderBytes", string(respOrderBytes))

	var respOrder helper.Response
	err = json.Unmarshal(respOrderBytes, &respOrder)
	assert.Nil(t, err)

	// Cast respOrder.Data to []interface{}
	dataInterface, ok := respOrder.Data.([]interface{})
	assert.True(t, ok)

	var orderhistory []orderu.OrderHistory
	for _, item := range dataInterface {
		// Convert each item to map and then marshal it to json
		itemBytes, err := json.Marshal(item)
		assert.Nil(t, err)

		// Unmarshal it into the target type
		var product orderu.OrderHistory
		err = json.Unmarshal(itemBytes, &product)
		assert.Nil(t, err)

		// Append to the orderhistory slice
		orderhistory = append(orderhistory, product)
	}

	// Now orderhistory is of type []orderu.OrderHistory
	fmt.Println("orderhistory[0]", orderhistory[0])

	assert.Greater(t, len(orderhistory), 0)

	t.Log("Get Order Detail")
	url := "/v1/order/" + orderhistory[0].ID

	requestDetail := httptest.NewRequest(http.MethodGet, url, nil)
	responseDetail, err := app.Test(requestDetail, int(600*time.Second))

	assert.Nil(t, err)
	assert.Equal(t, 200, responseDetail.StatusCode)

	respBytes, err := io.ReadAll(responseDetail.Body)
	assert.Nil(t, err)
	assert.NotNil(t, respBytes)
	fmt.Println("respBytes", string(respBytes))

	var resp helper.Response
	err = json.Unmarshal(respBytes, &resp)
	assert.Nil(t, err)
}
