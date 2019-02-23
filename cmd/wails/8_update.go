package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

func init() {

	// var forceRebuild = false
	checkSpinner := spinner.NewSpinner()
	checkSpinner.SetSpinSpeed(50)

	commandDescription := `This command checks if there are updates to Wails.`
	updateCmd := app.Command("update", "Check for Updates.").
		LongDescription(commandDescription)

	updateCmd.Action(func() error {

		message := "Checking for updates..."
		logger.PrintSmallBanner(message)
		fmt.Println()

		// Get versions
		checkSpinner.Start(message)
		resp, err := http.Get("https://api.github.com/repos/wailsapp/wails/tags")
		if err != nil {
			checkSpinner.Error(err.Error())
			return err
		}
		checkSpinner.Success()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			checkSpinner.Error(err.Error())
			return err
		}

		data := []map[string]interface{}{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}

		latestVersion := data[0]["name"].(string)
		fmt.Println()
		fmt.Println("Current Version: " + cmd.Version)
		fmt.Println(" Latest Version: " + latestVersion)
		if latestVersion != cmd.Version {
			updateSpinner := spinner.NewSpinner()
			updateSpinner.SetSpinSpeed(40)
			updateSpinner.Start("Updating to  : " + latestVersion)
			err = cmd.NewProgramHelper().RunCommandArray([]string{"go", "get", "-u", "-a", "github.com/wailsapp/wails/cmd/wails"})
			if err != nil {
				updateSpinner.Error(err.Error())
				return err
			}
			version, _, err := cmd.NewProgramHelper().GetOutputFromCommand("wails version")
			version = strings.TrimSpace(version)
			if err != nil {
				updateSpinner.Error(err.Error())
				return err
			}
			if version != latestVersion {
				updateSpinner.Error()
				logger.Yellow("Weird! Wails was supposed to update to %s but seems to be at %s instead.", latestVersion, version)
			} else {
				updateSpinner.Success()
				logger.Green("Success! Wails updated to " + version)
			}
		} else {
			logger.Yellow("Looks like you're up to date!")
		}
		return nil
	})
}
