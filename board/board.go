package board

import (
	"embed"
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

type line_id string
type hex_color string
type URI string

var LinesToURI map[line_id]URI
var LinesToColors map[line_id]hex_color
var StopsToFriendlies map[string]string
var StopsToDestinations map[string]string
var LinesToDirections map[string]string

//go:embed stops.txt
var stopsFS embed.FS

//go:embed boroughs.geojson
var boroughsFS embed.FS

var Terminals = map[string][2]string{
	"1":  {"Van Cortlandt Park-242 St", "South Ferry"},
	"2":  {"Wakefield-241 St", "Flatbush Av-Brooklyn College"},
	"3":  {"Harlem-148 St", "New Lots Av"},
	"4":  {"Woodlawn", "Crown Heights-Utica Av"},
	"5":  {"Eastchester-Dyre Av", "Flatbush Av-Brooklyn College"},
	"6":  {"Pelham Bay Park", "Brooklyn Bridge-City Hall"},
	"6X": {"Pelham Bay Park", "Brooklyn Bridge-City Hall"},
	"7":  {"Flushing-Main St", "34 St-Hudson Yards"},
	"A":  {"Inwood-207 St", "Far Rockaway-Mott Av"},
	"C":  {"168 St", "Euclid Av"},
	"E":  {"Jamaica Center-Parsons/Archer", "World Trade Center"},
	"B":  {"Bedford Park Blvd", "Brighton Beach"},
	"D":  {"Norwood-205 St", "Coney Island-Stillwell Av"},
	"F":  {"Jamaica-179 St", "Coney Island-Stillwell Av"},
	"FX": {"Jamaica-179 St", "Coney Island-Stillwell Av"},
	"M":  {"Forest Hills-71 Av", "Metropolitan Av"},
	"N":  {"Astoria-Ditmars Blvd", "Coney Island-Stillwell Av"},
	"Q":  {"96 St", "Coney Island-Stillwell Av"},
	"R":  {"Forest Hills-71 Av", "Bay Ridge-95 St"},
	"W":  {"Astoria-Ditmars Blvd", "Whitehall St"},
	"G":  {"Court Sq", "Church Av"},
	"J":  {"Jamaica Center-Parsons/Archer", "Broad St"},
	"Z":  {"Jamaica Center-Parsons/Archer", "Broad St"},
	"L":  {"8 Av", "Canarsie-Rockaway Pkwy"},
	// "GS": {"Grand Central-42 St", "Times Sq-42 St"},
	"GS": {"Times Sq-42 St", "Grand Central-42 St"},
	"FS": {"Franklin Av", "Prospect Park"},
	"H":  {"Broad Channel", "Rockaway Park"},

	/* TODO add shuttles */
}

type Arrival struct {
	TimeToArrival    int64  `json:"time_to_arrival"`
	Line             string `json:"line"` /* adheres to routes.txt */
	Direction        string `json:"direction"`
	Color            string `json:"color"`
	Stop             string `json:"stop"`
	Friendly_Stop    string `json:"friendly_stop,omitempty"` /* adheres to stops.txt */
	Destination_Stop string `json:"destination_stop,omitempty"`
	Destination_Dir  string `json:"destination_direction,omitempty"`
}

type borough struct {
	Name     string
	Geometry orb.Geometry
}

func init() {
	LinesToURI = make(map[line_id]URI)
	LinesToColors = make(map[line_id]hex_color)
	StopsToFriendlies = make(map[string]string)
	StopsToDestinations = make(map[string]string)

	init_linestoURI_linestoColors()
	init_stopstoFriendlies()
}

func init_linestoURI_linestoColors() {
	/* defined in https://www.mta.info/document/168976 */
	blue := hex_color("#0062CF")
	orange := hex_color("#EB6800")
	light_green := hex_color("#799534")
	brown := hex_color("#8E5C33")
	grey := hex_color("#7C858C")
	yellow := hex_color("#F6BC26")
	red := hex_color("#D82233")
	dark_green := hex_color("#009952")
	purple := hex_color("#9A38A1")
	// teal := hex_color("#0078C6")
	mta_blue := hex_color("#08179C")
	// isa_blue := hex_color("#0078C6")

	/* defined in ../subway_config/routes.txt */
	IND_eight_ave := []line_id{"A", "C", "E", "H"}
	IND_sixth_ave := []line_id{"B", "D", "F", "FX", "M", "FS"}
	IND_crosstown := []line_id{"G"}
	BMT_canarsie := []line_id{"L"}
	BMT_nassau := []line_id{"J", "Z"}
	BMT_broadway := []line_id{"N", "Q", "R", "W"} // my favorite
	IRT_broadway_seventh := []line_id{"1", "2", "3"}
	IRT_lex_ave := []line_id{"4", "5", "6", "6X"} // my second favorite
	IRT_flushing := []line_id{"7", "7X", "GS"}    // my third favorite
	// if only!
	// IND_second_ave := []string{"T"}

	SIR := []line_id{"SI"}

	for _, line := range IND_eight_ave {
		LinesToURI[line] = "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-ace"
		LinesToColors[line] = blue
	}
	for _, line := range IND_sixth_ave {
		LinesToURI[line] = "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-bdfm"
		LinesToColors[line] = orange
	}
	for _, line := range IND_crosstown {
		LinesToURI[line] = "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-g"
		LinesToColors[line] = light_green
	}
	for _, line := range BMT_canarsie {
		LinesToURI[line] = "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-l"
		LinesToColors[line] = grey
	}
	for _, line := range BMT_nassau {
		LinesToURI[line] = "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-jz"
		LinesToColors[line] = brown
	}
	for _, line := range BMT_broadway {
		LinesToURI[line] = "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-nqrw"
		LinesToColors[line] = yellow
	}

	for _, line := range IRT_broadway_seventh {
		LinesToColors[line] = red
	}
	for _, line := range IRT_lex_ave {
		LinesToColors[line] = dark_green
	}

	for _, line := range IRT_flushing {
		LinesToColors[line] = purple
	}

	trains := append(append(IRT_broadway_seventh, IRT_lex_ave...), IRT_flushing...)

	for _, line := range trains {
		LinesToURI[line] = "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs"
	}
	for _, line := range SIR {
		LinesToURI[line] = "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs-si"
		LinesToColors[line] = mta_blue
	}
}

func init_stopstoFriendlies() {
	data, err := stopsFS.ReadFile("stops.txt")
	if err != nil {
		panic("failed to open stops.txt: " + err.Error())
	}

	reader := csv.NewReader(strings.NewReader(string(data)))

	f_lines, err := reader.ReadAll() // [][]string
	if err != nil {
		panic("failed to read CSV: " + err.Error())
	}

	data1, err := boroughsFS.ReadFile("boroughs.geojson")
	if err != nil {
		panic("failed to read geoJSON: " + err.Error())
	}

	fc, err := geojson.UnmarshalFeatureCollection(data1)
	if err != nil {
		panic("failed to read geoJSON: " + err.Error())
	}

	boroughs := make([]borough, 0, len(fc.Features))
	for _, f := range fc.Features {
		name, _ := f.Properties["BoroName"].(string)
		boroughs = append(boroughs, borough{
			Name:     name,
			Geometry: f.Geometry,
		})
	}

	for _, line := range f_lines[1:] {
		stop_id := line[0]
		friendly := line[1]
		lat, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			panic("error: " + err.Error())
		}
		lon, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			panic("error: " + err.Error())
		}
		StopsToFriendlies[stop_id] = friendly
		StopsToDestinations[stop_id] = init_stopToDestination(lon, lat, boroughs)
	}
}

