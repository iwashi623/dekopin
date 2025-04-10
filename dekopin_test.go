package dekopin_test

import (
	"context"
	"testing"

	"github.com/iwashi623/dekopin"
	"github.com/stretchr/testify/assert"
)

type TestCase[TArgs, TArrangeResult, TActResult any] struct {
	Args    *TArgs
	Arrange func() TArrangeResult
	Assert  func(t *testing.T, assertArgs TArrangeResult, result TActResult)
}

func TestCreateRevisionTagName(t *testing.T) {
	type TestResult struct {
		Tag string
		Err error
	}

	type ArrangeResult struct {
		ctx context.Context
		tag string
	}

	cases := map[string]TestCase[any, ArrangeResult, TestResult]{
		"success_if_tag_is_not_empty_returns_the_tag_as_is": {
			Arrange: func() ArrangeResult {
				return ArrangeResult{
					ctx: context.Background(),
					tag: "test",
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.NoError(t, result.Err)
				assert.Equal(t, assertArgs.tag, result.Tag)
			},
		},
		"success_if_tag_is_empty_and_runner_is_not_local_returns_string_starting_with_tag": {
			Arrange: func() ArrangeResult {
				opt := dekopin.CmdOption{
					Runner: dekopin.RUNNER_GITHUB_ACTIONS,
				}

				ctx := dekopin.SetCmdOption(context.Background(), &opt)
				return ArrangeResult{
					ctx: ctx,
					tag: "",
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.NoError(t, result.Err)
				assert.Regexp(t, "^tag-.*", result.Tag)
			},
		},
		"error_if_tag_is_empty_and_runner_is_local_returns_error": {
			Arrange: func() ArrangeResult {
				opt := dekopin.CmdOption{
					Runner: dekopin.RUNNER_LOCAL,
				}

				ctx := dekopin.SetCmdOption(context.Background(), &opt)
				return ArrangeResult{
					ctx: ctx,
					tag: "",
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
			result, err := dekopin.CreateRevisionTagName(ar.ctx, ar.tag)
			c.Assert(t, ar, TestResult{
				Tag: result,
				Err: err,
			})
		})
	}
}

func TestGetCommitHash(t *testing.T) {
	type TestResult struct {
		CommitHash string
		Err        error
	}

	type ArrangeResult struct {
		ctx        context.Context
		env        map[string]string
		commitHash string
	}

	cases := map[string]TestCase[any, ArrangeResult, TestResult]{
		"success_github_actions_runner_with_valid_sha_longer_than_7_chars": {
			Arrange: func() ArrangeResult {
				hash := "abcdef1234567890"
				t.Setenv(dekopin.ENV_GITHUB_SHA, hash)

				opt := &dekopin.CmdOption{
					Runner: dekopin.RUNNER_GITHUB_ACTIONS,
				}
				ctx := dekopin.SetCmdOption(context.Background(), opt)
				return ArrangeResult{
					ctx:        ctx,
					env:        map[string]string{dekopin.ENV_GITHUB_SHA: hash},
					commitHash: hash[:dekopin.COMMIT_HASH_LENGTH],
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Equal(t, assertArgs.commitHash, result.CommitHash)
			},
		},
		"success_cloud_build_runner_with_valid_sha_longer_than_7_chars": {
			Arrange: func() ArrangeResult {
				hash := "1234567abcdef890"
				t.Setenv(dekopin.ENV_CLOUD_BUILD_SHA, hash)

				opt := &dekopin.CmdOption{
					Runner: dekopin.RUNNER_CLOUD_BUILD,
				}
				ctx := dekopin.SetCmdOption(context.Background(), opt)
				return ArrangeResult{
					ctx:        ctx,
					env:        map[string]string{dekopin.ENV_CLOUD_BUILD_SHA: hash},
					commitHash: hash[:dekopin.COMMIT_HASH_LENGTH],
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Equal(t, assertArgs.commitHash, result.CommitHash)
			},
		},
		"success_github_actions_runner_with_short_sha_less_than_or_equal_to_7_chars": {
			Arrange: func() ArrangeResult {
				hash := "abc123"
				t.Setenv(dekopin.ENV_GITHUB_SHA, hash)

				opt := &dekopin.CmdOption{
					Runner: dekopin.RUNNER_GITHUB_ACTIONS,
				}
				ctx := dekopin.SetCmdOption(context.Background(), opt)
				return ArrangeResult{
					ctx:        ctx,
					env:        map[string]string{dekopin.ENV_GITHUB_SHA: hash},
					commitHash: hash,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Equal(t, assertArgs.commitHash, result.CommitHash)
			},
		},
		"success_cloud_build_runner_with_short_sha_exactly_7_chars": {
			Arrange: func() ArrangeResult {
				hash := "1234567"
				t.Setenv(dekopin.ENV_CLOUD_BUILD_SHA, hash)

				opt := &dekopin.CmdOption{
					Runner: dekopin.RUNNER_CLOUD_BUILD,
				}
				ctx := dekopin.SetCmdOption(context.Background(), opt)
				return ArrangeResult{
					ctx:        ctx,
					env:        map[string]string{dekopin.ENV_CLOUD_BUILD_SHA: hash},
					commitHash: hash,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Equal(t, assertArgs.commitHash, result.CommitHash)
			},
		},
		"error_invalid_runner": {
			Arrange: func() ArrangeResult {
				opt := &dekopin.CmdOption{
					Runner: "invalid-runner",
				}
				ctx := dekopin.SetCmdOption(context.Background(), opt)
				return ArrangeResult{
					ctx: ctx,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Equal(t, "", result.CommitHash)
				assert.Error(t, result.Err)
			},
		},
		"error_environment_variable_not_set": {
			Arrange: func() ArrangeResult {
				opt := &dekopin.CmdOption{
					Runner: dekopin.RUNNER_GITHUB_ACTIONS,
				}
				ctx := dekopin.SetCmdOption(context.Background(), opt)
				return ArrangeResult{
					ctx: ctx,
					env: map[string]string{dekopin.ENV_GITHUB_SHA: ""},
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Equal(t, "", result.CommitHash)
				assert.Error(t, result.Err)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ar := c.Arrange()

			// Set environment variables
			for k, v := range ar.env {
				t.Setenv(k, v)
			}

			result, err := dekopin.GetCommitHash(ar.ctx)
			c.Assert(t, ar, TestResult{
				CommitHash: result,
				Err:        err,
			})

			// Clean up environment variables
			for k := range ar.env {
				t.Setenv(k, "")
			}
		})
	}
}

func TestGetRunnerRef(t *testing.T) {
	type TestResult struct {
		Ref string
		Err error
	}

	type ArrangeResult struct {
		ctx    context.Context
		runner string
		env    map[string]string
	}

	cases := map[string]TestCase[any, ArrangeResult, TestResult]{
		"success_github_actions_runner_with_valid_ref": {
			Arrange: func() ArrangeResult {
				opt := dekopin.CmdOption{
					Project: "test-project",
					Region:  "test-region",
					Service: "test-service",
					Runner:  dekopin.RUNNER_GITHUB_ACTIONS,
				}
				ctx := dekopin.SetCmdOption(context.Background(), &opt)

				return ArrangeResult{
					ctx:    ctx,
					runner: dekopin.RUNNER_GITHUB_ACTIONS,
					env: map[string]string{
						dekopin.ENV_GITHUB_REF: "refs/heads/main",
					},
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.NoError(t, result.Err)
				assert.Equal(t, "refs/heads/main", result.Ref)
			},
		},
		"success_cloud_build_runner_with_valid_ref": {
			Arrange: func() ArrangeResult {
				opt := dekopin.CmdOption{
					Project: "test-project",
					Region:  "test-region",
					Service: "test-service",
					Runner:  dekopin.RUNNER_CLOUD_BUILD,
				}
				ctx := dekopin.SetCmdOption(context.Background(), &opt)

				return ArrangeResult{
					ctx:    ctx,
					runner: dekopin.RUNNER_CLOUD_BUILD,
					env: map[string]string{
						dekopin.ENV_CLOUD_BUILD_REF: "main",
					},
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.NoError(t, result.Err)
				assert.Equal(t, "main", result.Ref)
			},
		},
		"error_local_runner_returns_error": {
			Arrange: func() ArrangeResult {
				opt := dekopin.CmdOption{
					Project: "test-project",
					Region:  "test-region",
					Service: "test-service",
					Runner:  dekopin.RUNNER_LOCAL,
				}
				ctx := dekopin.SetCmdOption(context.Background(), &opt)

				return ArrangeResult{
					ctx:    ctx,
					runner: dekopin.RUNNER_LOCAL,
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

			// Clean up environment variables
			for k := range ar.env {
				t.Setenv(k, "")
			}

			// Set environment variables
			for k, v := range ar.env {
				t.Setenv(k, v)
			}

			ref, err := dekopin.GetRunnerRef(ar.ctx)
			c.Assert(t, ar, TestResult{
				Ref: ref,
				Err: err,
			})
		})
	}
}

func TestValidateTag(t *testing.T) {
	type TestResult struct {
		Err error
	}

	type ArrangeResult struct {
		tag string
	}

	cases := map[string]TestCase[any, ArrangeResult, TestResult]{
		"success_tag_is_not_empty": {
			Arrange: func() ArrangeResult {
				return ArrangeResult{
					tag: "test",
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.NoError(t, result.Err)
			},
		},
		"success_empty_tag_does_not_cause_error": {
			Arrange: func() ArrangeResult {
				return ArrangeResult{
					tag: "",
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.NoError(t, result.Err)
			},
		},
		"success_tag_contains_only_lowercase_alphanumeric_and_hyphens": {
			Arrange: func() ArrangeResult {
				return ArrangeResult{
					tag: "test-tag12345",
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.NoError(t, result.Err)
			},
		},
		"error_tag_contains_characters_other_than_lowercase_alphanumeric_and_hyphens": {
			Arrange: func() ArrangeResult {
				return ArrangeResult{
					tag: "test.tag12345",
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
			err := dekopin.ValidateTag(ar.tag)
			c.Assert(t, ar, TestResult{
				Err: err,
			})
		})
	}
}
