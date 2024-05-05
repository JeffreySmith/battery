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

package battery_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/JeffreySmith/battery"
	"github.com/google/go-cmp/cmp"
)

func TestApmBatteryStatus(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("/usr/sbin/apm"); err != nil {
		t.Skipf("Unable to run 'apm' command, skipping: %v", err)
	}
	_, err := battery.ApmBatteryStat()
	if err != nil {
		t.Error(err)
	}
}
func TestBatteryLifeParse(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("testdata/apm.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := battery.Battery{Hours: 9, Minutes: 54}
	got := battery.Battery{}
	err = got.ParseApmBatteryLife(string(data))
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
func TestBatteryLifeUnknown(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("testdata/apm_unknown.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := battery.Battery{Minutes: -1}
	got := battery.Battery{}
	err = got.ParseApmBatteryLife(string(data))
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
func TestBatteryLifeParseBadInput(t *testing.T) {
	t.Parallel()
	got := battery.Battery{}
	err := got.ParseApmBatteryLife("bad input")
	if err == nil {
		t.Error(err)
	}
}
func TestApmCommandOutput(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("/usr/sbin/apm"); err != nil {
		t.Skipf("Unable to run 'apm' command, skipping: %v", err)
	}
	data, err := battery.GetApmOutput("/usr/sbin/apm")
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(data, "Battery state") {
		t.Skipf("No battery detected")
	}
	_, err = battery.GetApmOutput("/usr/sbin/apm")
	if err != nil {
		t.Fatal(err)
	}
}
func TestFailedApm(t *testing.T) {
	t.Parallel()
	_, err := battery.GetApmOutput("/usr/sbin/notapm")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
func TestApmCharging(t *testing.T) {
	t.Parallel()
	input, err := os.ReadFile("testdata/apm_charging.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := battery.Battery{Charging: true}
	got := battery.Battery{}
	err = got.ParseApmCharging(string(input))
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
func TestApmNotCharging(t *testing.T) {
	t.Parallel()
	input, err := os.ReadFile("testdata/apm.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := battery.Battery{Charging: false}
	got := battery.Battery{Charging: true}
	err = got.ParseApmCharging(string(input))
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
func TestApmNoCharging(t *testing.T) {
	t.Parallel()
	got := battery.Battery{}
	err := got.ParseApmCharging("bad input")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestApmBatteryPercent(t *testing.T) {
	t.Parallel()
	input, err := os.ReadFile("testdata/apm.txt")
	want := battery.Battery{ChargePercent: 90}
	got := battery.Battery{}
	err = got.ParseApmBatteryPercent(string(input))
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
func TestApmNoBatteryPercent(t *testing.T) {
	t.Parallel()
	got := battery.Battery{}
	err := got.ParseApmBatteryPercent("abc")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestApmInputToString(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("testdata/apm.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := battery.Battery{
		ChargePercent: 90,
		Hours:         9,
		Minutes:       54,
		Charging:      false,
	}
	got, err := battery.ParseApmOutput(string(data))
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
func TestApmInputToStringOnStruct(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("testdata/apm.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := battery.Battery{
		ChargePercent: 90,
		Hours:         9,
		Minutes:       54,
		Charging:      false,
	}
	got := battery.Battery{}
	err = got.ParseApmOutput(string(data))
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
func TestApmWithNoBatteryInput(t *testing.T) {
	t.Parallel()
	f, err := os.ReadFile("testdata/apm_nobattery.txt")
	_, err = battery.ParseApmOutput(string(f))
	if err == nil {
		t.Error(err)
	}
}
