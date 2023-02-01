package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistories(t *testing.T) {
	r, config, err := CreateConfig("TestHistories")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		args []string
		test func(*testing.T, []byte, []byte)
	}{
		{
			name: "without flags",
			args: []string{"list-logins", "--format", "csv", "use-pager=false"},
			test: func(t *testing.T, stdOut []byte, stdErr []byte) {
				assert.Equal(
					t,
					`id,path,when,destination,display,ip,source,targets,user_id,username,action,failure_type,interface
1,,2022-01-28T19:41:43-05:00,,,,,,1234,,,,
1,,2023-01-28T19:41:43-05:00,,,,,,1234,,,,
`,
					string(stdOut),
				)
				assert.Empty(t, string(stdErr))
			},
		},
		{
			name: "with flags",
			args: []string{"list-logins", "--format", "csv", "--start-at", "2023-01-28T19:41:43-05:00", "--end-at", "2023-01-30T07:43:37-05:00", "use-pager=false"},
			test: func(t *testing.T, stdOut []byte, stdErr []byte) {
				assert.Equal(
					t,
					`id,path,when,destination,display,ip,source,targets,user_id,username,action,failure_type,interface
1,,2023-01-28T19:41:43-05:00,,,,,,1234,,,,
`,
					string(stdOut),
				)
				assert.Empty(t, string(stdErr))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.args)
			stdOut, stdErr := callCmd(Histories(), config, tt.args)

			tt.test(t, stdOut, stdErr)
		})
	}
	r.Stop()
}
