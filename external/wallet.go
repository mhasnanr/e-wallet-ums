package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ExternalWallet struct{}

type WalletRequest struct {
	UserID int `json:"user_id"`
}

func (w *ExternalWallet) CreateWallet(userID int) error {
	var err error

	postBody, _ := json.Marshal(WalletRequest{UserID: userID})
	requestBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://localhost:8081/wallets/v1", "application/json", requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("wallet service error: status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
