package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	chart "github.com/wcharczuk/go-chart"
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

type peopleResponse struct {
	People []people `json:"people"`
}

type people struct {
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
}

type chartData struct {
	series chart.TimeSeries
	max    float64
	min    float64
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

func verifyPlayer(playerID int) (*people, error) {
	playerURL := fmt.Sprintf("people/%d", playerID)
	resp := new(peopleResponse)
	err := getJSON(hostURL+playerURL, resp)
	if err != nil {
		return nil, fmt.Errorf("could not verify player %d: %s", playerID, err)
	}
	return &resp.People[0], nil
}

func getPlayerData(playerID int, stat string) (*chartData, error) {
	player, err := verifyPlayer(playerID)
	if err != nil {
		return nil, err
	}

	playerStats := fmt.Sprintf("people/%d/stats?stats=gameLog", playerID)
	resp := new(statsResponse)
	err = getJSON(hostURL+playerStats, resp)
	if err != nil {
		return nil, err
	}

	splits := resp.Stats[0].Splits
	xValues := make([]time.Time, len(splits))
	yValues := make([]float64, len(splits))

	plot := 0
	max := float64(0)
	min := math.MaxFloat64
	for i := len(splits) - 1; i >= 0; i-- {
		xValues[i], err = time.Parse("2006-01-02", splits[i].Date)
		if err != nil {
			log.Fatalf("could not get parse date: %s", err)
		}
		switch strings.ToUpper(stat) {
		case "POINTS", "P":
			plot += splits[i].Stat.Points
		case "GOALS", "G":
			plot += splits[i].Stat.Goals
		case "ASSISTS", "A":
			plot += splits[i].Stat.Assists
		case "PLUSMINUS":
			plot += splits[i].Stat.PlusMinus
		case "PPG":
			plot += splits[i].Stat.PowerPlayGoals
		case "PPP":
			plot += splits[i].Stat.PowerPlayPoints
		case "SHG":
			plot += splits[i].Stat.ShortHandedGoals
		case "SHP":
			plot += splits[i].Stat.ShortHandedPoints
		case "SHOTS":
			plot += splits[i].Stat.Shots
		case "PIM":
			plot += splits[i].Stat.Pim
		case "GWG":
			plot += splits[i].Stat.GameWinningGoals
		case "OTG":
			plot += splits[i].Stat.OverTimeGoals
		default:
			plot += splits[i].Stat.Points
		}
		min = math.Min(min, float64(plot))
		max = math.Max(max, float64(plot))
		yValues[i] = float64(plot)
	}
	return &chartData{
		series: chart.TimeSeries{
			Name:    fmt.Sprintf("%s, %s", player.LastName, player.FirstName),
			XValues: xValues,
			YValues: yValues,
		},
		max: float64(max),
		min: float64(min),
	}, nil
}

func getPlayers(input []string) []int {
	a := make([]int, len(input))
	for i, v := range input {
		val, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("error converting input %s to int", v)
		}
		a[i] = val
	}
	return a
}

func getLines(min, max, count float64) []float64 {
	min = math.Min(0, min)
	rangeVal := max - min
	tickSize := math.Ceil(rangeVal / count)
	values := make([]float64, int(count)+1)
	for i := float64(0); i <= count; i++ {
		values[int(i)] = i*tickSize + min
	}
	return values
}

func getTicks(min, max, count float64) chart.Ticks {
	values := getLines(min, max, count)
	ticks := make([]chart.Tick, int(count)+1)
	for i := 0; i <= int(count); i++ {
		ticks[i] = chart.Tick{Value: values[i], Label: fmt.Sprintf("%.f", values[i])}
	}
	return ticks
}

func getGridLines(min, max, count float64) []chart.GridLine {
	values := getLines(min, max, count)
	gridLines := make([]chart.GridLine, int(count)+1)
	for i := 0; i <= int(count); i++ {
		gridLines[i] = chart.GridLine{Value: values[i]}
	}
	return gridLines
}

func main() {
	var stat = flag.String("stat", "points", "the stat to measure (i.e. points, goals)")
	var outputFile = flag.String("o", "leaders.png", "the file name i.e. 'top10_points.png'")
	flag.Parse()

	players := getPlayers(flag.Args())
	var series []chart.Series
	yAxisMax, yAxisMin := float64(0), math.MaxFloat64
	for _, player := range players {
		chartData, err := getPlayerData(player, *stat)
		if err != nil {
			log.Printf("%s", err)
			continue
		}
		series = append(series, chartData.series)
		yAxisMin = math.Min(yAxisMin, chartData.min)
		yAxisMax = math.Max(yAxisMax, chartData.max)
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Date",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:      *stat,
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			Ticks:     getTicks(yAxisMin, yAxisMax, 8),
			GridLines: getGridLines(yAxisMin, yAxisMax, 8),
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.ColorBlack,
				StrokeWidth: 0.2,
			},
		},
		Series: series,
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	f, err := os.Create(fmt.Sprintf("%s", *outputFile))
	if err != nil {
		log.Fatalf("error creating file %s", err)
	}
	defer f.Close()

	graph.Render(chart.PNG, f)
}
