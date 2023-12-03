package models

import "gorm.io/gorm"

// User is the user model
type User struct {
	gorm.Model
	Name         string         `json:"name"`
	CompanyName  string         `json:"company_name"`
	Password     string         `json:"password"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone"`
	TokenBalance int            `json:"token_balance"`
	ApiResponses []ApiResponses `gorm:"foreignKey:UserID"`
	ServiceLogs  []ServiceLogs  `gorm:"foreignKey:UserID"`
}

type Services struct {
	gorm.Model
	Name        string        `json:"service_name"`
	Price       int           `json:"price"`
	Description string        `json:"description"`
	ServiceURL  string        `json:"service_url"`
	ServiceLogs []ServiceLogs `gorm:"foreignKey:ServiceID"`
}

type ServiceLogs struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	Log       string `json:"log"`
	ServiceID uint   `json:"service_id"`
}

type ApiResponses struct {
	gorm.Model
	RequestHeaders    string `json:"request_headers"`
	RequestBody       string `json:"request_body"`
	RequestTime       string `json:"request_time"`
	RequestStatusCode string `json:"request_status_code"`
	RequestUrl        string `json:"request_url"`
	RequestMethod     string `json:"request_method"`
	RequestResponse   string `json:"request_response"`
	RequestPrice      int    `json:"request_price"`
	UserID            uint   `json:"user_id"`
}