func init_stopToDestination(lon float64, lat float64, boroughs []borough) string {
	pt := orb.Point{lon, lat}
	for _, b := range boroughs {
		switch g := b.Geometry.(type) {
		case orb.Polygon:
			if planar.PolygonContains(g, pt) {
				return b.Name
			}
		case orb.MultiPolygon:
			if planar.MultiPolygonContains(g, pt) {
				return b.Name
			}
		}
	}
	panic("Unknown stop")
	return ""
}

func fill_destination_stop(line string, direction string) string {
	terms, ok := Terminals[line]
	if !ok {
		return ""
	}

	if direction == "N" {
		return terms[0] // north terminal
	}
	return terms[1] // south terminal
}

func fill_destination_dir(line string, borough string, direction string) string {
	terms, ok := Terminals[line]
	if !ok {
		return ""
	}

	north := terms[0]
	south := terms[1]

	switch borough {
	case "Manhattan":
		if direction == "N" {
			return "Uptown to " + north
		}
		return "Downtown to " + south

	case "Brooklyn":
		if direction == "N" {
			return "To " + north
		}
		return "To " + south

	case "Queens":
		if direction == "N" {
			return "To " + north
		}
		return "To " + south

	case "Bronx":
		if direction == "N" {
			return "To " + north
		}
		return "To " + south

	default:
		return "To " + terms[1]
	}
}

