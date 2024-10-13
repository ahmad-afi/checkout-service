package integration_test

import (
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

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	fmt.Println("bytes", string(bytes))
	// assert.Equal(t, "Hello World", string(bytes))
}
