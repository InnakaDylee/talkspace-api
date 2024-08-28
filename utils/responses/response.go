package responses

// Success
type TResponseMeta struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

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
type TResponseMetaPage struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	CurrentPage     int    `json:"current_page"`
	TotalPages      int    `json:"total_pages"`
	TotalItems      int    `json:"total_items"`
	ItemsPerPage    int    `json:"items_per_page"`
	HasNextPage     bool   `json:"has_next_page"`
	HasPreviousPage bool   `json:"has_previous_page"`
}


type TSuccessResponsePage struct {
	Meta    TResponseMetaPage `json:"meta"`
	Results interface{}       `json:"results"`
}

func SuccessResponsePage(message string, page int, limit int, totaldata int64, data interface{}) TSuccessResponsePage {
	totalPages := int(totaldata / int64(limit))
	if totaldata % int64(limit) != 0 {
		totalPages++
	}

	return TSuccessResponsePage{
		Meta: TResponseMetaPage{
			Success:         true,
			Message:         message,
			CurrentPage:     page,
			TotalPages:      totalPages,
			TotalItems:      int(totaldata),
			ItemsPerPage:    limit,
			HasNextPage:     page < totalPages,
			HasPreviousPage: page > 1,
		},
		Results: data,
	}
}