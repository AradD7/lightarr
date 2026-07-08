package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"

	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/plex"
	"github.com/AradD7/lightarr/internal/wiz"
	"github.com/google/uuid"
)

func LoadRulesFromDb(db *database.Queries, logger *slog.Logger) ([]Rule, error) {
	rules := []Rule{}
	logger.Info("Loading rules from database...")
	rulesInDb, err := db.GetAllRules(context.Background())
	if err != nil {
		return nil, err
	}

	for _, rule := range rulesInDb {
		var tempCondition RuleCondition
		var tempWizAction []wiz.WizAction
		if err := json.Unmarshal([]byte(rule.Condition), &tempCondition); err != nil {
			logger.Info(fmt.Sprintf("Failed to unmarshal condition of rule %s", err.Error()))
			continue
		}
		if err := json.Unmarshal([]byte(rule.Action), &tempWizAction); err != nil {
			logger.Info(fmt.Sprintf("Failed to unmarshal Wiz Action of rule %s", err.Error()))
			continue
		}
		rules = append(rules, Rule{
			Name:      rule.Name.String,
			Id:        rule.ID,
			Condition: tempCondition,
			Action:    tempWizAction,
		})
	}
	if len(rules) == 0 {
		logger.Info("No rules are in database yet")
	} else {
		logger.Info("Rules loaded!")
	}
	return rules, nil
}

func AddRule(currentRules []Rule, event []string, account []plex.PlexAccount, player []plex.PlexDevice, actions []wiz.WizAction, db *database.Queries) ([]Rule, error) {
	newRule := Rule{}
	newRule.Id = uuid.New().String()
	newRule.Condition.Account = append(newRule.Condition.Account, account...)
	newRule.Condition.Event = append(newRule.Condition.Event, event...)
	newRule.Condition.Device = append(newRule.Condition.Device, player...)
	newRule.Action = append(newRule.Action, actions...)
	currentRules = append(currentRules, newRule)

	conditionData, err := json.Marshal(newRule.Condition)
	if err != nil {
		return currentRules, fmt.Errorf("Failed to marshal rules: %s", err)
	}

	actionData, err := json.Marshal(newRule.Action)
	if err != nil {
		return currentRules, fmt.Errorf("Failed to marshal rules: %s", err)
	}

	if _, err := db.AddRule(context.Background(), database.AddRuleParams{
		ID:        newRule.Id,
		Condition: string(conditionData),
		Action:    string(actionData),
	}); err != nil {
		return nil, fmt.Errorf("Failed to add rule to db: %s", err)
	}
	return currentRules, nil
}

func DeleteRule(currentRules []Rule, id string, db *database.Queries) ([]Rule, error) {
	for idx, rule := range currentRules {
		if rule.Id == id {
			currentRules = slices.Delete(currentRules, idx, idx+1)
			return currentRules, db.DeleteRule(context.Background(), id)
		}
	}
	return currentRules, fmt.Errorf("Failed to fund rule with id %s", id)
}

func TriggersRule(rules []Rule, payload plex.PlexPayload) ([]wiz.WizAction, string) {
	for _, rule := range rules {
		if slices.Contains(rule.Condition.Event, payload.Event) && isPayloadAccountInRule(rule, payload) && isPayloadDeviceInRule(rule, payload) {
			return rule.Action, rule.Id
		}
	}
	return nil, ""
}

func isPayloadDeviceInRule(rule Rule, payload plex.PlexPayload) bool {
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

func isPayloadAccountInRule(rule Rule, payload plex.PlexPayload) bool {
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
