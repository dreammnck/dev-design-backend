package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GatewayPaymentReq struct {
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	CardNumber   string  `json:"card_number"`
	ExpiryMonth  int     `json:"expiry_month"`
	ExpiryYear   int     `json:"expiry_year"`
	CVV          string  `json:"cvv"`
	MerchantID   string  `json:"merchant_id"`
	OrderID      string  `json:"order_id"`
}

type GatewayPaymentRes struct {
	Status        string  `json:"status"`
	Message       string  `json:"message"`
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	OrderID       string  `json:"order_id"`
}

type PaymentGatewayAdapter interface {
	ProcessPayment(req GatewayPaymentReq) (*GatewayPaymentRes, error)
}

type paymentGatewayAdapter struct {
	baseURL string
}

func NewPaymentGatewayAdapter(baseURL string) PaymentGatewayAdapter {
	return &paymentGatewayAdapter{baseURL: baseURL}
}

func (a *paymentGatewayAdapter) ProcessPayment(req GatewayPaymentReq) (*GatewayPaymentRes, error) {
	url := fmt.Sprintf("%s/v1/payments", a.baseURL)
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("payment gateway returned status: %d", resp.StatusCode)
	}

	var res GatewayPaymentRes
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