func fill_arrival(entity *gtfs.FeedEntity, line string, stop_id string, i int) []Arrival {
	if entity.TripUpdate == nil || entity.TripUpdate.Trip == nil {
		return []Arrival{}
	}

	if line != *entity.TripUpdate.Trip.RouteId {
		return []Arrival{}
	}

	var stop_time_update *gtfs.TripUpdate_StopTimeUpdate

	stop_time_update = nil
	for _, update := range entity.TripUpdate.StopTimeUpdate {
		if stop_time_update == nil && update.StopId != nil && stop_id == *update.StopId {
			stop_time_update = update
		}
	}

	if stop_time_update == nil {
		/* this train isn't coming to this station any time soon buddy */
		return []Arrival{{}}
	}
	var arrival_time int64

	if stop_time_update.Arrival != nil && stop_time_update.Arrival.Time != nil {
		arrival_time = *stop_time_update.Arrival.Time
	} else if stop_time_update.Departure != nil && stop_time_update.Departure.Time != nil {
		arrival_time = *stop_time_update.Departure.Time
	} else {
		panic("no arrival time found")
	}

	time_to_arrival := arrival_time - time.Now().Unix()

	var direction int
	var dir_str string
	if length := len(stop_id); stop_id[length-1] == 'N' {
		direction = 0
		dir_str = "N"
	} else {
		direction = 1
		dir_str = "S"
	}

	friendly_name := StopsToFriendlies[stop_id]

	return []Arrival{{
		Stop:             stop_id,
		TimeToArrival:    time_to_arrival,
		Direction:        strconv.Itoa(direction),
		Color:            string(LinesToColors[line_id(line)]),
		Line:             line,
		Friendly_Stop:    friendly_name,
		Destination_Stop: fill_destination_stop(line, dir_str),
		Destination_Dir:  fill_destination_dir(line, StopsToDestinations[stop_id], dir_str)}}
}

func request(line string, stop string) ([]Arrival, error) {
	URI, ok := LinesToURI[line_id(line)]
	if !ok {
		return nil, errors.New("Line not found: " + line)
	}
	resp, err := http.Get(string(URI))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // close at the end

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	feed := &gtfs.FeedMessage{}
	// fun fact: unmarshal means to take it "out of order" like military ppl
	if err := proto.Unmarshal(body, feed); err != nil {
		return nil, err
	}

	arrivals := make([]Arrival, 0)
	length := 0
	for i, entity := range feed.Entity {
		if new_arrivals := fill_arrival(entity, line, stop, i); len(new_arrivals) == 1 {
			new_arrival := new_arrivals[0]
			if new_arrival.Line != "" {
				arrivals = append(arrivals, new_arrival)
			}
		}
		length++
	}
	return arrivals, nil
}

func GetArrivals(line string, stop string, N int) ([]Arrival, error) {
	_, ok := StopsToFriendlies[stop]
	if !ok {
		return nil, errors.New("Stop not found: " + stop)
	}
	allArrivals, err := request(line, stop)
	if err != nil {
		return nil, errors.New("Line not found: " + line)
	}

	sort.Slice(allArrivals, func(i, j int) bool {
		return allArrivals[i].TimeToArrival < allArrivals[j].TimeToArrival
	})

	arrivals := make([]Arrival, 0, len(allArrivals))
	for _, a := range allArrivals {
		if a.TimeToArrival > 0 && a.Line != "" {
			// new_a := Arrival{TimeToArrival: a.TimeToArrival / 60, Line: a.Line, Color: a.Color, Direction: a.Direction, Stop: a.Stop}
			arrivals = append(arrivals, a)
		}
	}

	return arrivals[:min(N, len(arrivals))], nil
}
