package utils

type Response struct {
	Meta meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ResponseApi(code int, message string, data interface{}) Response {
	meta := meta{
		Code:    code,
		Message: message,
	}

	return Response{
		Meta: meta,
		Data: data,
	}
}
