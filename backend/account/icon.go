package account

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/fortee"
)

func FetchIcon(
	ctx context.Context,
	q *db.Queries,
	userID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Fetch user.
	user, err := q.GetUserByID(ctx, int32(userID))
	if err != nil {
		return fmt.Errorf("failed to fetch user icon (uid=%d): %w", userID, err)
	}
	// Fetch user icon URL.
	avatarURL, err := fortee.GetUserAvatarURL(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("failed to fetch user icon (uid=%d): %w", userID, err)
	}
	// Download user icon file.
	filePath := fmt.Sprintf("/files/img/%s/icon%s", url.PathEscape(user.Username), path.Ext(avatarURL))
	if err := downloadFile(ctx, fortee.Endpoint+avatarURL, "/data"+filePath); err != nil {
		return fmt.Errorf("failed to fetch user icon (uid=%d): %w", userID, err)
	}
	// Save user icon path.
	if err := q.UpdateUserIconPath(ctx, db.UpdateUserIconPathParams{
		UserID:   int32(userID),
		IconPath: &filePath,
	}); err != nil {
		return fmt.Errorf("failed to fetch user icon (uid=%d): %w", userID, err)
	}
	return nil
}

func downloadFile(ctx context.Context, url string, filePath string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to download file (%s): %w", url, err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file (%s): %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file (%s): status %d", url, res.StatusCode)
	}

	fileDir := filepath.Dir(filePath)
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory (%s): %w", fileDir, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file (%s): %w", filePath, err)
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return fmt.Errorf("failed to save file (%s): %w", filePath, err)
	}

	return nil
}
