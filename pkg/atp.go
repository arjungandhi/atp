package atp

// ATP is the main object that contians the automated task planning state
type ATP struct {
	Tasks         []Task `json:"tasks"`
	ActiveTaskDir string `json:"active_task_dir"`
	StorageDir    string `json:"storage_dir"`
}
