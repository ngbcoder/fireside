package apps

import "context"

// AppKind describes how an app is launched.
type AppKind int

const (
	AppWeb AppKind = iota
	AppFlatpak
	AppNative
)

// App is anything the shell can launch.
type App struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Kind       AppKind `json:"kind"`
	URL        string  `json:"url,omitempty"`        // web apps
	Exec       string  `json:"exec,omitempty"`       // native apps / flatpak ref
	Icon       string  `json:"icon,omitempty"`
	UserAgent  string  `json:"user_agent,omitempty"` // optional override
	ProfileDir string  `json:"profile_dir,omitempty"` // per-app browser profile
}

// Process represents a running app process.
type Process interface {
	Wait() error
	Stop() error
}

// Launcher launches apps.
type Launcher interface {
	Launch(ctx context.Context, app App) (Process, error)
}
