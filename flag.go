package dekopin

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type dekopinCommandKey struct{}

func SetDekopinCommand(ctx context.Context, cmd DekopinCommand) context.Context {
	return context.WithValue(ctx, dekopinCommandKey{}, cmd)
}

func GetDekopinCommand(ctx context.Context) (DekopinCommand, error) {
	cmd, ok := ctx.Value(dekopinCommandKey{}).(DekopinCommand)
	if !ok {
		return nil, fmt.Errorf("dekopin command not found")
	}
	return cmd, nil
}

type DekopinCommand interface {
	GetFileByFlag() (string, error)
	GetProjectByFlag() (string, error)
	GetRegionByFlag() (string, error)
	GetServiceByFlag() (string, error)
	GetRunnerByFlag() (string, error)
	GetImageByFlag() (string, error)
	GetTagByFlag() (string, error)
	GetRevisionByFlag() (string, error)
	GetCreateTagByFlag() (bool, error)
	GetRemoveTagsByFlag() (bool, error)
	GetUpdateTrafficByFlag() (bool, error)
}

type dekopinCommand struct {
	cobra.Command
}

var _ DekopinCommand = &dekopinCommand{}

func NewDekopinCommand(cmd *cobra.Command) DekopinCommand {
	return &dekopinCommand{
		Command: *cmd,
	}
}

func (c *dekopinCommand) GetFileByFlag() (string, error) {
	file, err := c.Flags().GetString("file")
	if err != nil {
		return "", fmt.Errorf("failed to get file flag: %w", err)
	}
	return file, nil
}

func (c *dekopinCommand) GetProjectByFlag() (string, error) {
	project, err := c.Flags().GetString("project")
	if err != nil {
		return "", fmt.Errorf("failed to get project flag: %w", err)
	}
	return project, nil
}

func (c *dekopinCommand) GetRegionByFlag() (string, error) {
	region, err := c.Flags().GetString("region")
	if err != nil {
		return "", fmt.Errorf("failed to get region flag: %w", err)
	}
	return region, nil
}

func (c *dekopinCommand) GetServiceByFlag() (string, error) {
	service, err := c.Flags().GetString("service")
	if err != nil {
		return "", fmt.Errorf("failed to get service flag: %w", err)
	}
	return service, nil
}

func (c *dekopinCommand) GetRunnerByFlag() (string, error) {
	runner, err := c.Flags().GetString("runner")
	if err != nil {
		return "", fmt.Errorf("failed to get runner flag: %w", err)
	}
	return runner, nil
}

func (c *dekopinCommand) GetImageByFlag() (string, error) {
	image, err := c.Flags().GetString("image")
	if err != nil {
		return "", fmt.Errorf("failed to get image flag: %w", err)
	}

	return image, nil
}

func (c *dekopinCommand) GetTagByFlag() (string, error) {
	tag, err := c.Flags().GetString("tag")
	if err != nil {
		return "", fmt.Errorf("failed to get tag flag: %w", err)
	}
	return tag, nil
}

func (c *dekopinCommand) GetRevisionByFlag() (string, error) {
	rv, err := c.Flags().GetString("revision")
	if err != nil {
		return "", fmt.Errorf("failed to get revision flag: %w", err)
	}

	return rv, nil
}

func (c *dekopinCommand) GetCreateTagByFlag() (bool, error) {
	createTag, err := c.Flags().GetBool("create-tag")
	if err != nil {
		return false, fmt.Errorf("failed to get create-tag flag: %w", err)
	}

	return createTag, nil
}

func (c *dekopinCommand) GetRemoveTagsByFlag() (bool, error) {
	removeTags, err := c.Flags().GetBool("remove-tags")
	if err != nil {
		return false, fmt.Errorf("failed to get remove-tags flag: %w", err)
	}
	return removeTags, nil
}

func (c *dekopinCommand) GetUpdateTrafficByFlag() (bool, error) {
	updateTraffic, err := c.Flags().GetBool("update-traffic")
	if err != nil {
		return false, fmt.Errorf("failed to get update-traffic flag: %w", err)
	}
	return updateTraffic, nil
}
