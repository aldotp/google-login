package util

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func ApiResponse(code int, message string, status string, data interface{}) Response {
	return Response{
		Meta: Meta{
			Code:    code,
			Message: message,
			Status:  status,
		},
		Data: data,
	}
}
