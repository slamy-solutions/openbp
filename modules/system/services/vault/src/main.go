package main

import (
	"os"

	"github.com/slamy-solutions/openbp/modules/system/services/vault/src/pkcs"

	log "github.com/sirupsen/logrus"
)

func main() {
	pkcsCtx, err := pkcs.NewPKCSFromEnv()
	if err != nil {
		log.Panic("Failed to load PKCS from env: " + err.Error())
		os.Exit(-1)
	}
	log.Info("Using [" + pkcsCtx.GetProviderName() + "] HSM provider.")

	err = pkcsCtx.Initialize()
	if err != nil {
		log.Panic("Failed to initialize PKCS: " + err.Error())
		os.Exit(-2)
	}
	defer pkcsCtx.Close()

	sealer := pkcs.NewSealer(pkcsCtx)
	sealer.Start()
	defer sealer.Stop()

}
