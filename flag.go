package dekopin

import (
	"fmt"

	"github.com/spf13/cobra"
)

func getFileByFlag(cmd *cobra.Command) (string, error) {
	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return "", fmt.Errorf("failed to get file flag: %w", err)
	}
	return file, nil
}

func getProjectByFlag(cmd *cobra.Command) (string, error) {
	project, err := cmd.Flags().GetString("project")
	if err != nil {
		return "", fmt.Errorf("failed to get project flag: %w", err)
	}
	return project, nil
}

func getRegionByFlag(cmd *cobra.Command) (string, error) {
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		return "", fmt.Errorf("failed to get region flag: %w", err)
	}
	return region, nil
}

func getServiceByFlag(cmd *cobra.Command) (string, error) {
	service, err := cmd.Flags().GetString("service")
	if err != nil {
		return "", fmt.Errorf("failed to get service flag: %w", err)
	}
	return service, nil
}

func getRunnerByFlag(cmd *cobra.Command) (string, error) {
	runner, err := cmd.Flags().GetString("runner")
	if err != nil {
		return "", fmt.Errorf("failed to get runner flag: %w", err)
	}
	return runner, nil
}

func getImageByFlag(cmd *cobra.Command) (string, error) {
	image, err := cmd.Flags().GetString("image")
	if err != nil {
		return "", fmt.Errorf("failed to get image flag: %w", err)
	}

	return image, nil
}

func getTagByFlag(cmd *cobra.Command) (string, error) {
	tag, err := cmd.Flags().GetString("tag")
	if err != nil {
		return "", fmt.Errorf("failed to get tag flag: %w", err)
	}
	return tag, nil
}

func getRevisionByFlag(cmd *cobra.Command) (string, error) {
	rv, err := cmd.Flags().GetString("revision")
	if err != nil {
		return "", fmt.Errorf("failed to get revision flag: %w", err)
	}

	return rv, nil
}

func getCreateTagByFlag(cmd *cobra.Command) (bool, error) {
	createTag, err := cmd.Flags().GetBool("create-tag")
	if err != nil {
		return false, fmt.Errorf("failed to get create-tag flag: %w", err)
	}

	return createTag, nil
}

func getRemoveTagByFlag(cmd *cobra.Command) (bool, error) {
	removeTag, err := cmd.Flags().GetBool("remove-tags")
	if err != nil {
		return false, fmt.Errorf("failed to get remove-tags flag: %w", err)
	}
	return removeTag, nil
}
