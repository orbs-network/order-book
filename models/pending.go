package models

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Pending struct {
	SwapId uuid.UUID `json:"swapId"`
	TxHash string    `json:"txHash"`
}

func (p *Pending) PendingToMap() map[string]interface{} {
	return map[string]interface{}{
		"swapId": p.SwapId.String(),
		"txHash": p.TxHash,
	}
}

func (p *Pending) MapToPending(data map[string]string) error {
	if len(data) == 0 {
		return fmt.Errorf("no data provided")
	}

	swapIdStr, exists := data["swapId"]
	if !exists {
		return fmt.Errorf("no swapId provided")
	}

	txHash, exists := data["txHash"]
	if !exists {
		return fmt.Errorf("no txHash provided")
	}

	swapId, err := uuid.Parse(swapIdStr)
	if err != nil {
		return fmt.Errorf("failed to parse swapId: %s", err)
	}

	p.SwapId = swapId
	p.TxHash = txHash

	return nil
}

func (p *Pending) ToJson() ([]byte, error) {
	jsonData, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pending: %s", err)
	}

	return jsonData, nil
}

func (p *Pending) FromJson(data []byte) error {
	err := json.Unmarshal(data, p)
	if err != nil {
		return fmt.Errorf("failed to unmarshal pending: %s", err)
	}

	return nil
}
