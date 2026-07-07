package rules

import (
	"github.com/AradD7/lightarr/internal/plex"
	"github.com/AradD7/lightarr/internal/wiz"
)

type RuleCondition struct {
	Event   []string           `json:"event"`
	Account []plex.PlexAccount `json:"account"`
	Device  []plex.PlexDevice  `json:"device"`
}

type Rule struct {
	Name      string          `json:"name"`
	Id        string          `json:"ruleID"`
	Condition RuleCondition   `json:"condition"`
	Action    []wiz.WizAction `json:"action"`
}
