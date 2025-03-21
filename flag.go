package dekopin

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
