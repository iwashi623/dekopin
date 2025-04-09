package dekopin

import (
	"context"
	"fmt"
	"slices"

	"github.com/spf13/cobra"
)

type CmdOption struct {
	Project string
	Region  string
	Service string
	Runner  string
}

type cmdOptionKey struct{}

func SetCmdOption(ctx context.Context, cmdOption *CmdOption) context.Context {
	return context.WithValue(ctx, cmdOptionKey{}, cmdOption)
}

func GetCmdOption(ctx context.Context) (*CmdOption, error) {
	cmdOption, ok := ctx.Value(cmdOptionKey{}).(*CmdOption)
	if !ok {
		return nil, fmt.Errorf("cmdOption not found")
	}
	return cmdOption, nil
}

func NewCmdOption(ctx context.Context, config *DekopinConfig, cmd *cobra.Command) (*CmdOption, error) {
	dekopinCmd, err := GetDekopinCommand(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dekopin command: %w", err)
	}

	project, err := dekopinCmd.GetProjectByFlag()
	if err != nil {
		return nil, fmt.Errorf("failed to get project flag: %w", err)
	}
	if project == "" {
		project = config.Project
	}

	region, err := dekopinCmd.GetRegionByFlag()
	if err != nil {
		return nil, fmt.Errorf("failed to get region flag: %w", err)
	}
	if region == "" {
		region = config.Region
	}

	service, err := dekopinCmd.GetServiceByFlag()
	if err != nil {
		return nil, fmt.Errorf("failed to get service flag: %w", err)
	}
	if service == "" {
		service = config.Service
	}

	runner, err := dekopinCmd.GetRunnerByFlag()
	if err != nil {
		return nil, fmt.Errorf("failed to get runner flag: %w", err)
	}
	if runner == "" {
		runner = config.Runner
	}

	option := &CmdOption{
		Project: project,
		Region:  region,
		Service: service,
		Runner:  runner,
	}

	if err := option.Validate(); err != nil {
		return nil, err
	}

	return option, nil
}

func (c *CmdOption) Validate() error {
	if c.Project == "" || c.Region == "" || c.Service == "" || c.Runner == "" {
		return fmt.Errorf("project, region, service, and runner are required")
	}

	if !slices.Contains(ValidRunners, c.Runner) {
		return fmt.Errorf("invalid runner type. Valid values: github-actions, cloud-build, local")
	}

	return nil
}
