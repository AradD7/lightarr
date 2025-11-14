package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/AradD7/lightarr/internal/wiz"
	"github.com/google/uuid"
)

const rulesCache = "rules.json"

type PlexAccount struct {
	Id		int 	`json:"id"`
	Title	string 	`json:"title"`
}

type PlexPlayer struct {
	PublicAddr 	string `json:"publicAddress"`
	Title 		string `json:"title"`
	Uuid 		string `json:"uuid"`
}

type PlexPayload struct {
	Event 	string 		`json:"event"`
	Account PlexAccount `json:"Account"`
	Player 	PlexPlayer	`json:"Player"`
}

type WizAction struct {
	Command 	[]wiz.WizCommand `json:"command"`
	BulbId		[]string		 `json:"BulbIds"`
}

type Rule struct {
	Id 		  uuid.UUID		`json:"ruleID"`
	Condition struct {
		Event 	[]string 		`json:"event"`
		Account []PlexAccount	`json:"account"`
		Player 	[]PlexPlayer	`json:"player"`
	} `json:"condition"`
	Action 	  []WizAction	`json:"action"`
}

func (cfg *config) loadRules() {
	data, err := os.ReadFile(rulesCache)
	if err != nil {
		fmt.Printf("Failed to read rules cache file: %s\n", err.Error())
		return
	}

	err = json.Unmarshal(data, &cfg.rules)
	if err != nil {
		fmt.Printf("Failed to unmarshal rules: %s\n", err)
		return
	}
}

func (cfg *config) addRule(event []string, account []PlexAccount, player []PlexPlayer, actions []WizAction) {
	newRule := Rule{}
	newRule.Id = uuid.New()
	newRule.Condition.Account = append(newRule.Condition.Account, account...)
	newRule.Condition.Event = append(newRule.Condition.Event, event...)
	newRule.Condition.Player = append(newRule.Condition.Player, player...)
	newRule.Action = append(newRule.Action, actions...)
	cfg.rules = append(cfg.rules, newRule)

	data, err := json.MarshalIndent(cfg.rules, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal rules: %s\n", err)
		return
	}

	err = os.WriteFile(rulesCache, data, 0644)
	if err != nil {
		fmt.Printf("Failed to write to %s file: %s\n", rulesCache, err)
		return
	}
}

func (cfg *config) deleteRule(id uuid.UUID) error {
	for idx, rule := range cfg.rules {
		if rule.Id == id {
			cfg.rules = slices.Delete(cfg.rules, idx, idx + 1)
			return nil
		}
	}
	return fmt.Errorf("Failed to fund rule with id %s", id.String())
}

func (cfg *config) triggersRule(payload PlexPayload) []WizAction {
	for _, rule := range cfg.rules {
		if slices.Contains(rule.Condition.Event, payload.Event) {
			if slices.Contains(rule.Condition.Player, payload.Player) || len(rule.Condition.Player) == 0 {
				if slices.Contains(rule.Condition.Account, payload.Account) || len(rule.Condition.Account) == 0 {
					return rule.Action
				}
			}
		}
	}
	return nil
}

func (cfg *config) getBulbByBulbId(bulbId string) *wiz.Bulb{
	for _, bulb := range cfg.bulbsMap {
		if bulb.Id == bulbId {
			return bulb
		}
	}
	return nil
}

func (cfg *config) executeActions(actions []WizAction) {
	for _, action := range actions {
		for _, cmd := range action.Command {
			for _, id := range action.BulbId {
				bulb := cfg.getBulbByBulbId(id)
				if bulb != nil {
					bulb.Execute(cfg.conn, cmd)
				} else {
					fmt.Printf("Failed to find a bulb with id %s\n", id)
				}
			}
		}
	}
}
