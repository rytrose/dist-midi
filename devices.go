package util

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"gitlab.com/gomidi/midi/mid"
)

// PromptForInput returns an input MIDI devices a user selects.
func PromptForInput(driver mid.Driver) (mid.In, error) {
	inputPorts, err := driver.Ins()
	if err != nil {
		return nil, fmt.Errorf("unable to get input ports from driver: %w", err)
	}

	var portNames []string
	for _, port := range inputPorts {
		portNames = append(portNames, port.String())
	}

	prompt := promptui.Select{
		Label: "Choose an Input MIDI Device",
		Items: portNames,
	}

	portIndex, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("unable to prompt for ports: %w", err)
	}

	return inputPorts[portIndex], nil
}

// PromptForOutput returns an output MIDI devices a user selects.
func PromptForOutput(driver mid.Driver) (mid.Out, error) {
	outputPorts, err := driver.Outs()
	if err != nil {
		return nil, fmt.Errorf("unable to get output ports from driver: %w", err)
	}

	var portNames []string
	for _, port := range outputPorts {
		portNames = append(portNames, port.String())
	}

	prompt := promptui.Select{
		Label: "Choose an Output MIDI Device",
		Items: portNames,
	}

	portIndex, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("unable to prompt for ports: %w", err)
	}

	return outputPorts[portIndex], nil
}
