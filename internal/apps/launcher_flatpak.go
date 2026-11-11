package apps

import "context"

// FlatpakLauncher launches flatpak apps.
type FlatpakLauncher struct{}

func (FlatpakLauncher) Launch(ctx context.Context, app App) (Process, error) {
	ref := app.Exec
	if ref == "" {
		ref = app.ID
	}
	return runCommand(ctx, "flatpak", "run", ref)
}
