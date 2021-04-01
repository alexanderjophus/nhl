package nhl

type franchisesResponse struct {
	Copyright  string      `json:"copyright"`
	Franchises []Franchise `json:"franchises"`
}

type Franchise struct {
	Franchiseid      int    `json:"franchiseId"`
	Firstseasonid    int    `json:"firstSeasonId"`
	Mostrecentteamid int    `json:"mostRecentTeamId"`
	Teamname         string `json:"teamName"`
	Locationname     string `json:"locationName"`
	Link             string `json:"link"`
	Lastseasonid     int    `json:"lastSeasonId,omitempty"`
}

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
