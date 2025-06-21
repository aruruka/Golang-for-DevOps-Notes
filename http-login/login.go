package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func doLoginRequest(client http.Client, requestURL, password string) (string, error) {
	loginRequest := LoginRequest{
		Password: password,
	}

	body, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("marshal error: %s", err)
	}

	response, err := client.Post(requestURL, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return "", fmt.Errorf("http Post error: %s", err)
	}

	defer response.Body.Close()

	resBody, err := io.ReadAll(response.Body)

	if err != nil {
		return "", fmt.Errorf("read response body error: %s", err)
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("invalid output (http code %s), response body: %s", response.Status, string(resBody))
	}

	if !json.Valid(resBody) {
		return "", RequestError{
			HTTPCode: response.Status,
			Body:     string(resBody),
			Err:      "no valid json returned",
		}
	}

	var loginResponse LoginResponse

	err = json.Unmarshal(resBody, &loginResponse)
	if err != nil {
		return "", RequestError{
			HTTPCode: response.Status,
			Body:     string(resBody),
			Err:      fmt.Sprintf("Page unmarshal error: %s", err),
		}
	}

	if loginResponse.Token == "" {
		return "", RequestError{
			HTTPCode: response.Status,
			Body:     string(resBody),
			Err:      "no token found",
		}
	}

	return loginResponse.Token, nil
}
