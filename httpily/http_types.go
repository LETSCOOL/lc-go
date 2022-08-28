// Package httpily
// 提供 http 輸入輸出的協助工具
package httpily

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// HttpErrorResponse 用來表示 http 回應錯誤的資料結構。
type HttpErrorResponse struct {
	// 狀態碼
	Code int `json:"code" example:"400"`

	// 狀態訊息
	Message string `json:"message" example:"bad request"`
}

// UndefinedStructure presents a request or response structure which are not defined yet.
// 當尚未定義 swagger 的 request/response 的型態時，可以暫時使用這個資料結構。
type UndefinedStructure struct {
	// 狀態訊息
	Message string `json:"message" example:"Undefined structure"`
}

// HttpGeneralResponse presents a response structure which are not defined yet.
// 當尚未定義 swagger 的 request/response 的型態時，可以暫時使用這個資料結構。
type HttpGeneralResponse struct {
	// 狀態碼
	Code int `json:"code" example:"400"`

	// 狀態訊息
	Message string `json:"message" example:"Undefined structure"`
}

// ExportJsonResponse 輸出標準的json response
// 如果 message 為 error 類型或者 status 為 4xx 或者 5xx，輸出為 HttpErrorResponse 結構。
// 否則一律輸出為 HttpGeneralResponse 結構。
// 	TODO: 也許有一天 ctx 要使用 generic type 來支援各種 web server
func ExportJsonResponse(ctx *gin.Context, status int, message interface{}) {
	switch v := message.(type) {
	case error:
		resp := HttpErrorResponse{
			Code:    status,
			Message: v.Error(),
		}
		ctx.JSON(status, resp)
		break
	default:
		switch status % 100 {
		case 4:
		case 5:
			// 目前 HttpErrorResponse 與 HttpGeneralResponse 相同，
			// 有必要特別處理嗎？
			resp := HttpErrorResponse{
				Code:    status,
				Message: fmt.Sprint(message),
			}
			ctx.JSON(status, resp)
			break
		default:
			resp := HttpGeneralResponse{
				Code:    status,
				Message: fmt.Sprint(message),
			}
			ctx.IndentedJSON(status, resp)
			break
		}
		break
	}
}

// ExportJsonError creates an error response
// err 若為 error 型態，則會藉由 error.Error() 當作回傳訊息，而其它的型別則都透過 fmt.Sprint() 輸出為字串。
//   TODO: 也許有一天 ctx 要使用 generic type 來支援各種 web server
func ExportJsonError(ctx *gin.Context, status int, err interface{}) {
	var message string
	switch v := err.(type) {
	case error:
		message = v.Error()
		break
	default:
		message = fmt.Sprint(err)
		break
	}

	er := HttpErrorResponse{
		Code:    status,
		Message: message,
	}

	ctx.JSON(status, er)
}
