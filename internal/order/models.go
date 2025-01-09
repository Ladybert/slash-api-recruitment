package order

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Order struct {
	ID              string      `json:"id"`
	UserID			string		`json:"user_id" validate:"required"`
	CustomerName    string      `json:"customer_name" validate:"required"`
	CustomerPhone   string      `json:"customer_phone" validate:"omitempty,e164"`
	CustomerAddress string      `json:"customer_address" validate:"omitempty"`
	TotalAmount     string      `json:"total_order" validate:"required,gt=0"`
	Status          string      `json:"status" validate:"required,oneof=pending done"`
	OrderTime       time.Time   `json:"order_time"`
	Expired         time.Time   `json:"expired" validate:"gtfield=OrderTime"`
	OrderItems      []OrderItem `json:"order_items" validate:"required,dive"`
}

type OrderItem struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
	Price     string `json:"price" validate:"required,gt=0"`
}

type OrderDetails struct {
    ID             string       `json:"id"`
    CustomerName   string       `json:"customer_name"`
    CustomerPhone  string       `json:"customer_phone"`
    CustomerAddress string      `json:"customer_address"`
    TotalAmount    string       `json:"total_order"`
    OrderTime      string   	`json:"order_time"`
    Expired        string   	`json:"expired"`
    Status         string       `json:"status"`
    Items          []OrderItem  `json:"items"`
}

var validate = validator.New()

func fieldErrorMessage(field, tag string) string {
	switch field {
	case "CustomerName":
		return "Customer name is required"
	case "CustomerPhone":
		if tag == "e164" {
			return "Customer phone must be in E.164 format"
		}
		return "Customer phone is optional"
	case "CustomerAddress":
		return "Customer address is optional"
	case "TotalAmount":
		if tag == "required" {
			return "Total amount is required"
		} else if tag == "gt" {
			return "Total amount must be greater than 0"
		}
	case "Status":
		if tag == "required" {
			return "Status is required"
		}
		if tag == "oneof" {
			return "Status must be either 'pending' or 'done'"
        }
	case "Expired":
		return "Expired time must be greater than order time"
	case "ProductID":
		return "Product ID is required"
	case "Quantity":
		if tag == "required" {
			return "Quantity is required"
		} else if tag == "gt" {
			return "Quantity must be greater than 0"
		}
	case "Price":
		if tag == "required" {
			return "Price is required"
		} else if tag == "gt" {
			return "Price must be greater than 0"
		}
	default:
		return fmt.Sprintf("Invalid value for field %s", field)
	}

    return fmt.Sprintf("Invalid value for field %s", field)
}

func ValidateOrder(o *Order) error {
	err := validate.Struct(o)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			return fmt.Errorf(fieldErrorMessage(fieldErr.Field(), fieldErr.Tag()))
		}
	}

	for _, item := range o.OrderItems {
		err := validate.Struct(item)
		if err != nil {
			for _, fieldErr := range err.(validator.ValidationErrors) {
				return fmt.Errorf(fieldErrorMessage(fieldErr.Field(), fieldErr.Tag()))
			}
		}
	}

	return nil
}