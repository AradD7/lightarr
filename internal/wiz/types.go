package wiz

type WizAction struct {
	Command  WizCommand `json:"command"`
	BulbsMac []string   `json:"bulbsMac"`
}
