package dekopin_test

import (
	"testing"

	"github.com/iwashi623/dekopin"
	"github.com/samber/lo"
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
