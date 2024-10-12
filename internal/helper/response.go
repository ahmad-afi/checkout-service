package helper

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	Success = "Success"
)

var EmptyMap = fiber.Map{}

func BuildResponse(ctx *fiber.Ctx, status bool, message string, data interface{}, code int) error {
	// var statusStr string
	// if status {
	// 	statusStr = Success
	// message = Success
	// } else {
	// statusStr =
	// switch code {
	// case 400:
	// 	statusStr = BadRequest
	// case 404:
	// 	statusStr = NotFound
	// default:
	// 	statusStr = "Something went wrong"
	// }
	// }

	return ctx.Status(code).JSON(&Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// for usecase
func HelperErrorResponse(err error, message ...string) *ErrorStruct {
	if message == nil {
		message = append(message, "Something went wrong")
	}
	if strings.Contains(err.Error(), "Not Found") {
		return &ErrorStruct{
			Code:    http.StatusNotFound,
			Err:     err,
			Message: message[0],
		}
	}
	return &ErrorStruct{
		Code:    http.StatusBadRequest,
		Err:     err,
		Message: message[0],
	}
}
