package nhl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	hostURL = "https://statsapi.web.nhl.com/api/v1"
)

type Client struct {
	c *http.Client
}

type Option func(Client)

func NewClient(opts ...Option) Client {
	defC := http.DefaultClient
	defC.Timeout = 10 * time.Second

	c := Client{c: defC}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) GetFranchises() ([]Franchise, error) {
	resp := franchisesResponse{}
	err := c.getJSON(fmt.Sprintf("%s/franchises", hostURL), &resp)
	if err != nil {
		return []Franchise{}, fmt.Errorf("getting franchises: %w", err)
	}
	return resp.Franchises, nil
}

func (c *Client) GetFranchise(franchiseID int) (Franchise, error) {
	resp := franchisesResponse{}
	err := c.getJSON(fmt.Sprintf("%s/franchises/%d", hostURL, franchiseID), &resp)
	if err != nil {
		return Franchise{}, fmt.Errorf("getting franchises: %w", err)
	}

	if len(resp.Franchises) != 1 {
		return Franchise{}, fmt.Errorf("cannot find franchise")
	}

	return resp.Franchises[0], nil
}

func (c *Client) GetPlayer(playerID string) (*Person, error) {
	resp := peopleResponse{}
	err := c.getJSON(fmt.Sprintf("%s/people/%s", hostURL, playerID), &resp)
	if err != nil {
		return nil, fmt.Errorf("could not verify player %s: %s", playerID, err)
	}
	return &resp.People[0], nil
}

func (c *Client) GameLogStats(p *Person) (statsResponse, error) {
	resp := statsResponse{}
	playerStatsURL := fmt.Sprintf("%s/people/%d/stats?stats=gameLog", hostURL, p.ID)
	if err := c.getJSON(playerStatsURL, &resp); err != nil {
		return statsResponse{}, err
	}
	return resp, nil
}

func (c *Client) GetTeamPlayerIDs(team int) (people []string, err error) {
	teamRosterURL := fmt.Sprintf("%s/teams/%d/roster", hostURL, team)
	resp := rosterResponse{}
	err = c.getJSON(teamRosterURL, &resp)
	if err != nil {
		return
	}

	for _, player := range resp.RosterEntry {
		people = append(people, strconv.Itoa(player.Person.ID))
	}
	return
}

func (c *Client) getJSON(url string, target interface{}) error {
	r, err := c.c.Get(url)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %s", r.Status)
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
