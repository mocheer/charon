package ws

import (
	"time"

	"github.com/mocheer/pluto/js/window"
)

func push() {
	window.SetInterval(getTyphoon, time.Minute)
}
