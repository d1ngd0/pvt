package pivotal

import "time"

type Project struct {
	ID                           int       `json:"id"`
	Kind                         string    `json:"kind"`
	Name                         string    `json:"name"`
	Version                      int       `json:"version"`
	IterationLength              int       `json:"iteration_length"`
	WeekStartDay                 string    `json:"week_start_day"`
	PointScale                   string    `json:"point_scale"`
	PointScaleIsCustom           bool      `json:"point_scale_is_custom"`
	BugsAndChoresAreEstimatable  bool      `json:"bugs_and_chores_are_estimatable"`
	AutomaticPlanning            bool      `json:"automatic_planning"`
	EnableTasks                  bool      `json:"enable_tasks"`
	TimeZone                     TimeZone  `json:"time_zone"`
	VelocityAveragedOver         int       `json:"velocity_averaged_over"`
	NumberOfDoneIterationsToShow int       `json:"number_of_done_iterations_to_show"`
	HasGoogleDomain              bool      `json:"has_google_domain"`
	EnableIncomingEmails         bool      `json:"enable_incoming_emails"`
	InitialVelocity              int       `json:"initial_velocity"`
	Public                       bool      `json:"public"`
	AtomEnabled                  bool      `json:"atom_enabled"`
	ProjectType                  string    `json:"project_type"`
	StartTime                    time.Time `json:"start_time"`
	CreatedAt                    time.Time `json:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at"`
	AccountID                    int       `json:"account_id"`
	CurrentIterationNumber       int       `json:"current_iteration_number"`
	EnableFollowing              bool      `json:"enable_following"`
}

// GetProjects will get all the projects you currently have
func (c *Client) GetProjects() ([]Project, error) {
	var p []Project
	err := c.get("projects", &p)
	return p, err
}
