package ultrasonic

import (
	"machine"
	"time"
)

const timeout = 23324 // 4m

// Ultrasonic holds the pins
type Ultrasonic struct {
	trigger machine.Pin
	echo    machine.Pin
}

// New creates a new ultrasonic sensor
func New(trigger, echo machine.Pin) Ultrasonic {
	return Ultrasonic{
		trigger: trigger,
		echo:    echo,
	}
}

// Configure will setup the ultrasonic sensor
func (u *Ultrasonic) Configure() {
	u.trigger.Configure(machine.PinConfig{Mode: machine.PinOutput})
	u.echo.Configure(machine.PinConfig{Mode: machine.PinInput})
}

// ReadDistance reads the distance and returns it in mm
func (u *Ultrasonic) ReadDistance() int32 {
	println("Hitting ReadDistance()")
	pulse := u.ReadPulse()

	return (pulse * 1715) / 10000 // mm
}

// ReadPulse will read the pulse of the ultrasonic sensor
func (u *Ultrasonic) ReadPulse() int32 {
	println("Hitting ReadPulse()")
	t := time.Now()
	println(t.String())
	u.trigger.Low()
	time.Sleep(2 * time.Microsecond)
	u.trigger.High()
	time.Sleep(10 * time.Microsecond)
	u.trigger.Low()
	i := uint8(0)
	for {
		if u.echo.Get() {
			t = time.Now()
			break
		}
		i++
		if i > 10 {
			if time.Since(t).Microseconds() > timeout {
				println("Distance timeout")
				return 0
			}
			i = 0
		}
	}
	i = 0
	for {
		if !u.echo.Get() {
			return int32(time.Since(t).Microseconds())
		}
		i++
		if i > 10 {
			if time.Since(t).Microseconds() > timeout {
				println("Distance timeout")
				return 0
			}
			i = 0
		}
	}
}
