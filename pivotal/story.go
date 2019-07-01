package pivotal

import (
	"fmt"
	"net/url"
	"time"
)

type Story struct {
	Kind          string     `json:"kind,omitempty"`
	ID            int        `json:"id,omitempty"`
	CurrentState  string     `json:"current_state,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	Description   string     `json:"description"`
	StoryType     string     `json:"story_type,omitempty"`
	Name          string     `json:"name"`
	RequestedByID int        `json:"requested_by_id,omitempty"`
	URL           string     `json:"url,omitempty"`
	OwnerIds      []int      `json:"owner_ids,omitempty"`
	Labels        []Label    `json:"labels,omitempty"`
	Blockers      []Blocker  `json:"blockers,omitempty"`
}

// GetFilteredStories will get all the stories matching the project and story
func (c *Client) GetFilteredStories(project int, filter string) ([]Story, error) {
	var s []Story
	err := c.get(fmt.Sprintf("projects/%d/stories?filter=%s&fields=%s", project, url.QueryEscape(filter), getFields(Story{})), &s)
	return s, err
}

func (c *Client) CreateStory(project int, story Story) error {
	return c.post(fmt.Sprintf("/projects/%d/stories", project), story)
}

func (c *Client) UpdateStory(project int, story Story) error {
	return c.put(fmt.Sprintf("/projects/%d/stories/%d", project, story.ID), story)
}

func (c *Client) GetStory(project int, storyId int) (Story, error) {
	var story Story
	err := c.get(fmt.Sprintf("/projects/%d/stories/%d?fields=%s", project, storyId, getFields(Story{})), &story)
	return story, err
}

func (c *Client) DeleteStory(project int, storyId int) error {
	return c.delete(fmt.Sprintf("/projects/%d/stories/%d", project, storyId))
}
