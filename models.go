package nhl

type franchisesResponse struct {
	Copyright  string      `json:"copyright"`
	Franchises []Franchise `json:"franchises"`
}

type Franchise struct {
	Franchiseid      int    `json:"franchiseId"`
	Firstseasonid    int    `json:"firstSeasonId,omitempty"`
	Mostrecentteamid int    `json:"mostRecentTeamId,omitempty"`
	Teamname         string `json:"teamName"`
	Locationname     string `json:"locationName,omitempty"`
	Link             string `json:"link"`
	Lastseasonid     int    `json:"lastSeasonId,omitempty"`
}

type teamsResponse struct {
	Copyright string `json:"copyright"`
	Teams     []Team `json:"teams"`
}
type Timezone struct {
	ID     string `json:"id"`
	Offset int    `json:"offset"`
	Tz     string `json:"tz"`
}
type Division struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}
type Conference struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}
type Venue struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Link     string   `json:"link"`
	City     string   `json:"city"`
	Timezone Timezone `json:"timeZone"`
}
type Team struct {
	ID              int        `json:"id"`
	Name            string     `json:"name"`
	Link            string     `json:"link"`
	Venue           Venue      `json:"venue,omitempty"`
	Abbreviation    string     `json:"abbreviation"`
	Teamname        string     `json:"teamName"`
	Locationname    string     `json:"locationName"`
	Firstyearofplay string     `json:"firstYearOfPlay"`
	Division        Division   `json:"division"`
	Conference      Conference `json:"conference"`
	Franchise       Franchise  `json:"franchise"`
	Shortname       string     `json:"shortName"`
	Officialsiteurl string     `json:"officialSiteUrl"`
	Franchiseid     int        `json:"franchiseId"`
	Active          bool       `json:"active"`
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
