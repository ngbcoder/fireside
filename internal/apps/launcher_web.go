package apps

import (
	"context"
	"fmt"
)

// WebLauncher launches web apps in kiosk Chromium.
type WebLauncher struct {
	ChromiumBin string // path to chromium binary; defaults to "chromium"
}

func (w WebLauncher) bin() string {
	if w.ChromiumBin != "" {
		return w.ChromiumBin
	}
	return "chromium"
}

func (w WebLauncher) Launch(ctx context.Context, app App) (Process, error) {
	if app.URL == "" {
		return nil, fmt.Errorf("web app %s has no URL", app.ID)
	}
	args := []string{
		"--app=" + app.URL,
		"--kiosk",
		"--ozone-platform=wayland",
	}
	if app.UserAgent != "" {
		args = append(args, "--user-agent="+app.UserAgent)
	}
	if app.ProfileDir != "" {
		args = append(args, "--user-data-dir="+app.ProfileDir)
	}
	return runCommand(ctx, w.bin(), args...)
}
