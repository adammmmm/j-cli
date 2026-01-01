package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/adammmmm/go-junos"
	"github.com/spf13/cobra"
)

var opCmd = &cobra.Command{
	Use:   "op",
	Short: "Operational Commands",
	Long:  `Runs operational commands passed as a string, for example "show version"`,
	Run:   operationalCommand,
}

func operationalCommand(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("no operational command provided")
		return
	}

	ctx := context.Background()
	var progress *Progress
	if !jsonOutput {
		progress = NewProgress(len(devices))
	}
	command := strings.Join(args, " ")

	cmdResult := CommandResult{
		Command: "op " + command,
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
			out, err := jnpr.Command(command, "text")
			if err != nil {
				return err
			}
			output = out
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
	rootCmd.AddCommand(opCmd)
}
