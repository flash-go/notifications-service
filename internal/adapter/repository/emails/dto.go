package adapter

import "encoding/json"

// Success codes

const (
	sendEmailSuccessCode     = 200
	sendEmailBadRequestCode  = 400
	sendEmailUnautorizedCode = 401
)

// Responses

type sendSuccessResponse struct {
	Result    bool   `json:"result"`
	Messageid string `json:"messageid"`
}

type sendErrorResponse struct {
	Result bool            `json:"result"`
	Errors json.RawMessage `json:"errors"`
}
