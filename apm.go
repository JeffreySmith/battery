package battery

import (
	"fmt"
	"regexp"
)
//These Values come from the apm man page
const (
	High     BatteryStatus = 0
	Low      BatteryStatus = 1
	Critical BatteryStatus = 2
	Charging BatteryStatus = 3
	Absent   BatteryStatus = 4
	Unknown  BatteryStatus = 255
)

var apmOutput = regexp.MustCompile("([0-9]+)%")
var minuteOutput = regexp.MustCompile("([0-9]+) minutes")
var state = regexp.MustCompile("Battery state: ([a-z]+)")


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
