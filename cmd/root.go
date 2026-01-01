package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/adammmmm/go-junos"
	"github.com/spf13/cobra"
)

type DeviceList struct {
	IPAddressList []string `json:"devices"`
}

var (
	authFile   string
	deviceFile string
	auth       *junos.AuthMethod
	devices    []string
	workers    int
	timeout    time.Duration

	retries    int
	backoff    time.Duration
	jsonOutput bool
)
var rootCmd = &cobra.Command{
	Use:   "j-cli",
	Short: "Junos Tool for Commands and Configuration",
	Long:  `A CLI tool for Junos devices, capable of configuration changes, configuration backups, operational commands and more.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		auth, err = readAuthJson(authFile)
		if err != nil {
			return err
		}

		devices, err = readDeviceJson(deviceFile)
		if err != nil {
			return err
		}

		if len(devices) == 0 {
			return fmt.Errorf("no devices found")
		}
		if workers < 1 {
			return fmt.Errorf("--workers must be >= 1")
		}
		if timeout <= 0 {
			return fmt.Errorf("--timeout must be > 0")
		}
		if retries < 0 {
			return fmt.Errorf("--retries must be >= 0")
		}
		if backoff <= 0 {
			return fmt.Errorf("--backoff must be > 0")
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func readDeviceJson(deviceFile string) ([]string, error) {
	jsonFile, err := os.Open(deviceFile)
	if err != nil {
		return nil, fmt.Errorf("error opening device file: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading device file: %v", err)
	}

	var deviceList DeviceList
	if err := json.Unmarshal(byteValue, &deviceList); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return deviceList.IPAddressList, nil
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&deviceFile,
		"devices",
		"d",
		"devices.json",
		"Device File in json format",
	)

	rootCmd.PersistentFlags().StringVarP(
		&authFile,
		"authentication",
		"a",
		"auth.json",
		"Authentication file in json format",
	)

	rootCmd.PersistentFlags().IntVarP(
		&workers,
		"workers",
		"w",
		10,
		"Number of concurrent device workers",
	)
	rootCmd.PersistentFlags().DurationVarP(
		&timeout,
		"timeout",
		"t",
		30*time.Second,
		"Timeout per device operation",
	)
	rootCmd.PersistentFlags().IntVar(
		&retries,
		"retries",
		2,
		"Number of retries per device on failure",
	)

	rootCmd.PersistentFlags().DurationVar(
		&backoff,
		"backoff",
		2*time.Second,
		"Initial backoff duration between retries",
	)
	rootCmd.PersistentFlags().BoolVar(
		&jsonOutput,
		"json",
		false,
		"Output results as JSON",
	)

}
