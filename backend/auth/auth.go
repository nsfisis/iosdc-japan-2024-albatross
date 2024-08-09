package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

var (
	ErrInvalidRegistrationToken = errors.New("invalid registration token")
	ErrNoRegistrationToken      = errors.New("no registration token")
	ErrForteeLoginFailed        = errors.New("fortee login failed")
)

func Login(
	ctx context.Context,
	queries *db.Queries,
	username string,
	password string,
	registrationToken *string,
) (int, error) {
	userAuth, err := queries.GetUserAuthByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err := signup(ctx, queries, username, password, registrationToken)
			if err != nil {
				return 0, err
			}
			return Login(ctx, queries, username, password, nil)
		}
		return 0, err
	}
	if userAuth.AuthType == "password" {
		passwordHash := userAuth.PasswordHash
		if passwordHash == nil {
			panic("inconsistant data")
		}
		err := bcrypt.CompareHashAndPassword([]byte(*passwordHash), []byte(password))
		if err != nil {
			return 0, err
		}
		return int(userAuth.UserID), nil
	} else if userAuth.AuthType == "fortee" {
		if err := verifyForteeAccount(ctx, username, password); err != nil {
			return 0, err
		}
		return int(userAuth.UserID), nil
	}
	panic(fmt.Sprintf("unexpected auth type: %s", userAuth.AuthType))
}

func signup(
	ctx context.Context,
	queries *db.Queries,
	username string,
	password string,
	registrationToken *string,
) error {
	if err := verifyRegistrationToken(ctx, queries, registrationToken); err != nil {
		return err
	}
	if err := verifyForteeAccount(ctx, username, password); err != nil {
		return err
	}

	// TODO: transaction
	userID, err := queries.CreateUser(ctx, username)
	if err != nil {
		return err
	}
	if err := queries.CreateUserAuth(ctx, db.CreateUserAuthParams{
		UserID:   userID,
		AuthType: "fortee",
	}); err != nil {
		return err
	}
	return nil
}

func verifyRegistrationToken(ctx context.Context, queries *db.Queries, registrationToken *string) error {
	if registrationToken == nil {
		return ErrNoRegistrationToken
	}
	exists, err := queries.IsRegistrationTokenValid(ctx, *registrationToken)
	if err != nil {
		return err
	}
	if !exists {
		return ErrInvalidRegistrationToken
	}
	return nil
}

func verifyForteeAccount(_ context.Context, username string, password string) error {
	reqData := url.Values{}
	reqData.Set("username", username)
	reqData.Set("password", password)
	reqBody := reqData.Encode()

	req, err := http.NewRequest("POST", "https://fortee.jp/api/user/login", bytes.NewBufferString(reqBody))
	if err != nil {
		return fmt.Errorf("http.NewRequest failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do failed: %v", err)
	}
	defer res.Body.Close()

	resData := struct {
		LoggedIn bool `json:"loggedIn"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return fmt.Errorf("json.Decode failed: %v", err)
	}

	if !resData.LoggedIn {
		return ErrForteeLoginFailed
	}
	return nil
}
