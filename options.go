package dekopin

import (
	"context"
	"fmt"

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

func NewCmdOption(config *DekopinConfig, cmd *cobra.Command) (*CmdOption, error) {
	project, err := getProjectByFlag(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get project flag: %w", err)
	}
	if project == "" {
		project = config.Project
	}

	region, err := getRegionByFlag(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get region flag: %w", err)
	}
	if region == "" {
		region = config.Region
	}

	service, err := getServiceByFlag(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get service flag: %w", err)
	}
	if service == "" {
		service = config.Service
	}

	runner, err := getRunnerByFlag(cmd)
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

	if err := option.validate(); err != nil {
		return nil, err
	}

	return option, nil
}

func (c *CmdOption) validate() error {
	if c.Project == "" || c.Region == "" || c.Service == "" || c.Runner == "" {
		return fmt.Errorf("project, region, service, and runner are required")
	}

	if !validRunners[c.Runner] {
		return fmt.Errorf("invalid runner type. Valid values: github-actions, cloud-build, local")
	}

	return nil
}
