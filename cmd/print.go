package cmd

import (
	"encoding/json"
	"fmt"
)

func printResults(cmdResult CommandResult) {
	if jsonOutput {
		enc, _ := json.MarshalIndent(cmdResult, "", "  ")
		fmt.Println(string(enc))
		return
	}

	for _, r := range cmdResult.Results {
		fmt.Println("****************************")
		fmt.Println(r.Device)
		fmt.Println("****************************")

		if r.OK {
			for _, line := range r.Output {
				fmt.Println(line)
			}
		} else {
			fmt.Println(r.Error)
		}
	}
}
