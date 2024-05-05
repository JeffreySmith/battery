package battery

import (
	"errors"
	"fmt"
	"math"
	"os/exec"
	"regexp"
	"strconv"
)

var batteryPercent = regexp.MustCompile("([0-9]+)%")
var minuteOutput = regexp.MustCompile("([0-9]+) minutes")
var state = regexp.MustCompile("Battery state: ([a-z]+)")

//These Values come from the apm man page
const (
	High     BatteryStatus = 0
	Low      BatteryStatus = 1
	Critical BatteryStatus = 2
	Charging BatteryStatus = 3
	Absent   BatteryStatus = 4
	Unknown  BatteryStatus = 255
)

	
func (b BatteryStatus) String() string {
	switch b {
	case High:
		return "High"
	case Low:
		return "Low"
	case Critical:
		return "Critical"
	case Charging:
		return "Charging"
	case Absent:
		return "Absent"
	case Unknown:
		return "Unknown"
	default:
		return fmt.Sprintf("%d", b)
	}
}

func GetApmOutput(cmd string) (string, error) {
	data, err := exec.Command(cmd).CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(data), err
}

func (b *Battery)ParseApmBatteryLife(input string) error {
	var hours, minutes int
	var err error
	matches := minuteOutput.FindStringSubmatch(input)
	if len(matches) == 2 {
		minutes, err = strconv.Atoi(matches[1])
		if err != nil {
			return err
		}
		hours = int(minutes / 60)
		minutes = int(math.Mod(float64(minutes), 60))
		b.Hours = hours
		b.Minutes = minutes
		return nil
	} 
	return errors.New("Couldn't find battery life minutes")
}
func (b *Battery)ParseApmCharging(input string) error {
	matches := state.FindStringSubmatch(input)
	if len(matches) == 2 {
		if matches[1] == "charging" {
			b.Charging = true
		} else{
			b.Charging = false
		}
		return nil
	}
	return errors.New("No status found")
}

func (b *Battery)ParseApmBatteryPercent(input string) error {
	matches := batteryPercent.FindStringSubmatch(input)
	if len(matches) == 2 {
		val,err := strconv.Atoi(matches[1])
		if err != nil {
			return err
		}
		b.ChargePercent = val
		return nil
	}
	return errors.New("Couldn't find battery percent")
}
func (b *Battery)ParseApmOutput(input string) error{
	battery,err := ParseApmOutput(input)
	if err != nil {
		return err
	}
	*b = battery
	return nil
}
func ParseApmOutput(input string) (Battery, error){
	battery := Battery{}
	err := battery.ParseApmBatteryLife(input)
	if err != nil {
		return Battery{},err
	}
	err = battery.ParseApmBatteryPercent(input)
	if err != nil {
		return Battery{},err
	}
	err = battery.ParseApmCharging(input)
	if err != nil {
		return Battery{},err
	}
	return battery, nil
}
