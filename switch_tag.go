package dekopin

import (
	"context"
	"fmt"

	run "cloud.google.com/go/run/apiv2"
	runpb "cloud.google.com/go/run/apiv2/runpb"
	"github.com/spf13/cobra"
)

var stDeployCmd = &cobra.Command{
	Use:     "st-deploy",
	Short:   "Switch Tag Deploy(Assign a Revision tag to a Cloud Run revision)",
	PreRunE: stDeployPreRun,
	RunE:    SwitchTagDeployCommand,
}

func SwitchTagDeployCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	opt, err := GetCmdOption(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cmdOption: %w", err)
	}

	tag, err := getTagByFlag(cmd)
	if err != nil {
		return fmt.Errorf("failed to get tag flag: %w", err)
	}

	rt, err := createRevisionTagName(ctx, tag)
	if err != nil {
		return fmt.Errorf("failed to get tag name: %w", err)
	}

	gc, err := GetGcloudCommand(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to get gcloud command: %w", err)
	}

	return switchTagDeploy(cmd.Context(), gc, opt, rt)
}

func switchTagDeploy(ctx context.Context, gc GcloudCommand, opt *CmdOption, tag string) error {
	// タグが存在するか確認する
	// gcp-goを使って確認する
	client, err := run.NewServicesClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create run client: %w", err)
	}
	defer client.Close()

	fullServiceName := fmt.Sprintf(SERVICE_FULL_NAME_FORMAT, opt.Project, opt.Region, opt.Service)
	services, err := client.GetService(ctx, &runpb.GetServiceRequest{
		Name: fullServiceName,
	})
	if err != nil {
		return fmt.Errorf("failed to get service: %w", err)
	}

	for _, service := range services.TrafficStatuses {
		fmt.Println("service: ", service)
	}

	return nil
}
