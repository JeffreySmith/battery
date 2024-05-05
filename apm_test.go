package battery_test

import (
	
	"os"
	"strings"
	"testing"
	"github.com/JeffreySmith/battery"
	"github.com/google/go-cmp/cmp"
)
func TestBatteryLifeParse(t *testing.T){
	t.Parallel()
	data, err := os.ReadFile("testdata/apm.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := battery.Battery{Hours:9,Minutes:54}
	got := battery.Battery{}
	err = got.ParseApmBatteryLife(string(data))
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want,got))
	}
}
func TestBatteryLifeParseBadInput(t *testing.T){
	t.Parallel()
	got := battery.Battery{}
	err := got.ParseApmBatteryLife("bad input")
	if err == nil {
		t.Error(err)
	}
}
func TestApmCommandOutput(t *testing.T) {
	t.Parallel()
	data, err := battery.GetApmOutput("/usr/sbin/apm")
	if err != nil {
		t.Skipf("Unable to run 'apm' command: %v", err)
	}
	if !strings.Contains(data,"Battery state"){
		t.Skipf("No battery detected")
	}
	_ , err = battery.GetApmOutput("/usr/sbin/apm")
	if err != nil {
		t.Fatal(err)
	}
}
func TestFailedApm(t *testing.T) {
	t.Parallel()
	_, err := battery.GetApmOutput("/usr/bin/notapm")
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
	want := battery.Battery{Charging:true}
	got := battery.Battery{}
	err = got.ParseApmCharging(string(input))
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(want, got){
		t.Error(cmp.Diff(want, got))
	}
}
func TestApmNotCharging(t *testing.T) {
	t.Parallel()
	input, err := os.ReadFile("testdata/apm.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := battery.Battery{Charging:false}
	got := battery.Battery{Charging:true}
	err = got.ParseApmCharging(string(input))
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(want, got){
		t.Error(cmp.Diff(want, got))
	}
}
func TestApmNoCharging(t *testing.T){
	t.Parallel()
	got := battery.Battery{}
	err := got.ParseApmCharging("bad input")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestApmBatteryPercent(t *testing.T){
	t.Parallel()
	input, err := os.ReadFile("testdata/apm.txt")
	want := battery.Battery{ChargePercent: 90}
	got := battery.Battery{}
	err = got.ParseApmBatteryPercent(string(input))
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got){
		t.Error(cmp.Diff(want, got))
	}
}
func TestApmNoBatteryPercent(t *testing.T){
	t.Parallel()
	got := battery.Battery{}
	err := got.ParseApmBatteryPercent("abc")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestApmInputToString(t *testing.T){
	t.Parallel()
	data, err := os.ReadFile("testdata/apm.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := battery.Battery {
		ChargePercent: 90,
		Hours:9,
		Minutes: 54,
		Charging: false,
	}
	got, err := battery.ParseApmOutput(string(data))
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
func TestApmInputToStringOnStruct(t *testing.T){
	t.Parallel()
	data, err := os.ReadFile("testdata/apm.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := battery.Battery {
		ChargePercent: 90,
		Hours:9,
		Minutes: 54,
		Charging: false,
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
func TestApmWithNoBatteryInput(t *testing.T){
	t.Parallel()
	f, err := os.ReadFile("testdata/apm_nobattery.txt")
	_, err = battery.ParseApmOutput(string(f))
	if err == nil {
		t.Error(err)
	}
}
