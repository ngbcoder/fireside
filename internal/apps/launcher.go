package apps

import "context"

// CompositeLauncher dispatches to the appropriate launcher by AppKind.
type CompositeLauncher struct {
	Web     WebLauncher
	MPV     MPVLauncher
	Flatpak FlatpakLauncher
}

func (c CompositeLauncher) Launch(ctx context.Context, a App) (Process, error) {
	switch a.Kind {
	case AppWeb:
		return c.Web.Launch(ctx, a)
	case AppFlatpak:
		return c.Flatpak.Launch(ctx, a)
	case AppNative:
		// Native apps are executed directly via Exec.
		return runCommand(ctx, a.Exec)
	}
	return nil, nil
}
