package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/trelore/nhl"
	"github.com/wcharczuk/go-chart"
)

type chartData struct {
	series chart.TimeSeries
	max    float64
	min    float64
}

func main() {
	stat := flag.String("stat", "points", "the stat to measure (i.e. points, goals)")
	outputFile := flag.String("output", "leaders", "the file name i.e. 'top10_points'")
	formatFlag := flag.String("format", "svg", "the file format SVG or PNG")
	teamFlag := flag.Int("team", -1, "Add team to output")
	flag.Parse()

	if err := run(*stat, *outputFile, *formatFlag, *teamFlag); err != nil {
		log.Fatal(err)
	}
}

func run(stat, outputFile, formatFlag string, teamFlag int) error {
	n := nhl.NewClient()

	var err error
	var playerIDs []string
	if teamFlag != -1 {
		playerIDs, err = n.GetTeamPlayerIDs(teamFlag)
		if err != nil {
			return err
		}
	}

	playerIDs = append(playerIDs, flag.Args()...)

	yAxisMax, yAxisMin := float64(0), math.MaxFloat64
	var series []chart.Series
	var wg sync.WaitGroup
	for _, playerID := range playerIDs {
		wg.Add(1)
		go func(playerID string) {
			defer wg.Done()
			chartData, err := getData(n, playerID, stat)
			if err != nil {
				log.Printf("%s", err)
				return
			}
			series = append(series, chartData.series)
			// possible race condition in multiple go routines setting minimum?
			yAxisMin = math.Min(yAxisMin, chartData.min)
			yAxisMax = math.Max(yAxisMax, chartData.max)
		}(playerID)
	}
	wg.Wait()

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Date",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:      stat,
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

	outputFormat := getFileExtension(formatFlag)

	f, err := os.Create(fmt.Sprintf("%s.%s", outputFile, outputFormat.name))
	if err != nil {
		log.Fatalf("error creating file %s", err)
	}
	defer f.Close()

	if err := graph.Render(outputFormat.renderer, f); err != nil {
		return err
	}

	return nil
}

type outputFormat struct {
	name     string
	renderer chart.RendererProvider
}

func getFileExtension(desiredExtension string) outputFormat {
	switch strings.ToLower(desiredExtension) {
	case "svg":
		return outputFormat{name: "svg", renderer: chart.SVG}
	case "png":
		return outputFormat{name: "png", renderer: chart.PNG}
	default:
		log.Printf("Desired extension '%s' not matched. Using svg.\n", desiredExtension)
		return outputFormat{name: "svg", renderer: chart.SVG}
	}
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

func getData(n nhl.Client, playerID string, stat string) (*chartData, error) {
	player, err := n.GetPlayer(playerID)
	if err != nil {
		return nil, err
	}

	resp, err := n.GameLogStats(player)
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
