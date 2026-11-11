package apps

// Folder groups apps on the home row.
type Folder struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	AppIDs []string `json:"app_ids"`
}

// Folders holds the set of folders.
type Folders struct {
	Items []Folder `json:"folders"`
}

// AppsInFolder returns the apps belonging to a folder.
func (f Folder) AppsInFolder(reg Registry) []App {
	var out []App
	for _, id := range f.AppIDs {
		if a, ok := reg.Find(id); ok {
			out = append(out, a)
		}
	}
	return out
}
