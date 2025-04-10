package dekopin_test

import (
	"context"
	"testing"

	"github.com/iwashi623/dekopin"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCmdOption_Validate(t *testing.T) {
	type TestResult struct {
		Err error
	}

	type ArrangeResult struct {
		option *dekopin.CmdOption
	}

	makeOption := func() *dekopin.CmdOption {
		runner := lo.Sample(dekopin.ValidRunners)

		return &dekopin.CmdOption{
			Project: "my-project",
			Region:  "us-central1",
			Service: "my-service",
			Runner:  runner,
		}
	}
	cases := map[string]TestCase[any, ArrangeResult, TestResult]{
		"success_case": {
			Arrange: func() ArrangeResult {
				option := makeOption()
				return ArrangeResult{
					option: option,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.NoError(t, result.Err)
			},
		},
		"error_empty_project": {
			Arrange: func() ArrangeResult {
				option := makeOption()
				option.Project = ""
				return ArrangeResult{
					option: option,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Error(t, result.Err)
			},
		},
		"error_empty_region": {
			Arrange: func() ArrangeResult {
				option := makeOption()
				option.Region = ""
				return ArrangeResult{
					option: option,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Error(t, result.Err)
			},
		},
		"error_empty_service": {
			Arrange: func() ArrangeResult {
				option := makeOption()
				option.Service = ""
				return ArrangeResult{
					option: option,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Error(t, result.Err)
			},
		},
		"error_empty_runner": {
			Arrange: func() ArrangeResult {
				option := makeOption()
				option.Runner = ""
				return ArrangeResult{
					option: option,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Error(t, result.Err)
			},
		},
		"error_invalid_runner_value": {
			Arrange: func() ArrangeResult {
				option := makeOption()
				option.Runner = "invalid-runner"
				return ArrangeResult{
					option: option,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Error(t, result.Err)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ar := c.Arrange()
			err := ar.option.Validate()
			c.Assert(t, ar, TestResult{
				Err: err,
			})
		})
	}
}

func TestGetCmdOption(t *testing.T) {
	type TestResult struct {
		Option *dekopin.CmdOption
		Err    error
	}

	type ArrangeResult struct {
		ctx context.Context
	}

	cases := map[string]TestCase[any, ArrangeResult, TestResult]{
		"success_option_exists_in_context": {
			Arrange: func() ArrangeResult {
				option := &dekopin.CmdOption{
					Project: "test-project",
					Region:  "test-region",
					Service: "test-service",
					Runner:  dekopin.RUNNER_GITHUB_ACTIONS,
				}
				ctx := dekopin.SetCmdOption(context.Background(), option)
				return ArrangeResult{
					ctx: ctx,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.NoError(t, result.Err)
				assert.NotNil(t, result.Option)
				assert.Equal(t, "test-project", result.Option.Project)
				assert.Equal(t, "test-region", result.Option.Region)
				assert.Equal(t, "test-service", result.Option.Service)
				assert.Equal(t, dekopin.RUNNER_GITHUB_ACTIONS, result.Option.Runner)
			},
		},
		"error_option_not_in_context": {
			Arrange: func() ArrangeResult {
				return ArrangeResult{
					ctx: context.Background(),
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Error(t, result.Err)
				assert.Nil(t, result.Option)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ar := c.Arrange()
			option, err := dekopin.GetCmdOption(ar.ctx)
			c.Assert(t, ar, TestResult{
				Option: option,
				Err:    err,
			})
		})
	}
}

func TestNewCmdOption(t *testing.T) {
	type TestResult struct {
		Option *dekopin.CmdOption
		Err    error
	}

	type ArrangeResult struct {
		ctx    context.Context
		config *dekopin.DekopinConfig
		cmd    *cobra.Command
		runner string
	}

	cases := map[string]TestCase[any, ArrangeResult, TestResult]{
		"success_using_flags": {
			Arrange: func() ArrangeResult {
				cmd := &cobra.Command{}
				runner := lo.Sample(dekopin.ValidRunners)
				cmd.Flags().String("project", "flag-project", "")
				cmd.Flags().String("region", "flag-region", "")
				cmd.Flags().String("service", "flag-service", "")
				cmd.Flags().String("runner", runner, "")

				ctx := dekopin.SetDekopinCommand(context.Background(), dekopin.NewDekopinCommand(cmd))

				return ArrangeResult{
					ctx:    ctx,
					config: nil,
					cmd:    cmd,
					runner: runner,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.NoError(t, result.Err)
				assert.NotNil(t, result.Option)
				assert.Equal(t, "flag-project", result.Option.Project)
				assert.Equal(t, "flag-region", result.Option.Region)
				assert.Equal(t, "flag-service", result.Option.Service)
				assert.Equal(t, assertArgs.runner, result.Option.Runner)
			},
		},
		"success_using_config": {
			Arrange: func() ArrangeResult {
				cmd := &cobra.Command{}
				runner := lo.Sample(dekopin.ValidRunners)
				cmd.Flags().String("project", "", "")
				cmd.Flags().String("region", "", "")
				cmd.Flags().String("service", "", "")
				cmd.Flags().String("runner", "", "")

				config := &dekopin.DekopinConfig{
					Project: "config-project",
					Region:  "config-region",
					Service: "config-service",
					Runner:  runner,
				}

				ctx := dekopin.SetDekopinCommand(context.Background(), dekopin.NewDekopinCommand(cmd))

				return ArrangeResult{
					ctx:    ctx,
					config: config,
					cmd:    cmd,
					runner: runner,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.NoError(t, result.Err)
				assert.NotNil(t, result.Option)
				assert.Equal(t, "config-project", result.Option.Project)
				assert.Equal(t, "config-region", result.Option.Region)
				assert.Equal(t, "config-service", result.Option.Service)
				assert.Equal(t, assertArgs.runner, result.Option.Runner)
			},
		},
		"success_flag_overrides_config": {
			Arrange: func() ArrangeResult {
				runner := lo.Sample(dekopin.ValidRunners)
				cmd := &cobra.Command{}
				cmd.Flags().String("project", "flag-project", "")
				cmd.Flags().String("region", "flag-region", "")
				cmd.Flags().String("service", "flag-service", "")
				cmd.Flags().String("runner", runner, "")

				config := &dekopin.DekopinConfig{
					Project: "config-project",
					Region:  "config-region",
					Service: "config-service",
					Runner:  "config-runner",
				}

				ctx := dekopin.SetDekopinCommand(context.Background(), dekopin.NewDekopinCommand(cmd))

				return ArrangeResult{
					ctx:    ctx,
					config: config,
					cmd:    cmd,
					runner: runner,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.NoError(t, result.Err)
				assert.NotNil(t, result.Option)
				assert.Equal(t, "flag-project", result.Option.Project)
				assert.Equal(t, "flag-region", result.Option.Region)
				assert.Equal(t, "flag-service", result.Option.Service)
				assert.Equal(t, assertArgs.runner, result.Option.Runner)
			},
		},
		"error_missing_required_values": {
			Arrange: func() ArrangeResult {
				cmd := &cobra.Command{}
				cmd.Flags().String("project", "", "")
				cmd.Flags().String("region", "", "")
				cmd.Flags().String("service", "", "")
				cmd.Flags().String("runner", "", "")

				ctx := dekopin.SetDekopinCommand(context.Background(), dekopin.NewDekopinCommand(cmd))

				return ArrangeResult{
					ctx:    ctx,
					config: nil,
					cmd:    cmd,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Error(t, result.Err)
				assert.Nil(t, result.Option)
			},
		},
		"error_invalid_runner": {
			Arrange: func() ArrangeResult {
				cmd := &cobra.Command{}
				cmd.Flags().String("project", "test-project", "")
				cmd.Flags().String("region", "test-region", "")
				cmd.Flags().String("service", "test-service", "")
				cmd.Flags().String("runner", "invalid-runner", "")

				ctx := dekopin.SetDekopinCommand(context.Background(), dekopin.NewDekopinCommand(cmd))

				return ArrangeResult{
					ctx:    ctx,
					config: nil,
					cmd:    cmd,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Error(t, result.Err)
				assert.Nil(t, result.Option)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ar := c.Arrange()
			option, err := dekopin.NewCmdOption(ar.ctx, ar.config, ar.cmd)
			c.Assert(t, ar, TestResult{
				Option: option,
				Err:    err,
			})
		})
	}
}
