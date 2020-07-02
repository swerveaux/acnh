package main

import (
	"bufio"
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

var ErrNoMoreInput = errors.New("end of input")
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
	var acnh ACNH
	err := loadJSON(&acnh)
	if err != nil {
		log.Fatal(err)
	}

	var bugMode, fishMode bool
	if len(os.Args) > 1 {
		if os.Args[1] == "-b" {
			bugMode = true
		}
		if os.Args[1] == "-f" {
			fishMode = true
		}
	}

	if bugMode {
		for {
			bug, err := getBug()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatal(err)
			}
			acnh.Bugs = append(acnh.Bugs, bug)
		}
	}

	if fishMode {
		for {
			fish, err := getFish()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatal(err)
			}
			acnh.Fishes = append(acnh.Fishes, fish)
		}
	}

	err = saveJSON(acnh)
	if err != nil {
		log.Fatal(err)
	}
}

func loadJSON(acnh *ACNH) error {
	f, err := os.Open("acnh.json")
	if err != nil {
		return fmt.Errorf("unable to open acnh.json: %w", err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(acnh)
	if err != nil {
		return fmt.Errorf("unable to decode acnh.json: %w", err)
	}

	return nil
}

func getFish() (Fish, error) {
	var fish Fish

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Fish's name? ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return fish, fmt.Errorf("failed reading name: %w", err)
	}
	fish.Name = strings.TrimSpace(name)

	fmt.Print("...price? ")
	price, err := reader.ReadString('\n')
	if err != nil {
		return fish, fmt.Errorf("failed reading price: %w", err)
	}
	p, err := strconv.Atoi(strings.TrimSpace(price))
	if err != nil {
		return fish, ErrBadInput
	}
	fish.Price = p

	fmt.Print("...location? ")
	location, err := reader.ReadString('\n')
	if err != nil {
		return fish, fmt.Errorf("failed reading location: %w", err)
	}
	fish.Location = strings.TrimSpace(location)

	fmt.Print("...hours (24 hour ints in range, e.g. 8-17 for 8am to 5pm, or 'all'? ")
	hoursStr, err := reader.ReadString('\n')
	hoursStr = strings.ToLower(strings.TrimSpace(hoursStr))
	if err != nil {
		return fish, fmt.Errorf("failed reading hours: %w", err)
	}
	if hoursStr == "all" {
		fish.Hours = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	} else {
		h, err := rangeString(hoursStr, 24)
		if err != nil {
			return fish, err
		}
		fish.Hours = h
	}

	fmt.Print("...during months (e.g., 'March-September', 'aug-oct', 'Sept-apr', or 'all' if all)? ")
	monthsStr, err := reader.ReadString('\n')
	monthsStr = strings.ToLower(strings.TrimSpace(monthsStr))
	if err != nil {
		return fish, fmt.Errorf("failed reading months: %w", err)
	}
	if monthsStr == "all" {
		fish.Months = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	} else {
		fish.Months, err = rangeMonths(monthsStr)
		if err != nil {
			return fish, ErrBadInput
		}
	}

	fmt.Print("...shadow size? ")
	shadowSize, err := reader.ReadString('\n')
	if err != nil {
		return fish, fmt.Errorf("failed reading shadow size: %w", err)
	}
	fish.ShadowSize = strings.TrimSpace(shadowSize)

	return fish, nil
}

func getBug() (Bug, error) {
	var bug Bug

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Bugs's name? ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return bug, fmt.Errorf("failed reading name: %w", err)
	}
	bug.Name = strings.TrimSpace(name)

	fmt.Print("...price? ")
	price, err := reader.ReadString('\n')
	if err != nil {
		return bug, fmt.Errorf("failed reading price: %w", err)
	}
	p, err := strconv.Atoi(strings.TrimSpace(price))
	if err != nil {
		return bug, ErrBadInput
	}
	bug.Price = p

	fmt.Print("...during months (e.g., 'March-September', 'aug-oct', 'Sept-apr', or 'all' if all)? ")
	monthsStr, err := reader.ReadString('\n')
	monthsStr = strings.ToLower(strings.TrimSpace(monthsStr))
	if err != nil {
		return bug, fmt.Errorf("failed reading months: %w", err)
	}
	if monthsStr == "all" {
		bug.Months = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	} else {
		bug.Months, err = rangeMonths(monthsStr)
		if err != nil {
			return bug, ErrBadInput
		}
	}

	fmt.Print("...hours (24 hour ints in range, e.g. 8-17 for 8am to 5pm, or 'all'? ")
	hoursStr, err := reader.ReadString('\n')
	hoursStr = strings.ToLower(strings.TrimSpace(hoursStr))
	if err != nil {
		return bug, fmt.Errorf("failed reading hours: %w", err)
	}
	if hoursStr == "all" {
		bug.Hours = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	} else {
		h, err := rangeString(hoursStr, 24)
		if err != nil {
			return bug, err
		}
		bug.Hours = h
	}

	fmt.Print("...location? ")
	location, err := reader.ReadString('\n')
	if err != nil {
		return bug, fmt.Errorf("failed reading location: %w", err)
	}
	bug.Location = strings.TrimSpace(location)

	return bug, nil
}

func saveJSON(acnh ACNH) error {
	f, err := os.Create("acnh.json")
	if err != nil {
		return fmt.Errorf("unable to open acnh.json for writing: %w", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(acnh)
	if err != nil {
		return fmt.Errorf("unable to encode acnh.json: %w", err)
	}

	return nil
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
