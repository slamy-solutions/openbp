package pkcs

import (
	"errors"

	pkcs11 "github.com/miekg/pkcs11"
)

type SoftHSM2PKCSHandle struct {
	DynamicPKCSHandle
}

func findSlotBylabel(libraryPath string, tokenLabel string) (uint, error) {
	p := pkcs11.New(libraryPath)
	if p == nil {
		return 0, errors.New("failed to open PKCS library")
	}
	err := p.Initialize()
	if err != nil {
		return 0, errors.New("failed to initialize PKCS library: " + err.Error())
	}
	defer p.Destroy()
	defer p.Finalize()

	slots, err := p.GetSlotList(true)
	if err != nil {
		return 0, errors.New("failed to initialize PKCS library: " + err.Error())
	}

	for _, slot := range slots {
		tokenInfo, err := p.GetTokenInfo(slot)
		if err != nil {
			return 0, errors.New("error while getting token info: " + err.Error())
		}

		if tokenInfo.Label == tokenLabel {
			return slot, nil
		}
	}

	return 0, errors.New("slot with token and specified token label wasnt founded")
}

func NewSoftHSM2PKCSHandle(libraryPath string, tokenLabel string) (*SoftHSM2PKCSHandle, error) {
	// Search fo slot ID.\
	// failed to identify SOFTHSM slot:
	slot, err := findSlotBylabel(libraryPath, tokenLabel)
	if err != nil {
		return nil, errors.New("failed to identify SOFTHSM slot: " + err.Error())
	}

	return &SoftHSM2PKCSHandle{
		DynamicPKCSHandle: *NewDynamicPKCSHandle(libraryPath, slot),
	}, nil
}

func (h *SoftHSM2PKCSHandle) GetProviderName() string {
	return "softhsm2"
}
