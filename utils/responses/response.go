package responses

// Meta
type TResponseMeta struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type TResponseMetaPage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Total   int64  `json:"total"`
}

// Success
type TSuccessResponse struct {
	Meta    TResponseMeta `json:"meta"`
	Results interface{}   `json:"results"`
}

func SuccessResponse(message string, data interface{}) interface{} {
	if data == nil {
		return TErrorResponse{
			Meta: TResponseMeta{
				Success: true,
				Message: message,
			},
		}
	} else {
		return TSuccessResponse{
			Meta: TResponseMeta{
				Success: true,
				Message: message,
			},
			Results: data,
		}
	}
}


// Error
type TErrorResponse struct {
	Meta TResponseMeta `json:"meta"`
}

func ErrorResponse(message string) interface{} {
	return TErrorResponse{
		Meta: TResponseMeta{
			Success: false,
			Message: message,
		},
	}
}

// Pagination
type TSuccessResponsePage struct {
	Meta    TResponseMetaPage `json:"meta"`
	Results interface{}       `json:"results"`
}

func SuccessResponsePage(message string, page int, limit int, total int64, data interface{}) interface{} {
	return TSuccessResponsePage{
		Meta: TResponseMetaPage{
			Success: true,
			Message: message,
			Page:    page,
			Limit:   limit,
			Total:   total,
		},
		Results: data,
	}
}