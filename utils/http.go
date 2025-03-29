package utils

type HttpResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func FormatSuccessResponse(data interface{}) HttpResponse {
	return HttpResponse{
		Success: true,
		Data:    data,
	}
}
