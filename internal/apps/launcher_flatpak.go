package apps

import (
	"context"
	"fmt"
	"runtime"
)

// FlatpakLauncher launches flatpak apps.
type FlatpakLauncher struct{}

func (FlatpakLauncher) Launch(ctx context.Context, app App) (Process, error) {
	if runtime.GOOS == "windows" {
		return nil, fmt.Errorf("flatpak not available on Windows")
	}
	ref := app.Exec
	if ref == "" {
		ref = app.ID
	}
	return runCommand(ctx, "flatpak", "run", ref)
}
