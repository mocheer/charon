package ws

import (
	"time"

	"github.com/mocheer/pluto/js"
)

func push() {
	js.SetInterval(getTyphoon, time.Minute)
}
