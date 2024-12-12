package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SmirnovND/gofermart/internal/domain"
	http2 "github.com/SmirnovND/gofermart/internal/pkg/http"
	"io"
	"net/http"
)

const Processed = "PROCESSED"

type ProcessingService struct {
	apiClient http2.APIClient
}

func NewProcessingService(accrualSystemAddress string, APIClient http2.APIClient) *ProcessingService {
	APIClient.SetBaseURL(accrualSystemAddress)
	return &ProcessingService{
		apiClient: APIClient,
	}
}

func (p *ProcessingService) GetOrder(number string) (*domain.OrderProcessing, error) {
	url := fmt.Sprintf("%s/api/orders/%s", p.apiClient.BaseURL, number)
	resp, err := p.apiClient.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var orderResp *domain.OrderProcessing
		if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return orderResp, nil

	case http.StatusNoContent:
		return nil, errors.New("order not registered in the system")

	case http.StatusTooManyRequests:
		retryAfter, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("too many requests: retry after %s", string(retryAfter))

	case http.StatusInternalServerError:
		return nil, errors.New("internal server error")

	default:
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}
}
