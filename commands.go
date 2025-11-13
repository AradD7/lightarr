package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/AradD7/lightarr/internal/wiz"
)

func commandExit(cfg *config, param ...string) error {
	fmt.Println("Closing Lightarr... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, param ...string) error {
	fmt.Printf("Welcome to the Lightarr!\n\nUsage:\n\n")
	for key, value := range getCommands() {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	return nil
}

func commandPrintBulbs(cfg *config, param ...string) error {
	for _, bulb := range cfg.bulbsMap {
		fmt.Println(*bulb)
	}
	return nil
}

func commandTurnBulbs(cfg *config, param ...string) error {
	if len(param) < 3 {
		return fmt.Errorf("Turn command take at least 3 inputs. turn [on/off] [...bulb ids/bulb names/all]")
	}
	wizCommand := wiz.NewSetPilotWizCommand()
	switch param[1] {
	case "off":
		wizCommand.TurnOff()
	case "on":
		wizCommand.TurnOn()
	default:
		return fmt.Errorf("Second argument for Turn command is either off or on")
	}

	for _, bulb := range cfg.bulbsMap {
		if param[2] == "all" || slices.Contains(param, bulb.Id) {
			err := bulb.Execute(cfg.conn, wizCommand)
			if err != nil {
				fmt.Printf("Failed to turn %s bulb with id %s\n", bulb.Id, param[1])
			}
		}
	}
	return nil
}
