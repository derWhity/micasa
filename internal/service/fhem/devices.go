package fhem

import "time"

// The Device interface represents a device registered inside the FHEM instance
type Device interface {
	Type() string
}

// The Switch interface describes a switch device which can be switched on and off
type Switch interface {
	// On switches the device on
	On() error
	// Off switches the device off
	Off() error
	// IsOn checks if the device is switched on
	IsOn() (bool, error)
}

// The ContactSensor interface describes a sensor that checks if a door or window has been
// opened or closed
type ContactSensor interface {
	// IsOpen checks if the sensor reports that the door or window is currently open
	IsOpen() (bool, error)
}

// The TemperatureSensor interface describes a temperature sensor
type TemperatureSensor interface {
	// GetTemperature reads the current temperature measured by the sensor
	GetTemperature() (*Temperature, error)
}

// The HumiditySensor interface describes a humidity sensor that can measure relative humidity
type HumiditySensor interface {
	// GetHumidity reads the current relative humidity measured by the sensor
	GetHumidity() (float64, error)
}

// The ClimateControl interface describes a device that can change the temperature inside a room like a
// heater or AC
type ClimateControl interface {
	TemperatureSensor

	SetTemperature(t Temperature) error
}

// The ClimatePlanner interface describes a device that can set-up one or more climate plans
type ClimatePlanner interface {
	NumPlans() (uint, error)
	GetPlan(number uint) (ClimatePlan, error)
	SetPlan(number uint, plan ClimatePlan) error
}

// A ClimatePlan stores climate settings for every day of the week
type ClimatePlan interface {
	GetFor(day time.Weekday) (DayPlan, error)
	SetFor(day time.Weekday, plan DayPlan) error
}

// A DayPlan represents the planned temperature changes over one day. It needs to be able to return a series of
// TemperatureEntries that contain all temperature changes for this day
type DayPlan interface {
	GetSeries() []TemperatureEntry
}

// A TemperatureEntry describes the change to a specific temperature starting at a specific time of the day
type TemperatureEntry struct {
	// Time the new temperature starts. The date part of the timestamp is ignored
	StartTime time.Time
	// The temperature to set when the start time comes
	Value Temperature
}
