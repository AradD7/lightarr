package config

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/plex"
	"github.com/AradD7/lightarr/internal/rules"
	"github.com/AradD7/lightarr/internal/wiz"
	"github.com/google/uuid"
)

func (cfg *Config) LoadRules() error {
	cfg.Logger.Info("Loading rules...")
	rulesInDb, err := cfg.Db.GetAllRules(context.Background())
	if err != nil {
		return err
	}

	for _, rule := range rulesInDb {
		var tempCondition rules.RuleCondition
		var tempWizAction []wiz.WizAction
		if err := json.Unmarshal([]byte(rule.Condition), &tempCondition); err != nil {
			cfg.Logger.Info(fmt.Sprintf("Failed to unmarshal condition of rule %s", err.Error()))
			continue
		}
		if err := json.Unmarshal([]byte(rule.Action), &tempWizAction); err != nil {
			cfg.Logger.Info(fmt.Sprintf("Failed to unmarshal Wiz Action of rule %s", err.Error()))
			continue
		}
		cfg.Rules = append(cfg.Rules, rules.Rule{
			Name:      rule.Name.String,
			Id:        rule.ID,
			Condition: tempCondition,
			Action:    tempWizAction,
		})
	}
	if len(cfg.Rules) == 0 {
		cfg.Logger.Info("No rules in db")
	} else {
		cfg.Logger.Info("Rules loaded!")
	}
	return nil
}

func (cfg *Config) AddRule(event []string, account []plex.PlexAccount, player []plex.PlexDevice, actions []wiz.WizAction) error {
	newRule := rules.Rule{}
	newRule.Id = uuid.New().String()
	newRule.Condition.Account = append(newRule.Condition.Account, account...)
	newRule.Condition.Event = append(newRule.Condition.Event, event...)
	newRule.Condition.Device = append(newRule.Condition.Device, player...)
	newRule.Action = append(newRule.Action, actions...)
	cfg.Rules = append(cfg.Rules, newRule)

	conditionData, err := json.Marshal(newRule.Condition)
	if err != nil {
		return fmt.Errorf("Failed to marshal rules: %s", err)
	}

	actionData, err := json.Marshal(newRule.Action)
	if err != nil {
		return fmt.Errorf("Failed to marshal rules: %s", err)
	}

	if _, err := cfg.Db.AddRule(context.Background(), database.AddRuleParams{
		ID:        newRule.Id,
		Condition: string(conditionData),
		Action:    string(actionData),
	}); err != nil {
		return fmt.Errorf("Failed to add rule to db: %s", err)
	}
	return nil
}

func (cfg *Config) DeleteRule(id string) error {
	for idx, rule := range cfg.Rules {
		if rule.Id == id {
			cfg.Rules = slices.Delete(cfg.Rules, idx, idx+1)
			return cfg.Db.DeleteRule(context.Background(), id)
		}
	}
	return fmt.Errorf("Failed to fund rule with id %s", id)
}

func (cfg *Config) TriggersRule(payload plex.PlexPayload) ([]wiz.WizAction, string) {
	for _, rule := range cfg.Rules {
		if slices.Contains(rule.Condition.Event, payload.Event) && isPayloadAccountInRule(rule, payload) && isPayloadDeviceInRule(rule, payload) {
			return rule.Action, rule.Id
		}
	}
	return nil, ""
}

func isPayloadDeviceInRule(rule rules.Rule, payload plex.PlexPayload) bool {
	if len(rule.Condition.Device) == 0 {
		return true
	}
	for _, device := range rule.Condition.Device {
		if device.Id == payload.Player.Uuid {
			return true
		}
	}
	return false
}

func isPayloadAccountInRule(rule rules.Rule, payload plex.PlexPayload) bool {
	if len(rule.Condition.Account) == 0 {
		return true
	}
	for _, account := range rule.Condition.Account {
		if account.Id == payload.Account.Id {
			return true
		}
	}
	return false
}

func (cfg *Config) getBulbByBulbMac(bulbMac string) *wiz.Bulb {
	for _, bulb := range cfg.BulbsMap {
		if bulb.Mac == bulbMac {
			return bulb
		}
	}
	return nil
}

func (cfg *Config) ExecuteActions(actions []wiz.WizAction) {
	for _, action := range actions {
		for _, mac := range action.BulbsMac {
			bulb := cfg.getBulbByBulbMac(mac)
			if bulb != nil {
				bulb.Execute(cfg.Conn, action.Command)
			} else {
				cfg.Logger.Info(fmt.Sprintf("Failed to find a bulb with id %s", mac))
			}
		}
	}
}
