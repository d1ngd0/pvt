package pivotal

import (
	"fmt"
	"time"
)

type LabelBlocker struct {
	Matches string `json:"matches"`
	Blocks  string `json:"blocks"`
}

type Blocker struct {
	ID          int       `json:"id,omitempty"`
	StoryID     int       `json:"story_id,omitempty"`
	PersonID    int       `json:"person_id,omitempty"`
	Description string    `json:"description"`
	Resolved    bool      `json:"resolved,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Kind        string    `json:"kind,omitempty"`
}

func (c *Client) CreateLabelBlocker(project int, b LabelBlocker) error {
	matchingStories, err := c.GetFilteredStories(project, fmt.Sprintf("labels:\"%s\"", b.Matches))
	if err != nil {
		return err
	}

	blockers := make([]Blocker, len(matchingStories))
	for x, l := 0, len(matchingStories); x < l; x++ {
		blockers[x] = Blocker{Description: fmt.Sprintf("#%d", matchingStories[x].ID)}
	}

	blockedStories, err := c.GetFilteredStories(project, fmt.Sprintf("labels:\"%s\"", b.Blocks))
	if err != nil {
		return err
	}

	for x, l := 0, len(blockedStories); x < l; x++ {
		for y, l := 0, len(blockers); y < l; y++ {
			if err := c.CreateBlocker(project, blockedStories[x].ID, blockers[y]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) CreateBlocker(project, story int, blocker Blocker) error {
	return c.post(fmt.Sprintf("/projects/%d/stories/%d/blockers", project, story), blocker)
}
