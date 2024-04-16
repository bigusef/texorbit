package city

type cityInput struct {
	NameEn   string `json:"name_en" validate:"required"`
	NameAr   string `json:"name_ar" validate:"required"`
	IsActive *bool  `json:"is_active" validate:"required"`
}

type cityResponse struct {
	ID       int64  `json:"id"`
	NameEn   string `json:"name_en"`
	NameAr   string `json:"name_ar"`
	IsActive bool   `json:"is_active"`
}

type activeCityResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
