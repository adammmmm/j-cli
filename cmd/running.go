package cmd

import (
	"context"

	"github.com/adammmmm/go-junos"
	"github.com/spf13/cobra"
)

var runningCmd = &cobra.Command{
	Use:   "running",
	Short: "Show running configuration",
	Long:  `Show running configuration`,
	Run:   runCommand,
}

func runCommand(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	var progress *Progress
	if !jsonOutput {
		progress = NewProgress(len(devices))
	}

	cmdResult := CommandResult{
		Command: "running",
		Workers: workers,
		Timeout: timeout.String(),
		Retries: retries,
	}

	results := runOnDevices(ctx, devices, workers, func(ctx context.Context, device string) workerResult {
		var config string

		devCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		jnpr, err := junos.NewSession(device, auth)
		if err != nil {
			return workerResult{device: device, err: err}
		}
		defer jnpr.Close()

		err = withRetry(devCtx, retries, backoff, func() error {
			cfg, err := jnpr.GetConfig("text")
			if err != nil {
				return err
			}
			config = cfg
			return nil
		})

		if err != nil {
			return workerResult{device: device, err: err}
		}

		return workerResult{device: device, output: config}
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
	rootCmd.AddCommand(runningCmd)
}
