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

func LoginFortee(ctx context.Context, username string, password string) error {
	client, err := NewClientWithResponses(apiEndpoint, WithRequestEditorFn(addAcceptHeader))
	if err != nil {
		return err
	}
	res, err := client.PostLoginWithFormdataBodyWithResponse(ctx, PostLoginFormdataRequestBody{
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return ErrLoginFailed
	}
	if !res.JSON200.LoggedIn {
		return ErrLoginFailed
	}
	return nil
}

// fortee API denies requests without Accept header.
func addAcceptHeader(_ context.Context, req *http.Request) error {
	req.Header.Set("Accept", "application/json")
	return nil
}
