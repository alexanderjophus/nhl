package hockeystats

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	hostURL = "https://statsapi.web.nhl.com/api/v1/"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

type statsResponse struct {
	Stats []stats `json:"stats"`
}

type stats struct {
	Splits []split `json:"splits"`
}

type split struct {
	Date string `json:"date"`
	Stat stat   `json:"stat"`
}

type stat struct {
	Points            int `json:"points"`
	Goals             int `json:"goals"`
	Assists           int `json:"assists"`
	PlusMinus         int `json:"plusMinus"`
	PowerPlayGoals    int `json:"powerPlayGoals"`
	PowerPlayPoints   int `json:"powerPlayPoints"`
	ShortHandedGoals  int `json:"shortHandedGoals"`
	ShortHandedPoints int `json:"shortHandedPoints"`
	Shots             int `json:"shots"`
	Pim               int `json:"pim"`
	GameWinningGoals  int `json:"gameWinningGoals"`
	OverTimeGoals     int `json:"overTimeGoals"`
}

type rosterResponse struct {
	RosterEntry []rosterEntry `json:"roster"`
}

type rosterEntry struct {
	Person       Person `json:"person"`
	JerseyNumber string `json:"jerseyNumber"`
}

type peopleResponse struct {
	People []Person `json:"people"`
}

type Person struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func GetPlayer(playerID string) (*Person, error) {
	playerURL := fmt.Sprintf("people/%s", playerID)
	resp := peopleResponse{}
	err := getJSON(hostURL+playerURL, &resp)
	if err != nil {
		return nil, fmt.Errorf("could not verify player %s: %s", playerID, err)
	}
	return &resp.People[0], nil
}

func (p *Person) GameLogStats() (statsResponse, error) {
	playerStats := fmt.Sprintf("people/%d/stats?stats=gameLog", p.ID)
	resp := statsResponse{}
	if err := getJSON(hostURL+playerStats, &resp); err != nil {
		return statsResponse{}, err
	}
	return resp, nil
}

func GetTeamPlayers(team int) (people []string, err error) {
	teamRoster := fmt.Sprintf("teams/%d/roster", team)
	resp := new(rosterResponse)
	err = getJSON(hostURL+teamRoster, resp)
	if err != nil {
		return
	}

	for _, player := range resp.RosterEntry {
		people = append(people, strconv.Itoa(player.Person.ID))
	}
	return
}

func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", r.Status)
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
