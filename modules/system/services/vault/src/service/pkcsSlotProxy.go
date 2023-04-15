package service

import (
	"errors"
	"time"

	"github.com/miekg/pkcs11"
)

var ErrSlotNotFound = errors.New("slot not found")

type PKCSSlotProxy struct {
	cache *slotCache

	PKCS11Ctx *pkcs11.Ctx
}

func NewPKCSSlotProxy(pkcsCtx *pkcs11.Ctx) *PKCSSlotProxy {
	return &PKCSSlotProxy{
		cache:     newSlotCache(time.Minute * 5),
		PKCS11Ctx: pkcsCtx,
	}
}

func (p *PKCSSlotProxy) GetSlotByName(name string) (uint, error) {
	// Check if slot is inside cache
	if slot, slotInCache := p.cache.Get(name); slotInCache {
		return slot, nil
	}

	slotList, err := p.PKCS11Ctx.GetSlotList(true)

	if err != nil {
		return 0, errors.New("failed to get pkcs11 slots list: " + err.Error())
	}

	for _, slot := range slotList {
		slotInfo, err := p.PKCS11Ctx.GetSlotInfo(slot)
		if err != nil {
			return 0, errors.New("failed to get pkcs11 slot info: " + err.Error())
			// Handle error
		}

		// Compare the slot name with the desired name
		if slotInfo.SlotDescription == name {
			p.cache.Set(name, slot, time.Minute)
			return slot, nil
		}
	}

	return 0, ErrSlotNotFound
}
