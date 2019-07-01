package pivotal

import (
	"fmt"
	"net/url"
	"time"
)

type Epic struct {
	ID          int        `json:"id,omitempty"`
	Kind        string     `json:"kind,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	ProjectID   int        `json:"project_id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	URL         string     `json:"url,omitempty"`
	Label       *Label     `json:"label,omitempty"`
}

func (c *Client) CreateEpic(project int, epic Epic) error {
	return c.post(fmt.Sprintf("/projects/%d/epics", project), epic)
}

func (c *Client) UpdateEpic(project int, epic Epic) error {
	return c.put(fmt.Sprintf("/projects/%d/epics/%d", project, epic.ID), epic)
}

func (c *Client) GetEpic(project int, epicId int) (Epic, error) {
	var epic Epic
	err := c.get(fmt.Sprintf("/projects/%d/epics/%d?fields=%s", project, epicId, getFields(Epic{})), &epic)
	return epic, err
}

// GetProjects will get all the projects you currently have
func (c *Client) GetFilteredEpics(project int, filter string) ([]Epic, error) {
	var e []Epic
	err := c.get(fmt.Sprintf("projects/%d/epics?filter=%s", project, url.QueryEscape(filter)), &e)
	return e, err
}

func (c *Client) DeleteEpic(project int, id int) error {
	return c.delete(fmt.Sprintf("/projects/%d/epics/%d", project, id))
}
