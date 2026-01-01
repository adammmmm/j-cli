package cmd

import (
	"context"

	"github.com/adammmmm/go-junos"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show difference between current and previous commit",
	Long:  `Show difference between current and previous commit`,
	Run:   diffConfigCommand,
}

func diffConfigCommand(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	var progress *Progress
	if !jsonOutput {
		progress = NewProgress(len(devices))
	}
	cmdResult := CommandResult{
		Command: "diff",
		Workers: workers,
		Timeout: timeout.String(),
		Retries: retries,
	}

	results := runOnDevices(ctx, devices, workers, func(ctx context.Context, device string) workerResult {
		var output string

		devCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		jnpr, err := junos.NewSession(device, auth)
		if err != nil {
			return workerResult{device: device, err: err}
		}
		defer jnpr.Close()

		err = withRetry(devCtx, retries, backoff, func() error {
			diff, err := jnpr.Diff(1)
			if err != nil {
				return err
			}
			output = diff
			return nil
		})

		if err != nil {
			return workerResult{device: device, err: err}
		}

		return workerResult{device: device, output: output}
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
	rootCmd.AddCommand(diffCmd)
}
