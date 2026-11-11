package apps

import (
	"context"
	"runtime"
)

// MPVLauncher launches media via mpv.
type MPVLauncher struct{}

func (MPVLauncher) Launch(ctx context.Context, app App) (Process, error) {
	bin := "mpv"
	if runtime.GOOS == "windows" {
		bin = "mpv.exe"
	}
	args := []string{"--fullscreen", "--force-window=yes"}
	if app.URL != "" {
		args = append(args, app.URL)
	}
	return runCommand(ctx, bin, args...)
}
