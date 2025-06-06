package scratch

import (
	"testing"

	"github.com/michaelmdeng/mdcli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestNewAction_ArgumentCount(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		expectedErr string
	}{
		{
			name:        "NoArguments",
			args:        []string{},
			expectedErr: "exactly one argument <name> must be provided",
		},
		{
			name:        "TooManyArguments",
			args:        []string{"arg1", "arg2"},
			expectedErr: "exactly one argument <name> must be provided",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := cli.NewApp()
			app.ExitErrHandler = func(_ *cli.Context, err error) {
			}
			app.Commands = []*cli.Command{newCommand}

			err := app.Run(append([]string{"mdcli", "new"}, tc.args...))
			assert.Error(t, err)
			assert.ErrorContains(t, err, tc.expectedErr)
		})
	}
}

func TestNewAction_ScratchPath(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		cfg           func(*config.Config)
		expectedError string
	}{
		{
			name: "NeitherFlagNorConfig",
			args: []string{"name"},
			cfg: func(cfg *config.Config) {
				cfg.Scratch.ScratchPath = ""
			},
			expectedError: "configuration not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := cli.NewApp()
			app.ExitErrHandler = func(_ *cli.Context, err error) {
			}
			app.Commands = []*cli.Command{newCommand}
			err := app.Run(append([]string{"mdcli", "new"}, tc.args...))
			assert.Error(t, err)
			assert.ErrorContains(t, err, tc.expectedError)
		})
	}
}
