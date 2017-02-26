package constructor

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Constructor-Engine
var (
	constructorEngineHeight = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "dimwit_constructor_engine_current_height",
		Help: "Current height the constructor engine is at",
	})
)

func InitPrometheus() {
	// Constructor-Engine
	prometheus.MustRegister(constructorEngineHeight)
}
