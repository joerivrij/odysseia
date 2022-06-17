package app

import (
	"github.com/odysseia/aristoteles/configs"
)

type PeriklesHandler struct {
	Config *configs.PeriklesConfig
}

func (p *PeriklesHandler) Flow() {
	wait := make(chan struct{})
	for {
		go p.CheckForAnnotations(wait)
		<-wait
	}
}
