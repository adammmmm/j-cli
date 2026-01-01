package cmd

type DeviceResult struct {
	Device string   `json:"device"`
	OK     bool     `json:"ok"`
	Output []string `json:"output,omitempty"`
	Error  string   `json:"error,omitempty"`
}

type CommandResult struct {
	Command string         `json:"command"`
	Workers int            `json:"workers"`
	Timeout string         `json:"timeout"`
	Retries int            `json:"retries"`
	Results []DeviceResult `json:"results"`
}
