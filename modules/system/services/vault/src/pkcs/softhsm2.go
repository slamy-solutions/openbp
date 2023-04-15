package pkcs

type SoftHSM2PKCSHandle struct {
	DynamicPKCSHandle
}

func NewSoftHSM2PKCSHandle(libraryPath string, slot uint) *SoftHSM2PKCSHandle {
	return &SoftHSM2PKCSHandle{
		DynamicPKCSHandle: *NewDynamicPKCSHandle(libraryPath, slot),
	}
}

func (h *SoftHSM2PKCSHandle) GetProviderName() string {
	return "softhsm2"
}
