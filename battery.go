package battery

type BatteryStatus int
type Mode int

type Battery struct {
	ChargePercent int
	Battery       BatteryStatus
	Hours         int
	Minutes       int
	Charging      bool
	PerformanceMode Mode
}
