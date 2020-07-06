package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type ACNH struct {
	Bugs   []Bug  `json:"bugs"`
	Fishes []Fish `json:"fishes"`
}

type Fish struct {
	Bug
	ShadowSize string `json:"shadow_size"`
}

type Bug struct {
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Months   []int  `json:"months"`
	Hours    []int  `json:"hours"`
	Location string `json:"location"`
}

var ErrBadInput = errors.New("that's some bad input")

var months = map[string]int{
	"jan": 0,
	"feb": 1,
	"mar": 2,
	"apr": 3,
	"may": 4,
	"jun": 5,
	"jul": 6,
	"aug": 7,
	"sep": 8,
	"oct": 9,
	"nov": 10,
	"dec": 11,
}

func main() {
	bugs, err := processBugs()
	if err != nil {
		log.Fatal(err)
	}

	fishes, err := processFish()
	if err != nil {
		log.Fatal(err)
	}

	acnh := ACNH{
		Bugs:   bugs,
		Fishes: fishes,
	}

	outFull, err := os.Create("acnh.json")
	if err != nil {
		log.Fatal(err)
	}
	defer outFull.Close()
	json.NewEncoder(outFull).Encode(acnh)
}

func processBugs() ([]Bug, error) {
	var bugs []Bug
	inFile, err := os.Open("bugs.csv")
	if err != nil {
		return bugs, fmt.Errorf("unable to open bugs CSV file: %w", err)
	}
	defer inFile.Close()

	r := csv.NewReader(inFile)
	// Read off header line
	_, err = r.Read()
	if err != nil {
		return bugs, fmt.Errorf("somehow errored reading header line on bugs input file: %w", err)
	}

	for {
		fields, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Reached end of bugs input file.")
				break
			}
		}

		fmt.Printf("processing %s\n", fields[0])
		price, err := strconv.Atoi(fields[1])
		if err != nil {
			return bugs, fmt.Errorf("price '%s' in '%s' was not a valid int: %w", fields[1], fields[0], err)
		}
		months, err := parseMonths(fields[2])
		if err != nil {
			return bugs, fmt.Errorf("months '%s' in '%s' was not a valid month range: %w", fields[2], fields[0], err)
		}
		hours, err := parseHours(fields[3])
		if err != nil {
			return bugs, fmt.Errorf("hours '%s' in '%s' was not a valid hour range: %w", fields[3], fields[0], err)
		}

		bug := Bug{
			fields[0],
			price,
			months,
			hours,
			fields[4],
		}
		bugs = append(bugs, bug)
	}

	return bugs, nil
}

func processFish() ([]Fish, error) {
	var fishes []Fish
	inFile, err := os.Open("fish.csv")
	if err != nil {
		return fishes, fmt.Errorf("unable to open fish CSV file: %w", err)
	}
	defer inFile.Close()

	r := csv.NewReader(inFile)
	// Read off header line
	_, err = r.Read()
	if err != nil {
		return fishes, fmt.Errorf("somehow errored reading header line on fishes input file: %w", err)
	}

	for {
		fields, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Reached end of fish input file.")
				break
			}
		}

		fmt.Printf("processing %s\n", fields[0])
		price, err := strconv.Atoi(fields[1])
		if err != nil {
			return fishes, fmt.Errorf("price '%s' in '%s' was not a valid int: %w", fields[1], fields[0], err)
		}
		months, err := parseMonths(fields[4])
		if err != nil {
			return fishes, fmt.Errorf("months '%s' in '%s' was not a valid month range: %w", fields[4], fields[0], err)
		}
		hours, err := parseHours(fields[3])
		if err != nil {
			return fishes, fmt.Errorf("hours '%s' in '%s' was not a valid hour range: %w", fields[3], fields[0], err)
		}

		fish := Fish{
			Bug{
				fields[0],
				price,
				months,
				hours,
				fields[2],
			},
			fields[5],
		}
		fishes = append(fishes, fish)
	}

	return fishes, nil
}

func rangeString(r string, mod int) ([]int, error) {
	pair := strings.Split(strings.TrimSpace(r), "-")
	if len(pair) != 2 {
		return []int{}, errors.New("range must be two numbers separated by a '-'")
	}

	min, err := strconv.Atoi(pair[0])
	if err != nil {
		return []int{}, fmt.Errorf("min in range was not an int: %w", err)
	}

	max, err := strconv.Atoi(pair[1])
	if err != nil {
		return []int{}, fmt.Errorf("max in range was not an int: %w", err)
	}

	return rangeNums(min, max, mod), nil
}

