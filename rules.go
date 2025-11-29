package main

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/wiz"
	"github.com/google/uuid"
)

type PlexAccount struct {
	Id		  int 	 `json:"id"`
	Title	  string `json:"title"`
	Thumbnail string `json:"thumb"`
}

type PlexDevice struct {
	Id 		 int 	`json:"id"`
	Name   	 string `json:"name"`
	LastSeen string `json:"lastSeenAt"`
}

type PlexPayload struct {
	Event 	string 		`json:"event"`
	Account PlexAccount `json:"Account"`
	Player 	PlexDevice	`json:"Player"`
}

type WizAction struct {
	Command 	[]wiz.WizCommand `json:"command"`
	BulbsMac	[]string		 `json:"BulbsMac"`
}

type RuleCondition struct {
	Event 	[]string 		`json:"event"`
	Account []PlexAccount	`json:"account"`
	Player 	[]PlexDevice	`json:"player"`
}

type Rule struct {
	Id 		  string		`json:"ruleID"`
	Condition RuleCondition `json:"condition"`
	Action 	  []WizAction	`json:"action"`
}

func (cfg *config) loadRules() error {
	fmt.Println("Loading rules...")
	rules, err := cfg.db.GetAllRules(context.Background())
	if err != nil {
		return err
	}

	for _, rule := range rules {
		var tempCondition RuleCondition
		var tempWizAction []WizAction
		if err := json.Unmarshal([]byte(rule.Condition), &tempCondition); err != nil {
			fmt.Printf("Failed to unmarshal condition of rule %s\n", err.Error())
			continue
		}
		if err := json.Unmarshal([]byte(rule.Action), &tempWizAction); err != nil {
			fmt.Printf("Failed to unmarshal Wiz Action of rule %s\n", err.Error())
			continue
		}
		cfg.rules = append(cfg.rules, Rule{
			Id: 		rule.ID,
			Condition:  tempCondition,
			Action: 	tempWizAction,
		})
	}
	if len(cfg.rules) == 0 {
		fmt.Println("No rules in db")
	} else {
		fmt.Println("Rules loaded!")
	}
	return nil
}

func (cfg *config) addRule(event []string, account []PlexAccount, player []PlexDevice, actions []WizAction) error {
	newRule := Rule{}
	newRule.Id = uuid.New().String()
	newRule.Condition.Account = append(newRule.Condition.Account, account...)
	newRule.Condition.Event = append(newRule.Condition.Event, event...)
	newRule.Condition.Player = append(newRule.Condition.Player, player...)
	newRule.Action = append(newRule.Action, actions...)
	cfg.rules = append(cfg.rules, newRule)

	conditionData, err := json.Marshal(newRule.Condition)
	if err != nil {
		return fmt.Errorf("Failed to marshal rules: %s", err)
	}

	actionData, err := json.Marshal(newRule.Action)
	if err != nil {
		return fmt.Errorf("Failed to marshal rules: %s", err)
	}

	if _, err := cfg.db.AddRule(context.Background(), database.AddRuleParams{
		ID: 		newRule.Id,
		Condition:  string(conditionData),
		Action: 	string(actionData),
	}); err != nil {
		return fmt.Errorf("Failed to add rule to db: %s", err)
	}
	return nil
}

func (cfg *config) deleteRule(id string) error {
	for idx, rule := range cfg.rules {
		if rule.Id == id {
			cfg.rules = slices.Delete(cfg.rules, idx, idx + 1)
			return cfg.db.DeleteRule(context.Background(), id)
		}
	}
	return fmt.Errorf("Failed to fund rule with id %s", id)
}

func (cfg *config) triggersRule(payload PlexPayload) ([]WizAction, string) {
	for _, rule := range cfg.rules {
		if slices.Contains(rule.Condition.Event, payload.Event) {
			if slices.Contains(rule.Condition.Player, payload.Player) || len(rule.Condition.Player) == 0 {
				if slices.Contains(rule.Condition.Account, payload.Account) || len(rule.Condition.Account) == 0 {
					return rule.Action, rule.Id
				}
			}
		}
	}
	return nil, ""
}

func (cfg *config) getBulbByBulbMac(bulbMac string) *wiz.Bulb{
	for _, bulb := range cfg.bulbsMap {
		if bulb.Mac == bulbMac {
			return bulb
		}
	}
	return nil
}

func (cfg *config) executeActions(actions []WizAction) {
	for _, action := range actions {
		for _, cmd := range action.Command {
			for _, mac := range action.BulbsMac {
				bulb := cfg.getBulbByBulbMac(mac)
				if bulb != nil {
					bulb.Execute(cfg.conn, cmd)
				} else {
					fmt.Printf("Failed to find a bulb with id %s\n", mac)
				}
			}
		}
	}
}
