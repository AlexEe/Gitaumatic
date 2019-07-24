package jira

// CreateProject creates a JIRA project with the provided name
func CreateProject(jiraProjectName string) {
}

// CreateRole creates a new role in JIRA with the provided name
func CreateRole(roleName string) {

}

// AddUser adds a user to JIRA
func AddUser(username string, name string, email string, jiraProjectName string) {
}

// ListProject lists projects in JIRA matching the provided glob
func ListProject(projectNameGlob string) {

}

// ListRole lists the roles in JIRA matching the provided glob.
func ListRole(roleGlob string) {

}

// ListUser lists the users with similar names to the provided user name
func ListUser(username string) {

}

// SuspendUser marks the specified user as suspending in JIRA with the reasons provided.
func SuspendUser(username string, reasons string) {

}