func rangeNums(min, max, mod int) []int {
	if max < min {
		max += mod
	}

	nums := make([]int, max-min+1)
	for i := range nums {
		nums[i] = (min + i) % mod
	}
	return nums
}

func rangeMonths(r string) ([]int, error) {
	pair := strings.Split(strings.TrimSpace(r), "-")
	if len(pair) != 2 {
		return []int{}, errors.New("range must be 'month-month'")
	}

	minMonth := strings.ToLower(strings.TrimSpace(pair[0]))[:3]
	min, ok := months[minMonth]
	if !ok {
		return []int{}, errors.New("couldn't find month " + minMonth)
	}

	maxMonth := strings.ToLower(strings.TrimSpace(pair[1]))[:3]
	max, ok := months[maxMonth]
	if !ok {
		return []int{}, errors.New("couldn't find month " + maxMonth)
	}

	return rangeNums(min, max, 12), nil
}

func parseMonths(ms string) ([]int, error) {
	if strings.TrimSpace(strings.ToLower(ms)) == "all" {
		return rng(0, 11), nil
	}

	splits := strings.Split(ms, ",")
	for i := range splits {
		splits[i] = strings.TrimSpace(splits[i])
	}

	// Single month
	if len(splits) == 1 && !strings.HasPrefix(strings.ToLower(splits[0]), "all except ") {
		mi, ok := months[strings.ToLower(splits[0])[:3]]
		if !ok {
			return []int{}, errors.New(fmt.Sprintf("single element in months that was neither a valid month or 'all': '%s'", splits[0]))
		}
		return []int{mi}, nil
	}

	if strings.HasPrefix(strings.ToLower(splits[0]), "all except ") {
		splits[0] = splits[0][11:]
		splits = invertMonths(splits)
	}

	// Multiple months
	var mons []int
	for _, m := range splits {
		if m == "" {
			continue
		}
		mi, ok := months[strings.ToLower(m)[:3]]
		if !ok {
			return []int{}, errors.New(fmt.Sprintf("element in list was neither a valid month or 'all': '%s'", m))
		}
		mons = append(mons, mi)
	}

	return mons, nil
}

func parseHours(hs string) ([]int, error) {
	hs = strings.TrimSpace(strings.ToLower(hs))
	if hs == "all" {
		return rng(0, 23), nil
	}

	hours := make([]int, 0, 24)
	splits := strings.Split(hs, ",")
	for i := range splits {
		pair := strings.Split(strings.TrimSpace(splits[i]), "-")
		if len(pair) != 2 {
			return hours, errors.New(fmt.Sprintf("unknown format: '%s'", splits[i]))
		}
		start, err := parseTime(pair[0])
		if err != nil {
			return hours, err
		}
		end, err := parseTime(pair[1])
		if err != nil {
			return hours, err
		}
		// Weirdness if time boundary goes past midnight
		if end < start {
			hours = append(hours, rng(start, 23)...)
			hours = append(hours, rng(0, end-1)...)
		} else {
			hours = append(hours, rng(start, end-1)...)
		}
	}

	return hours, nil
}

func parseTime(ts string) (int, error) {
	ts = strings.TrimSpace(strings.ToLower(ts))
	if !strings.HasSuffix(ts, "am") && !strings.HasSuffix(ts, "pm") {
		return 0, errors.New(fmt.Sprintf("time must end with either 'am' or 'pm': '%s'", ts))
	}
	// Special case 12am and 12pm because time is weird.
	if ts == "12pm" {
		return 12, nil
	}
	if ts == "12am" {
		return 0, nil
	}
	var isPM bool
	if strings.HasSuffix(ts, "pm") {
		isPM = true
	}
	hour, err := strconv.Atoi(ts[:len(ts)-2])
	if err != nil {
		return 0, err
	}
	if isPM {
		hour += 12
	}
	return hour % 24, nil
}

func invertMonths(ms []string) []string {
	allMonths := make(map[string]bool)
	mons := []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}
	for _, m := range mons {
		allMonths[m] = true
	}

	for _, m := range ms {
		allMonths[strings.ToLower(m)[:3]] = false
	}

	var invertedMonths []string
	for k, v := range allMonths {
		if v {
			invertedMonths = append(invertedMonths, k)
		}
	}

	return invertedMonths
}

func rng(min, max int) []int {
	if min > max {
		min, max = max, min
	}
	r := make([]int, max-min+1)
	for i := 0; i < max-min+1; i++ {
		r[i] = i + min
	}
	return r
}
