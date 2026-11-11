package apps

import "context"

// MPVLauncher launches media via mpv.
type MPVLauncher struct{}

func (MPVLauncher) Launch(ctx context.Context, app App) (Process, error) {
	args := []string{"--fullscreen", "--force-window=yes"}
	if app.URL != "" {
		args = append(args, app.URL)
	}
	return runCommand(ctx, "mpv", args...)
}
