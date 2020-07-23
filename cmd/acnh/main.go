package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

type Bug struct {
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Months   []int  `json:"months"`
	Hours    []int  `json:"hours"`
	Location string `json:"location"`
	HourMap  map[int]bool
	Timing   Timing
}

func (b *Bug) SetHourMap(m map[int]bool) {
	b.HourMap = m
}

func (b *Bug) GetHours() []int {
	return b.Hours
}

type Fish struct {
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Months     []int  `json:"months"`
	Hours      []int  `json:"hours"`
	Location   string `json:"location"`
	ShadowSize string `json:"shadow_size"`
	HourMap    map[int]bool
	Timing     Timing
}

type Umbrella struct {
	Name               string `json:"name"`
	DIY                string `json:"diy"`
	BuyPrice           string `json:"buy_price"`
	SellPrice          string `json:"sell_price"`
	HHABase            string `json:"hha_base"`
	Color1             string `json:"color_1"`
	Color2             string `json:"color_2"`
	Size               string `json:"size"`
	MilesPrice         string `json:"miles_price"`
	Source             string `json:"source"`
	SourceNotes        string `json:"source_notes"`
	VillagerEquippable string `json:"villager_equippable"`
	Catalog            string `json:"catalog"`
}

func (f *Fish) SetHourMap(m map[int]bool) {
	f.HourMap = m
}

func (f *Fish) GetHours() []int {
	return f.Hours
}

type SeaCreature struct {
	Name    string `json:"name"`
	Price   int    `json:"price"`
	Hours   []int  `json:"hours"`
	Months  []int  `json:"months"`
	HourMap map[int]bool
	Timing  Timing
}

func (s *SeaCreature) SetHourMap(m map[int]bool) {
	s.HourMap = m
}

func (s *SeaCreature) GetHours() []int {
	return s.Hours
}

type HourMapper interface {
	SetHourMap(map[int]bool)
	GetHours() []int
}

type ACNH struct {
	Bugs         []Bug         `json:"bugs"`
	Fishes       []Fish        `json:"fishes"`
	SeaCreatures []SeaCreature `json:"sea_creatures"`
	Umbrellas    []Umbrella    `json:"umbrellas"`
}

type Timing struct {
	AvailableNow    bool
	AvailableAt     int
	AvailableUntil  int
	AvailableAllDay bool
	CurrentHour     int
}

func (t *Timing) DisplayAt() string {
	if t.AvailableAt == 0 {
		return "12AM"
	}
	if t.AvailableAt == 12 {
		return "12PM"
	}
	if t.AvailableAt > 12 {
		return fmt.Sprintf("%dPM", t.AvailableAt-12)
	}
	return fmt.Sprintf("%dAM", t.AvailableAt)
}

func (t *Timing) DisplayUntil() string {
	if t.AvailableUntil == 0 {
		return "12AM"
	}
	if t.AvailableUntil == 12 {
		return "12PM"
	}
	if t.AvailableUntil > 12 {
		return fmt.Sprintf("%dPM", t.AvailableUntil-12)
	}
	return fmt.Sprintf("%dAM", t.AvailableUntil)
}

func main() {
	logger := StdLogger{}
	critters, err := loadCritters(logger)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := loadTemplate(logger)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", mainHandler(critters, tmpl, logger))
	http.HandleFunc("/sortable.js", sortableHandler(logger))
	http.HandleFunc("/style.css", cssHandler(logger))
	http.HandleFunc("/acnh.js", jsHandler(logger))
	logger.Log("Starting server", "port", "80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func sortableHandler(logger Logger) http.HandlerFunc {
	file, err := os.Open("js/sortable.js")
	if err != nil {
		logger.Log("failed to open sortable.js", "error", err)
		return http.NotFound
	}
	defer file.Close()
	js, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Log("failed to read sortable.js", "error", err)
		return http.NotFound
	}
	jsStr := string(js)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprint(w, jsStr)
	}
}

func cssHandler(logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("css/style.css")
		if err != nil {
			logger.Log("failed to open style.css", "error", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer file.Close()
		css, err := ioutil.ReadAll(file)
		if err != nil {
			logger.Log("failed to read style.css", "error", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		cssStr := string(css)
		w.Header().Set("Content-Type", "text/css")
		fmt.Fprint(w, cssStr)
	}
}

func jsHandler(logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("js/acnh.js")
		if err != nil {
			logger.Log("failed to open style.css", "error", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer file.Close()
		js, err := ioutil.ReadAll(file)
		if err != nil {
			logger.Log("failed to open style.css", "error", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		jsStr := string(js)
		w.Header().Set("Content-type", "application/javascript")
		fmt.Fprint(w, jsStr)
	}
}

func mainHandler(critters ACNH, tmpl *template.Template, logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if tmpl == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "the Los Angeles time zone no longer exists...that's bad.")
			logger.Log("failed loading timezone data", "error", err)
			return
		}
		t := time.Now().In(loc)
		var filteredCritters ACNH

		for _, bug := range critters.Bugs {
			if contains(bug.Months, int(t.Month())-1) {
				b := bug
				b.Timing = timing(bug.HourMap, t.Hour())
				filteredCritters.Bugs = append(filteredCritters.Bugs, b)
			}
		}
		for _, fish := range critters.Fishes {
			if contains(fish.Months, int(t.Month())-1) {
				fish.Timing = timing(fish.HourMap, t.Hour())
				filteredCritters.Fishes = append(filteredCritters.Fishes, fish)
			}
		}
		for _, sc := range critters.SeaCreatures {
			if contains(sc.Months, int(t.Month())-1) {
				sc.Timing = timing(sc.HourMap, t.Hour())
				filteredCritters.SeaCreatures = append(filteredCritters.SeaCreatures, sc)
			}
		}

		filteredCritters.Umbrellas = critters.Umbrellas

		tmpl.Execute(w, filteredCritters)
	}
}

func loadCritters(logger Logger) (ACNH, error) {
	var critters ACNH
	file, err := os.Open("acnh.json")
	if err != nil {
		return critters, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&critters)
	for i := range critters.Bugs {
		setHourMap(&critters.Bugs[i])
	}
	for i := range critters.Fishes {
		setHourMap(&critters.Fishes[i])
	}
	for i := range critters.SeaCreatures {
		setHourMap(&critters.SeaCreatures[i])
	}
	return critters, err
}

func setHourMap(h HourMapper) {
	hours := h.GetHours()
	hourMap := make(map[int]bool)
	for i := range hours {
		hourMap[hours[i]] = true
	}

	h.SetHourMap(hourMap)
}

func loadTemplate(logger Logger) (*template.Template, error) {
	file, err := os.Open("templates/main.html")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	text, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return template.New("main").Parse(string(text))
}

func contains(s []int, n int) bool {
	for i := range s {
		if s[i] == n {
			return true
		}
	}
	return false
}

func timing(s map[int]bool, n int) Timing {
	if len(s) == 24 {
		return Timing{
			AvailableNow:    true,
			AvailableAllDay: true,
			CurrentHour:     n,
		}
	}

	availableNow := s[n]
	var availableUntil int
	var availableAt int
	if availableNow {
		for i := n; i < n+24; i++ {
			if !s[i%24] {
				availableUntil = i % 24
				break
			}
		}
	} else {
		for i := n; i < n+24; i++ {
			if s[i%24] {
				availableAt = i % 24
				break
			}
		}
	}

	return Timing{
		AvailableNow:   availableNow,
		AvailableUntil: availableUntil,
		AvailableAt:    availableAt,
		CurrentHour:    n,
	}
}
