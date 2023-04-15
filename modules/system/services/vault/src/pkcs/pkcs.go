package pkcs

import (
	"crypto/aes"
	"crypto/rand"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/miekg/pkcs11"
	log "github.com/sirupsen/logrus"
)

// Default secret used to encrypt master key
var default_vault_decryption_secret = []byte("default")

// Total size of master key
const vault_master_key_size = 128

// Bigger the padding size - less probability of "false" decryption
const vault_master_key_padding_size = 32

// Slot wasnt found
var ErrSlotNotFound = errors.New("slot not found")

type PKCS11Handle struct {
	MasterKeyLock   sync.RWMutex
	MasterKey       []byte
	MasterKeyLoaded bool

	PKCS11Ctx *pkcs11.Ctx

	cache *slotCache
}

func InitializePKCS() (*PKCS11Handle, error) {
	libraryPath := os.Getenv("PKCS_LIBRARY_PATH")
	if os.Getenv("EMULATE_HSM") == "true" {
		libraryPath = os.Getenv("SOFTHSM_LIBRARY_PATH")
		log.Warn("Emulating PKCS11 with SoftHSM")
	}

	log.Info("Loading PKCS11 library from " + libraryPath)
	p := pkcs11.New(libraryPath)
	err := p.Initialize()
	if err != nil {
		return nil, errors.New("failed to initialize library: " + err.Error())
	}

	err = EnsureDefaultVaultMasterKey()
	if err != nil {
		return nil, errors.New("failed to ensure default master key: " + err.Error())
	}

	handle := &PKCS11Handle{
		MasterKeyLock:   sync.RWMutex{},
		MasterKey:       []byte{},
		MasterKeyLoaded: false,
		PKCS11Ctx:       p,
		cache:           newSlotCache(time.Minute),
	}

	handle.TryLoadVaultMasterKeyWithDefaultDecriptionSecret()

	return handle, nil
}

var ErrFailedToDecryptMasterKey = errors.New("failed to decrypt master key. Decryption key may not be valid or data is corrupted")

func EnsureDefaultVaultMasterKey() error {
	masterKeyFilePath := os.Getenv("MASTER_KEY_FILE_PATH")

	if _, err := os.Stat(masterKeyFilePath); err != nil {
		if os.IsNotExist(err) {
			// Generate new master key
			masterKey := make([]byte, vault_master_key_size)
			_, err = rand.Read(masterKey)
			if err != nil {
				return errors.New("failed to securely generate random master key: " + err.Error())
			}
			for i := 1; i <= vault_master_key_padding_size; i++ {
				masterKey[i] = 0
			}

			// Encrypt master key and save to the file
			c, err := aes.NewCipher(default_vault_decryption_secret)
			if err != nil {
				return errors.New("failed to create AES cipher: " + err.Error())
			}
			encryptedKey := make([]byte, len(masterKey))
			c.Encrypt(encryptedKey, masterKey)
			err = os.WriteFile(masterKeyFilePath, encryptedKey, 0)
			if err != nil {
				return errors.New("failed to write encrypted master key to the disk: " + err.Error())
			}
			log.Info("Created master key with default decryption key.")

			return nil
		} else {
			return errors.New("failed to read master key file: " + err.Error())
		}
	}

	return nil
}

func (p *PKCS11Handle) LoadVaultMasterKey(decryptionSecret []byte) error {
	masterKeyFilePath := os.Getenv("MASTER_KEY_FILE_PATH")
	var rawData []byte
	var err error

	// Load encrypted master key from the dics
	if rawData, err = os.ReadFile(masterKeyFilePath); err != nil {
		return errors.New("failed to read master key file: " + err.Error())
	}

	// Decrypt master key
	c, err := aes.NewCipher(decryptionSecret)
	if err != nil {
		return errors.New("failed to create AES cipher: " + err.Error())
	}
	result := make([]byte, len(rawData))
	c.Decrypt(result, rawData)

	// Check if decrypted master key is valid or not
	for i := 1; i <= vault_master_key_padding_size; i++ {
		if result[i] != 0 {
			return ErrFailedToDecryptMasterKey
		}
	}

	p.MasterKeyLock.Lock()
	defer p.MasterKeyLock.Unlock()

	p.MasterKeyLoaded = true
	p.MasterKey = result

	return nil
}

func (p *PKCS11Handle) TryLoadVaultMasterKeyWithDefaultDecriptionSecret() error {
	return p.LoadVaultMasterKey(default_vault_decryption_secret)
}

func (p *PKCS11Handle) GetSlotByName(name string) (uint, error) {
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
