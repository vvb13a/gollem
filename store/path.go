package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	appDirName = "app"
	dbFileName = "app.db"
)

// DataDir returns the directory that holds the SQLite database. APP_DATA_DIR
// overrides the platform default, which is useful for local development and
// container volumes:
//
//	macOS   ~/Library/Application Support/app
//	Linux   $XDG_DATA_HOME/app (~/.local/share/app)
//	Windows %LOCALAPPDATA%\app
func DataDir() (string, error) {
	if dir := strings.TrimSpace(os.Getenv("APP_DATA_DIR")); dir != "" {
		return dir, nil
	}
	switch runtime.GOOS {
	case "windows":
		if dir := os.Getenv("LOCALAPPDATA"); dir != "" {
			return filepath.Join(dir, appDirName), nil
		}
		return "", errors.New("LOCALAPPDATA is not set")
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home dir: %w", err)
		}
		return filepath.Join(home, "Library", "Application Support", appDirName), nil
	default:
		if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
			return filepath.Join(dir, appDirName), nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home dir: %w", err)
		}
		return filepath.Join(home, ".local", "share", appDirName), nil
	}
}

// DatabasePath returns the full path of the SQLite database file, creating its
// parent directory if needed.
func DatabasePath() (string, error) {
	dir, err := DataDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create data dir: %w", err)
	}
	return filepath.Join(dir, dbFileName), nil
}
