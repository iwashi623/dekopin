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
		"正常系_引数のtagが空文字でなければ、tagをそのまま返す": {
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
		"正常系_引数のtagが空文字で、runnerがlocalでなければ、`tag-`から始まる文字列を返す": {
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
		"異常系_引数のtagが空文字で、runnerがlocalであれば、エラーを返す": {
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
	}

	type ArrangeResult struct {
		runner     string
		env        map[string]string
		commitHash string
	}

	cases := map[string]TestCase[any, ArrangeResult, TestResult]{
		"正常系_GitHub Actions Runnerで有効なSHA(7文字より長い)": {
			Arrange: func() ArrangeResult {
				hash := "abcdef1234567890"
				t.Setenv(dekopin.ENV_GITHUB_SHA, hash)
				return ArrangeResult{

					runner:     dekopin.RUNNER_GITHUB_ACTIONS,
					env:        map[string]string{dekopin.ENV_GITHUB_SHA: hash},
					commitHash: hash[:dekopin.COMMIT_HASH_LENGTH],
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Equal(t, assertArgs.commitHash, result.CommitHash)
			},
		},
		"正常系_Cloud Build Runnerで有効なSHA(7文字より長い)": {
			Arrange: func() ArrangeResult {
				hash := "1234567abcdef890"
				t.Setenv(dekopin.ENV_CLOUD_BUILD_SHA, hash)
				return ArrangeResult{
					runner:     dekopin.RUNNER_CLOUD_BUILD,
					env:        map[string]string{dekopin.ENV_CLOUD_BUILD_SHA: hash},
					commitHash: hash[:dekopin.COMMIT_HASH_LENGTH],
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Equal(t, assertArgs.commitHash, result.CommitHash)
			},
		},
		"正常系_GitHub Actions Runnerで短いSHA(7文字以下)": {
			Arrange: func() ArrangeResult {
				hash := "abc123"
				t.Setenv(dekopin.ENV_GITHUB_SHA, hash)
				return ArrangeResult{
					runner:     dekopin.RUNNER_GITHUB_ACTIONS,
					env:        map[string]string{dekopin.ENV_GITHUB_SHA: hash},
					commitHash: hash,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Equal(t, assertArgs.commitHash, result.CommitHash)
			},
		},
		"正常系_Cloud Build Runnerで短いSHA(7文字)": {
			Arrange: func() ArrangeResult {
				hash := "1234567"
				t.Setenv(dekopin.ENV_CLOUD_BUILD_SHA, hash)
				return ArrangeResult{
					runner:     dekopin.RUNNER_CLOUD_BUILD,
					env:        map[string]string{dekopin.ENV_CLOUD_BUILD_SHA: hash},
					commitHash: hash,
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Equal(t, assertArgs.commitHash, result.CommitHash)
			},
		},
		"異常系_無効なRunner": {
			Arrange: func() ArrangeResult {
				return ArrangeResult{
					runner: "invalid-runner",
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Equal(t, "", result.CommitHash)
			},
		},
		"異常系_環境変数が設定されていない場合": {
			Arrange: func() ArrangeResult {
				// 環境変数は設定しない
				return ArrangeResult{
					runner: dekopin.RUNNER_GITHUB_ACTIONS,
					env:    map[string]string{dekopin.ENV_GITHUB_SHA: ""},
				}
			},
			Assert: func(t *testing.T, assertArgs ArrangeResult, result TestResult) {
				assert.Equal(t, "", result.CommitHash)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ar := c.Arrange()

			// 環境変数を設定
			for k, v := range ar.env {
				t.Setenv(k, v)
			}

			result := dekopin.GetCommitHash(ar.runner)
			c.Assert(t, ar, TestResult{
				CommitHash: result,
			})

			// 環境変数を削除
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
		"正常系_GitHub Actions Runnerで有効なRef": {
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
		"正常系_Cloud Build Runnerで有効なRef": {
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
		"異常系_Local Runnerではエラーを返す": {
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

			// 環境変数を削除
			for k := range ar.env {
				t.Setenv(k, "")
			}

			// 環境変数を設定
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
