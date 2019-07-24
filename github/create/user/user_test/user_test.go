package user_test

import (
	"context"
	createUser "omniactl/github/create/user"
	deleteUser "omniactl/github/delete/user"
	githubLogin "omniactl/login/github"
	"testing"
)

func TestCheckLogin(t *testing.T) {
	err := githubLogin.CheckGithubLogin()
	if err != nil {
		t.Error("Connection to Github failed.")
	}
}

func TestGetUsername(t *testing.T) {
	type test struct {
		data   string
		answer string
	}

	tests := []test{
		test{"e223344", "e223344"},
		test{"0000000", ""},
		test{"newUser", ""},
	}

	for _, v := range tests {
		x := createUser.GetUsername(v.data)
		if x != v.answer {
			t.Error("Expected", v.answer, "Got", x)
		}
	}
}

func TestCheckUsernameFormat(t *testing.T) {
	type test struct {
		data   string
		answer bool
	}

	tests := []test{
		test{"e223344", true},
		test{"0000000", false},
		test{"newUser", false},
	}

	for _, v := range tests {
		x := createUser.CheckUsernameFormat(v.data)
		if x != v.answer {
			t.Error("Expected", v.answer, "Got", x)
		}
	}
}

func TestCheckIfUserExists(t *testing.T) {
	type test struct {
		data   string
		answer bool
	}

	tests := []test{
		test{"e661018", true},
		test{"0000000", false},
		test{"newUser", false},
	}

	for _, v := range tests {
		x := createUser.CheckIfUserExists(v.data)
		if x != v.answer {
			t.Error("Expected", v.answer, "Got", x)
		}
	}
}

func TestGetEmail(t *testing.T) {
	type test struct {
		data   string
		answer string
	}

	tests := []test{
		test{"example@statestreet.com", "example@statestreet.com"},
		test{"notanaddress", ""},
		test{"1234", ""},
	}

	for _, v := range tests {
		x := createUser.GetEmail(v.data)
		if x != v.answer {
			t.Error("Expected", v.answer, "Got", x)
		}
	}
}

func TestCheckEmailFormat(t *testing.T) {
	type test struct {
		data   string
		answer bool
	}

	tests := []test{
		test{"example@statestreet.com", true},
		test{"notanaddress", false},
		test{"1234", false},
	}

	for _, v := range tests {
		x := createUser.CheckEmailFormat(v.data)
		if x != v.answer {
			t.Error("Expected", v.answer, "Got", x)
		}
	}
}

func TestCheckGetOrgs(t *testing.T) {
	type test struct {
		org   string
		role  string
		teams []string
	}

	tests := []test{
		test{
			org:   "testing",
			role:  "member",
			teams: []string{"test_team"},
		},
	}

	for _, v := range tests {
		x := createUser.GetOrgs(v.org, v.role, v.teams)
		for k, vv := range x {
			if k != v.org {
				t.Error("Expected", v.org, "Got", k)
			}
			if vv.Role != v.role {
				t.Error("Expected", v.role, "Got", vv.Role)
			}
			for _, vvv := range vv.Teams {
				for _, vvvv := range v.teams {
					if vvv.Name != vvvv {
						t.Error("Expected", vvvv, "Got", vvv.Name)
					}
				}
			}
		}
	}
}

func TestCheckIfOrgExists(t *testing.T) {
	type test struct {
		data   string
		answer bool
	}

	tests := []test{
		test{"testing", true},
		test{"falseOrg", false},
	}

	for _, v := range tests {
		x := createUser.CheckOrgExists(v.data)
		if x != v.answer {
			t.Error("Expected", v.answer, "Got", x)
		}
	}
}

func TestCheckIfRoleExists(t *testing.T) {
	type test struct {
		data   string
		answer bool
	}

	tests := []test{
		test{"member", true},
		test{"admin", true},
		test{"falseRole", false},
	}

	for _, v := range tests {
		x := createUser.CheckRoleExists(v.data)
		if x != v.answer {
			t.Error("Expected", v.answer, "Got", x)
		}
	}
}

func TestCheckTeamExists(t *testing.T) {
	type test struct {
		org    string
		team   string
		answer bool
	}

	tests := []test{
		test{"testing", "test_team", true},
		test{"testing", "false_team", false},
	}

	for _, v := range tests {
		x := createUser.CheckTeamExists(v.org, v.team)
		if x != v.answer {
			t.Error("Expected", v.answer, "Got", x)
		}
	}
}

