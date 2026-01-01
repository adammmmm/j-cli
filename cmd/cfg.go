package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/adammmmm/go-junos"
	"github.com/spf13/cobra"
)

var cfgCmd = &cobra.Command{
	Use:   "cfg",
	Short: "Configuration Commands",
	Long:  `Runs configuration commands passed as a string, for example "set snmp location Stockholm"`,
	Run:   configCommand,
}

func configCommand(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("no configuration command provided")
		return
	}

	ctx := context.Background()
	var progress *Progress
	if !jsonOutput {
		progress = NewProgress(len(devices))
	}

	config := strings.Join(args, " ")
	cmdResult := CommandResult{
		Command: "cfg " + config,
		Workers: workers,
		Timeout: timeout.String(),
		Retries: retries,
	}

	results := runOnDevices(ctx, devices, workers, func(ctx context.Context, device string) workerResult {
		devCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		jnpr, err := junos.NewSession(device, auth)
		if err != nil {
			return workerResult{device: device, err: err}
		}
		defer jnpr.Close()

		err = withRetry(devCtx, retries, backoff, func() error {
			if resp := jnpr.Config(config, "set", true); resp != nil {
				return fmt.Errorf("%v", resp)
			}
			return nil
		})

		if err != nil {
			return workerResult{device: device, err: err}
		}

		return workerResult{device: device, output: "Configuration applied"}
	})

	for r := range results {
		if progress != nil {
			progress.Increment()
		}

		if r.err != nil {
			cmdResult.Results = append(cmdResult.Results, DeviceResult{
				Device: r.device,
				OK:     false,
				Error:  r.err.Error(),
			})
		} else {
			cmdResult.Results = append(cmdResult.Results, DeviceResult{
				Device: r.device,
				OK:     true,
				Output: normalizeOutput(r.output),
			})
		}
	}

	printResults(cmdResult)
}

func init() {
	rootCmd.AddCommand(cfgCmd)
}
