package pivotal

import (
	"time"
)

type Label struct {
	ID        int       `json:"id,omitempty"`
	ProjectID int       `json:"project_id,omitempty"`
	Kind      string    `json:"kind,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// func (c *Client) CreateOrGetLabel(project int, name string) (Label, error) {
//
// }
//
// func (c *Client) GetLabel(project int, name string) (Label, error) {
// 	var l Label
// 	err := c.get(fmt.Sprintf("/projects/{project_id}/labels", project), &l)
// 	return l, err
// }
