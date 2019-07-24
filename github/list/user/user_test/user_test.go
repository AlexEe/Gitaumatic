package user_test

import (
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	createUser "omniactl/github/create/user"
	listUser "omniactl/github/list/user"
	"testing"
)

func TestCheckUsername(t *testing.T) {
	tests := []struct {
		data   string
		result string
	}{
		{"e666666", "e666666"},
		{"fake_user", ""},
	}
	for _, test := range tests {
		result := listUser.CheckUsername(test.data)
		assert.Equal(t, test.result, result)
	}
}

func TestGetGithubUser(t *testing.T) {
	var githubUser *github.User

	tests := []struct {
		data   string
		result *github.User
	}{
		{"e666666", githubUser},
	}

	for _, test := range tests {
		result := listUser.GetGithubUser(test.data)
		assert.IsType(t, test.result, result)
		assert.Equal(t, test.data, result.GetLogin())
	}
}

func TestGetOrgsForUser(t *testing.T) {
	var orgMap map[string]createUser.Org

	tests := []struct {
		username string
		org      string
		team     string
		result   map[string]createUser.Org
	}{
		{"e666666", "testing", "test_team", orgMap},
	}

	for _, test := range tests {
		result := listUser.GetOrgsForUser(test.username)
		assert.IsType(t, test.result, result)
		for _, v := range result {
			assert.Equal(t, test.org, v.Name)
			for _, vv := range v.Teams {
				assert.Equal(t, test.team, vv.Name)
			}
		}
	}
}

func TestGetTeamsForUser(t *testing.T) {
	var teamMap []createUser.Team

	tests := []struct {
		username string
		org      string
		team     string
		id       int64
		result   []createUser.Team
	}{
		{"e666666", "testing", "test_team", 90, teamMap},
	}

	for _, test := range tests {
		result := listUser.GetTeamsForUser(test.username, test.org)
		assert.IsType(t, test.result, result)
		for _, v := range result {
			assert.Equal(t, test.team, v.Name)
			assert.Equal(t, test.id, v.ID)
		}
	}
}
