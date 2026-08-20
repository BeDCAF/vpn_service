package panel

import (
	"encoding/json"
	"errors"
)

type Client struct {
	Email      string `json:"email"`
	TotalGB    int64  `json:"totalGB"`
	ExpiryTime int64  `json:"expiryTime"`
	LimitIp    int64  `json:"limitIp"`
	Enable     bool   `json:"enable"`
}

type Response struct {
	Client     Client  `json:"client"`
	InboundIds []int64 `json:"inboundIds"`
}

func (r Response) ResponseToByte() ([]byte, error) {
	return json.Marshal(r)
}

type AddUserResponse struct {
	Success bool     `json:"success"`
	Obj     []string `json:"obj"`
}

func (c *Client) ValidateClient() error {
	if len(c.Email) == 0 {
		return errors.New("The email can not be empty")
	}

	return nil
}
