package cmd

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func toString(list []*cobra.Command) []string {
	newList := make([]string, len(list))
	for i, v := range list {
		newList[i] = v.Use
	}
	return newList
}

func suggestionToString(list []prompt.Suggest) []string {
	newList := make([]string, len(list))
	for i, v := range list {
		newList[i] = v.Text
	}
	return newList
}

func branchToString(list []Branch) []string {
	newList := make([]string, len(list))
	for i, v := range list {
		newList[i] = v.Name
	}
	return newList
}

func TestCommonCommandsList(t *testing.T) {
	expects := []string{"pull --rebase origin master", "commit -a --amend --no-edit"}
	reality := toString(CommonCommandsList())
	for _, e := range expects {
		assert.Contains(t, reality, e)
	}
}

func TestBranchList(t *testing.T) {
	expects := []string{"master"}
	notexpects := []string{"origin/master", "origin/HEAD"}
	reality := branchToString(BranchList())
	for _, e := range expects {
		assert.Contains(t, reality, e)
	}
	for _, ne := range notexpects {
		assert.NotContains(t, reality, ne)
	}
}

func TestToStructuredBranchList(t *testing.T) {
	expects :=
		[]struct {
			raw             string
			expectedFirstBranchName string
		}{
			{
				`'2020-10-06; John Doe; bf84c09; origin/other-branch; (3 days ago)'
'2020-10-06; John Doe; bf84c09; origin/master; (3 days ago)'`,
				"origin/other-branch",
			},
			{
				`warning: ignoring broken ref refs/remotes/origin/HEAD
'2020-10-02; John Doe; e5cffc5; origin/release-v2.11.0; (7 days ago)'
'2020-10-01; John Doe; 2f41d5e; origin/feature_FD-5860; (8 days ago)'`,
				"origin/release-v2.11.0",
			},
		}
	for _, e := range expects {
		fmt.Println(e.expectedFirstBranchName)
		list := toStructuredBranchList(e.raw)
		assert.Greaterf(t, len(list), 0, e.expectedFirstBranchName)
		reality := list[0].Name
		assert.Equal(t, reality, e.expectedFirstBranchName)
	}
}


// Tests AllBitAndGitSubCommands has common commands, git sub commands, git aliases, git-extras and bit commands
func TestAllBitAndGitSubCommands(t *testing.T) {
	expects := []string{"pull --rebase origin master", "commit -a --amend --no-edit", "add", "push", "fetch", "pull", "co", "lg", "release", "info", "save", "sync"}
	reality := toString(AllBitAndGitSubCommands(ShellCmd))
	for _, e := range expects {
		assert.Contains(t, reality, e)
	}
}

func TestParseManPage(t *testing.T) {
	reality := parseManPage("rebase")
	assert.NotContains(t, reality, "GIT-REBASE(1)")
}

func TestFlagSuggestionsForCommand(t *testing.T) {
	// fixme add support for all git sub commands
	expects :=
		[]struct {
			cmd             string
			expectedOptions []string
			expectedFlags   []string
		}{
			{
				"rebase",
				[]string{"-i"},
				[]string{"--continue", "--abort", "--merge"},
			},
			{
				"push",
				[]string{"-f"},
				[]string{"--force", "--dry-run", "--porcelain", "--delete", "--tags"},
			},
			{
				"pull",
				[]string{"-q"},
				[]string{"--ff-only", "--no-ff", "--no-edit"},
			},
		}
	for _, e := range expects {
		realityFlags := suggestionToString(FlagSuggestionsForCommand(e.cmd, "--"))
		for _, ee := range e.expectedFlags {
			assert.Contains(t, realityFlags, ee)
		}
		realityOptions := suggestionToString(FlagSuggestionsForCommand(e.cmd, "-"))
		for _, ee := range e.expectedOptions {
			assert.Contains(t, realityOptions, ee)
		}

	}
}

func BenchmarkAllBitAndGitSubCommands(b *testing.B) {
	for n := 0; n < b.N; n++ {
		AllBitAndGitSubCommands(ShellCmd)
	}
}
