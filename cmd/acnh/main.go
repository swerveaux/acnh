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
	critters, err := loadCritters()
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := loadTemplate()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", mainHandler(critters, tmpl))
	http.HandleFunc("/sortable.js", sortableHandler())
	http.HandleFunc("/style.css", cssHandler())
	http.HandleFunc("/acnh.js", jsHandler())
	log.Fatal(http.ListenAndServe(":80", nil))
}

func sortableHandler() http.HandlerFunc {
	file, _ := os.Open("js/sortable.js")
	defer file.Close()
	js, _ := ioutil.ReadAll(file)
	jsStr := string(js)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprint(w, jsStr)
	}
}

func cssHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _ := os.Open("css/style.css")
		defer file.Close()
		css, _ := ioutil.ReadAll(file)
		cssStr := string(css)
		w.Header().Set("Content-Type", "text/css")
		fmt.Fprint(w, cssStr)
	}
}

func jsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _ := os.Open("js/acnh.js")
		defer file.Close()
		js, _ := ioutil.ReadAll(file)
		jsStr := string(js)
		w.Header().Set("Content-type", "application/javascript")
		fmt.Fprint(w, jsStr)
	}
}

func mainHandler(critters ACNH, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if tmpl == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "the Los Angeles time zone no longer exists...that's bad.")
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

		tmpl.Execute(w, filteredCritters)
	}
}

func loadCritters() (ACNH, error) {
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

func loadTemplate() (*template.Template, error) {
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