func TestDeleteFlagOrgs(t *testing.T) {
	testOrg := map[string]createUser.Org{
		"testing": createUser.Org{
			Name: "testing",
			ID:   569,
			Role: "member",
			Teams: []createUser.Team{
				createUser.Team{
					Name: "test_team",
					ID:   90,
				},
			},
		},
	}

	x := createUser.DeleteFlagOrgs(testOrg)
	for k := range testOrg {
		_, ok := x[k]
		if ok {
			t.Error("Organisation set by flag was not deleted from list.")
		}
	}
}

func TestCreateOrgList(t *testing.T) {
	testOrg := map[string]createUser.Org{
		"testing": createUser.Org{
			Name: "testing",
			ID:   569,
			Role: "member",
			Teams: []createUser.Team{
				createUser.Team{
					Name: "test_team",
					ID:   90,
				},
			},
		},
	}

	x := createUser.CreateOrgList(testOrg)
	for _, v := range x {
		if v != "testing" {
			t.Error("Expected 'testing' Got", x)
		}
	}
}

func TestGetAllOrgs(t *testing.T) {
	testOrgs := map[string]int64{
		"MSF":      29,
		"testing":  569,
		"causeway": 32,
		"devtools": 684,
		"Security": 376,
	}

	x := createUser.GetAllOrgs()
	for k, v := range testOrgs {
		value, ok := x[k]
		if ok {
			if value.ID != v {
				t.Errorf("Retrieved organisation ID is not correct: Expected '%v' Got '%v'", v, value.ID)
			}
		} else {
			t.Errorf("'%v' Github organisation exists, but was not retrieved.", k)
		}
	}
}

func TestGetTeamsForOrg(t *testing.T) {
	org := "testing"
	teams := map[string]int64{
		"test_team":    90,
		"public_team3": 122,
		"team2":        108,
	}

	x := createUser.GetTeamsForOrg(org)
	for k, v := range teams {
		value, ok := x[k]
		if ok {
			if value.ID != v {
				t.Errorf("Retrieved team ID is not correct: Expected '%v' Got '%v'", v, value.ID)
			}
		} else {
			t.Errorf("'%v' Github team exists in organisation '%v', but was not retrieved.", k, org)
		}
	}
}

func TestCreateTeamList(t *testing.T) {
	Teams := map[string]createUser.Team{
		"test_team": createUser.Team{
			Name: "test_team",
			ID:   90,
		},
	}

	x := createUser.CreateTeamList(Teams)
	for k := range Teams {
		for _, v := range x {
			if v != k {
				t.Errorf("Expected '%v' Got '%v'", k, x)
			}
		}
	}
}

func TestCreateUser(t *testing.T) {
	type test struct {
		username string
		email    string
		answer   string
	}

	tests := []test{
		test{"e444444", "test4@statestreet.com", "e444444"},
	}

	for _, v := range tests {
		x, _ := createUser.CreateUser(v.username, v.email)
		if x != v.answer {
			t.Errorf("Expected '%v' Got '%v'", v.answer, x)
		}
	}
}

func TestAddUserToOrgs(t *testing.T) {
	username := "e444444"
	testOrg := map[string]createUser.Org{
		"testing": createUser.Org{
			Name: "testing",
			ID:   569,
			Role: "member",
			Teams: []createUser.Team{
				createUser.Team{
					Name: "test_team",
					ID:   90,
				},
			},
		},
	}

	createUser.AddUserToOrgs(username, testOrg)

	Client := githubLogin.CreateClient()
	for k, v := range testOrg {
		isMember, _, err := Client.Organizations.IsMember(context.Background(), k, username)
		if err != nil {
			deleteUser.DeleteFromGithub(username)
			t.Errorf("Checking member status of '%v' in organisation '%v' failed", username, k)
		}
		if isMember == false {
			deleteUser.DeleteFromGithub(username)
			t.Errorf("New user '%v' could not be added to organisation '%v'", username, k)
		}
		for _, vv := range v.Teams {
			isMember, _, err := Client.Teams.IsTeamMember(context.Background(), vv.ID, username)
			if err != nil {
				deleteUser.DeleteFromGithub(username)
				t.Errorf("Checking member status of '%v' in team '%v' failed", username, vv.Name)
			}
			if isMember == false {
				deleteUser.DeleteFromGithub(username)
				t.Errorf("New user '%v' could not be added to team '%v'", username, vv.Name)
			}
			deleteUser.DeleteFromGithub(username)
		}
	}
}
