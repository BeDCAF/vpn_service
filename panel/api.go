package panel

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type VPNService interface {
	AddUser(ctx context.Context, params Response) ([]string, error)
	DeleteUser(ctx context.Context, params Response) error
}

type ClientVPN struct {
	httpClient *http.Client
	url        string
	token      string
}

func NewClientVPN(url, token string) (*ClientVPN, error) {
	if url == "" || token == "" {
		return nil, errors.New("URL, method, and token must not be empty")
	}

	return &ClientVPN{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		url:   url,
		token: token,
	}, nil
}

func (r *ClientVPN) AddUser(ctx context.Context, params Response) ([]string, error) {
	data, err := params.ResponseToByte()
	if err != nil {
		return nil, fmt.Errorf("marshal params: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.url+"/panel/api/clients/add", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("create add request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do add request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("add user failed: status %d, body: %s", resp.StatusCode, body)
	}
	io.Copy(io.Discard, resp.Body)

	req, err = http.NewRequestWithContext(ctx, "GET", r.url+"/panel/api/clients/links/"+params.Client.Email, nil)
	if err != nil {
		return nil, fmt.Errorf("create get links request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do get links request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get links failed: status %d, body: %s", resp.StatusCode, body)
	}

	var apiResp AddUserResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response JSON: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: success=false, body: %s", body)
	}

	return apiResp.Obj, nil
}

func (r *ClientVPN) DeleteUser(ctx context.Context, params Response) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.url+"/panel/api/clients/del/"+params.Client.Email, nil)
	if err != nil {
		return fmt.Errorf("create delete request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do delete request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete user failed: status %d, body: %s", resp.StatusCode, body)
	}

	var apiResp AddUserResponse
	if err := json.Unmarshal(body, &apiResp); err == nil {
		if !apiResp.Success {
			return fmt.Errorf("API error: success=false, body: %s", body)
		}
	}

	return nil
}
