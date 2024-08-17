package fortee

import (
	"context"
	"errors"
	"net/http"
)

const (
	apiEndpoint = "https://fortee.jp"
)

var (
	ErrLoginFailed = errors.New("fortee login failed")
)

func Login(ctx context.Context, username string, password string) (string, error) {
	client, err := NewClientWithResponses(apiEndpoint, WithRequestEditorFn(addAcceptHeader))
	if err != nil {
		return "", err
	}
	res, err := client.PostLoginWithFormdataBodyWithResponse(ctx, PostLoginFormdataRequestBody{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	if res.StatusCode() != http.StatusOK {
		return "", ErrLoginFailed
	}
	resOk := res.JSON200
	if !resOk.LoggedIn {
		return "", ErrLoginFailed
	}
	if resOk.User == nil {
		return "", ErrLoginFailed
	}
	return resOk.User.Username, nil
}

// fortee API denies requests without Accept header.
func addAcceptHeader(_ context.Context, req *http.Request) error {
	req.Header.Set("Accept", "application/json")
	return nil
}
