/*BSD 3-Clause License

Copyright (c) 2024, Jeffrey Smith

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package battery

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var batteryPercent = regexp.MustCompile("([0-9]+)%")
var minuteOutput = regexp.MustCompile("([0-9]+) minutes")
var state = regexp.MustCompile("Battery state: ([a-z]+)")

// These Values come from the apm man page
const (
	High     BatteryStatus = 0
	Low      BatteryStatus = 1
	Critical BatteryStatus = 2
	Charging BatteryStatus = 3
	Absent   BatteryStatus = 4
	Unknown  BatteryStatus = 255
)

const (
	UnknownAdapter Adapter = 0
	Disconnected   Adapter = 1
	Connected      Adapter = 2
)

func (a Adapter) String() string {
	switch a {
	case UnknownAdapter:
		return "Unknown"
	case Disconnected:
		return "Not Connected"
	case Connected:
		return "Connected"
	default:
		return ""
	}
}
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

func ApmBatteryStat() (BatteryStatus, error) {
	data, err := exec.Command("/usr/sbin/apm", "-b").CombinedOutput()
	if err != nil {
		return Unknown, err
	}
	ret, err := strconv.Atoi(strings.Trim(string(data), "\n"))
	if err != nil {
		return Unknown, err
	}
	if ret == 0 || ret == 1 || ret == 2 || ret == 3 || ret == 4 || ret == 255 {
		return BatteryStatus(ret), nil
	} else {
		return Unknown, errors.New("Unexpected return value")
	}
}

func GetApmOutput(cmd string) (string, error) {
	data, err := exec.Command(cmd).CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(data), err
}
func (b *Battery) ParseBatteryState(input string) {
	matches := state.FindStringSubmatch(input)
	if len(matches) == 2 {
		switch matches[1] {
		case "charging":
			b.Battery = Charging
		case "absent":
			b.Battery = Absent
		case "high":
			b.Battery = High
		case "low":
			b.Battery = Low
		case "critical":
			b.Battery = Critical
		}
	} else {
		b.Battery = Unknown
	}
}
func (b *Battery) ParseAdapterStatus(input string) {
	var status Adapter
	if strings.Contains(input, "not connected") {
		status = Disconnected
	} else if strings.Contains(input, "connected") {
		status = Connected
	} else {
		status = UnknownAdapter
	}
	b.AdapterConnected = status
}

func (b *Battery) ParseApmBatteryLife(input string) error {
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
	} else if strings.Contains(input, "Battery state: absent") {
		b.Hours = 0
		b.Minutes = -1
		return nil
	} else if strings.Contains(input, "unknown life estimate") {
		b.Hours = 0
		b.Minutes = -1
		return nil
	} else if strings.Contains(input, "unknown recharge time") {
		b.Hours = 0
		b.Minutes = -1
		return nil
	}
	return errors.New("Unable to read remaining battery minutes")

}
func (b *Battery) ParseApmCharging(input string) error {
	matches := state.FindStringSubmatch(input)
	if len(matches) == 2 {
		if matches[1] == "charging" {
			b.Charging = true
		} else {
			b.Charging = false
		}
		return nil
	}
	return errors.New("No status found")
}

func (b *Battery) ParseApmBatteryPercent(input string) error {
	matches := batteryPercent.FindStringSubmatch(input)
	if len(matches) == 2 {
		val, err := strconv.Atoi(matches[1])
		if err != nil {
			return err
		}
		b.ChargePercent = val
		return nil
	}
	return errors.New("Couldn't find battery percent")
}
func (b *Battery) ToJSON() (string,error) {
	data, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func (b *Battery) ParseApmOutput(input string) error {
	battery, err := ParseApmOutput(input)
	if err != nil {
		return err
	}
	*b = battery
	return nil
}
func ParseApmOutput(input string) (Battery, error) {
	battery := Battery{}
	err := battery.ParseApmBatteryLife(input)
	if err != nil {
		return Battery{}, err
	}
	err = battery.ParseApmBatteryPercent(input)
	if err != nil {
		return Battery{}, err
	}
	err = battery.ParseApmCharging(input)
	if err != nil {
		return Battery{}, err
	}
	battery.ParseAdapterStatus(input)
	battery.ParseBatteryState(input)
	return battery, nil
}
func (b *Battery) PrintTimeRemaining() {
	if b.Minutes < 0 {
		fmt.Println("Estimated remaining time: Unknown")
	} else if b.Minutes < 10 {
		fmt.Printf("Estimated remaining time: %dh0%dm\n", b.Hours, b.Minutes)
	} else {
		fmt.Printf("Estimated remaining time: %dh%dm\n", b.Hours, b.Minutes)
	}

}
func OpenBSDMain() int {

	OS := runtime.GOOS
	if OS != "openbsd" {
		return 2
	}
	timeRemaining := flag.Bool("t", true, "Show time remaining")
	chargeStatus := flag.Bool("s", false, "Show whether the computer is charging")
	adapterStatus := flag.Bool("a", false, "Show the status of the adapter")
	percentRemaining := flag.Bool("p", false, "Show estimated remaining battery percent")
	jsonOutput := flag.Bool("j", false, "Display output in JSON")
	flag.Parse()
	apm_output, err := GetApmOutput("/usr/sbin/apm")
	battery, err := ParseApmOutput(apm_output)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if *jsonOutput {
		j,err := battery.ToJSON()
		if err != nil {
			return 1
		}
		fmt.Println(j)
		return 0
	}
	if *chargeStatus {
		if battery.Charging {
			fmt.Printf("Status: Charging\n")
		} else {
			fmt.Printf("Status: Not Charging\n")
		}
	}
	if *adapterStatus {
		fmt.Printf("Adapter Status: %v\n", battery.AdapterConnected)
	}
	if *percentRemaining {
		fmt.Printf("Remaining battery: %d%%\n", battery.ChargePercent)
	}
	if *timeRemaining {
		battery.PrintTimeRemaining()
	}
	return 0
}

