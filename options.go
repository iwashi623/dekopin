package dekopin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type cmdOptionKey struct{}

type CmdOption struct {
	Project string
	Region  string
	Service string
	Runner  string
}

func NewCmdOption(config *DekopinConfig, cmd *cobra.Command) (*CmdOption, error) {
	project, err := cmd.Flags().GetString("project")
	if err != nil {
		return nil, fmt.Errorf("failed to get project flag: %w", err)
	}
	if project == "" {
		project = config.Project
	}

	region, err := cmd.Flags().GetString("region")
	if err != nil {
		return nil, fmt.Errorf("failed to get region flag: %w", err)
	}
	if region == "" {
		region = config.Region
	}

	service, err := cmd.Flags().GetString("service")
	if err != nil {
		return nil, fmt.Errorf("failed to get service flag: %w", err)
	}
	if service == "" {
		service = config.Service
	}

	runner, err := cmd.Flags().GetString("runner")
	if err != nil {
		return nil, fmt.Errorf("failed to get runner flag: %w", err)
	}
	if runner == "" {
		runner = config.Runner
	}

	if project == "" || region == "" || service == "" || runner == "" {
		return nil, fmt.Errorf("project, region, service, and runner are required")
	}

	return &CmdOption{
		Project: project,
		Region:  region,
		Service: service,
		Runner:  runner,
	}, nil
}

func GetCmdOption(ctx context.Context) (*CmdOption, error) {
	cmdOption, ok := ctx.Value(cmdOptionKey{}).(*CmdOption)
	if !ok {
		return nil, fmt.Errorf("cmdOption not found")
	}
	return cmdOption, nil
}
