package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	pkcs11Ctx, err := InitializePKCS()
	if err != nil {
		log.Panic("Failed to initialize PKCS: " + err.Error())
		os.Exit(-1)
	}
	defer pkcs11Ctx.Destroy()
	defer pkcs11Ctx.Finalize()

	err = EnsureDefaultVaultMasterKey()
	if err != nil {
		log.Panic("Failed to ensure default vault master key: " + err.Error())
		os.Exit(-2)
	}

	var masterKeyLoaded bool = false
	vaultMasterKey, err := TryLoadVaultMasterKeyWithDefaultDecriptionSecret()
	if err != nil {
		if err != ErrFailedToDecryptMasterKey {
			log.Panic("Unknown error while trying to load vault master key using default decription secret: " + err.Error())
			os.Exit(-3)
		}
	} else {
		masterKeyLoaded = true
	}
}
