package app

import (
	"crypto/rsa"
	"github.com/odysseia/aristoteles/configs"
)

type PeriklesHandler struct {
	Config     *configs.PeriklesConfig
	Ca         []byte
	PrivateKey *rsa.PrivateKey
}

func (p *PeriklesHandler) Flow() {
	wait := make(chan struct{})
	for {
		go p.CheckForAnnotations(wait)
		<-wait
	}
}
