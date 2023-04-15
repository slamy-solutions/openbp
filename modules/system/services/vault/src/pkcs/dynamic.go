package pkcs

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"

	"github.com/miekg/pkcs11"
)

/*
	Dynamic PKCS handler that can load other PKCS libraries (based on softhsm2 compatibility)
*/
type DynamicPKCSHandle struct {
	PKCS11Ctx *pkcs11.Ctx

	libraryPath string
	slot        uint
	loggedIn    bool
	session     pkcs11.SessionHandle

	lock sync.Mutex

	closed bool
}

func NewDynamicPKCSHandle(libraryPath string, slot uint) *DynamicPKCSHandle {
	return &DynamicPKCSHandle{
		PKCS11Ctx:   nil,
		libraryPath: libraryPath,
		slot:        slot,
		loggedIn:    false,
		lock:        sync.Mutex{},
		closed:      false,
	}
}

func (h *DynamicPKCSHandle) Initialize() error {
	h.lock.Lock()
	defer h.lock.Unlock()

	p := pkcs11.New(h.libraryPath)
	err := p.Initialize()
	if err != nil {
		return errors.New("failed to initialize PKCS library: " + err.Error())
	}
	h.PKCS11Ctx = p

	// Ensure that slot exists
	if _, err := p.GetSlotInfo(h.slot); err != nil {
		return errors.New("failed to get PKCS slot information. Maybe it is wrong slot. Error: " + err.Error())
	}

	// Ensure that token exists in slot
	if _, err = p.GetTokenInfo(h.slot); err != nil {
		return errors.New("failed to get PKCS token information. Maybe token is not initialized in the slot. Error: " + err.Error())
	}

	return nil
}

func (h *DynamicPKCSHandle) GetProviderName() string {
	return "dynamic"
}

func (h *DynamicPKCSHandle) EnsureSessionAndLogIn(password string) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.loggedIn {
		h.PKCS11Ctx.CloseSession(h.session)
		h.loggedIn = false
	}

	// Open session to the PKCS token
	s, err := h.PKCS11Ctx.OpenSession(h.slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		return errors.New("failed to open PKCS slot session: " + err.Error())
	}
	h.session = s

	err = h.PKCS11Ctx.Login(h.session, pkcs11.CKU_USER, password)
	if err != nil {
		h.PKCS11Ctx.CloseSession(s)
	}

	h.loggedIn = true
	return nil
}
func (h *DynamicPKCSHandle) LogOutAndCloseSession() error {
	h.lock.Lock()
	defer h.lock.Unlock()

	errString := ""
	if closeSessionErr := h.PKCS11Ctx.CloseSession(h.session); closeSessionErr != nil {
		errString = errString + "error while closing session: " + closeSessionErr.Error()
	}
	if logoutErr := h.PKCS11Ctx.Logout(h.session); logoutErr != nil {
		errString = errString + "error while logging out: " + logoutErr.Error()
	}
	h.session = 0
	h.loggedIn = false

	if errString != "" {
		return errors.New(errString)
	}

	return nil
}

func (h *DynamicPKCSHandle) EnsureRSAKeyPair(ctx context.Context, name string) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	//TODO: ensure this is safe between multiple instances of the vault

	if !h.loggedIn {
		return ErrPKCSNotLoggedIn
	}

	searchTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, name),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
	}
	err := h.PKCS11Ctx.FindObjectsInit(h.session, searchTemplate)
	if err != nil {
		go h.LogOutAndCloseSession()
		return errors.New("failed to initialize search for RSA key-pair object in PKCS: " + err.Error())
	}
	defer h.PKCS11Ctx.FindObjectsFinal(h.session)

	objs, _, err := h.PKCS11Ctx.FindObjects(h.session, 1)
	if err != nil {
		go h.LogOutAndCloseSession()
		return errors.New("failed to search for RSA key-pair object in PKCS: " + err.Error())
	}

	// There is no RSA key-pair
	if len(objs) == 0 {
		mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_KEY_PAIR_GEN, nil)}
		publicKeyTemplate := []*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
			pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
			pkcs11.NewAttribute(pkcs11.CKA_LABEL, name),
			pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
			pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
			pkcs11.NewAttribute(pkcs11.CKA_MODULUS_BITS, 2048),
		}
		privateKeyTemplate := []*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
			pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
			pkcs11.NewAttribute(pkcs11.CKA_LABEL, name),
			pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
			pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
			pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),
			pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
			pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		}

		_, _, err := h.PKCS11Ctx.GenerateKeyPair(h.session, mechanism, publicKeyTemplate, privateKeyTemplate)
		if err != nil {
			go h.LogOutAndCloseSession()
			return errors.New("error while generating RSA key-pair in the PKCS: " + err.Error())
		}
	}

	return nil
}

/*func (h *DynamicPKCSHandle) GetRSAPublicKey(ctx context.Context, name string) ([]byte, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.PKCS11Ctx.obje
}*/
func (h *DynamicPKCSHandle) SignRSA(ctx context.Context, name string, message io.Reader) ([]byte, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	privateKey, _, err := h.findRSAKeyPair(name)
	if err != nil {
		return []byte{}, err
	}

	mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS, nil)}
	err = h.PKCS11Ctx.SignInit(h.session, mechanism, privateKey)
	if err != nil {
		go h.LogOutAndCloseSession()
		return []byte{}, errors.New("error while initializing RSA signing algorithm: " + err.Error())
	}

	dataBuffer := make([]byte, 2048)
	for {
		readed, err := message.Read(dataBuffer)

		if readed != 0 {
			if updateErr := h.PKCS11Ctx.SignUpdate(h.session, dataBuffer[:readed]); updateErr != nil {
				h.PKCS11Ctx.SignFinal(h.session)
				go h.LogOutAndCloseSession()
				return []byte{}, errors.New("error while signing block of data: " + err.Error())
			}
		}

		if err != nil {
			if err == io.EOF {
				signature, err := h.PKCS11Ctx.SignFinal(h.session)

				if err != nil {
					go h.LogOutAndCloseSession()
					return []byte{}, errors.New("error while finalizing signature generation: " + err.Error())
				}

				return signature, nil
			}

			h.PKCS11Ctx.SignFinal(h.session)
			return []byte{}, errors.New("error while reading data for signing: " + err.Error())
		}
	}
}
func (h *DynamicPKCSHandle) VerifyRSA(ctx context.Context, name string, message io.Reader, signature []byte) (bool, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	_, publicKey, err := h.findRSAKeyPair(name)
	if err != nil {
		return false, err
	}

	mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS, nil)}
	err = h.PKCS11Ctx.VerifyInit(h.session, mechanism, publicKey)
	if err != nil {
		go h.LogOutAndCloseSession()
		return false, errors.New("error while initializing RSA verification algorithm: " + err.Error())
	}

	dataBuffer := make([]byte, 2048)
	for {
		readed, err := message.Read(dataBuffer)

		if readed != 0 {
			if updateErr := h.PKCS11Ctx.SignUpdate(h.session, dataBuffer[:readed]); updateErr != nil {
				h.PKCS11Ctx.VerifyFinal(h.session, []byte{})
				go h.LogOutAndCloseSession()
				return false, errors.New("error while verifiying block of data: " + err.Error())
			}
		}

		if err != nil {
			if err == io.EOF {
				err := h.PKCS11Ctx.VerifyFinal(h.session, signature)
				return err == nil, nil
			}

			h.PKCS11Ctx.VerifyFinal(h.session, []byte{})
			return false, errors.New("error while reading data for verification: " + err.Error())
		}
	}
}

// Tries to find key-pair. returns private and public key handles
func (h *DynamicPKCSHandle) findRSAKeyPair(name string) (pkcs11.ObjectHandle, pkcs11.ObjectHandle, error) {
	findTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, name),
	}

	if err := h.PKCS11Ctx.FindObjectsInit(h.session, findTemplate); err != nil {
		go h.LogOutAndCloseSession()
		return 0, 0, errors.New("error while initializing PKCS search of RSA keys: " + err.Error())
	}
	defer h.PKCS11Ctx.FindObjectsFinal(h.session)

	objs, _, err := h.PKCS11Ctx.FindObjects(h.session, 2)
	if err != nil {
		go h.LogOutAndCloseSession()
		return 0, 0, errors.New("error while performing PKCS search of RSA keys: " + err.Error())
	}

	if len(objs) == 2 {
		return 0, 0, ErrRSAKeyDoesntExist
	}

	var rsaPublicKeyHandle pkcs11.ObjectHandle
	var rsaPrivateKeyHandle pkcs11.ObjectHandle
	pubKeyFounded := false
	privKeyFounded := false

	pubKeyAttribute := pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY)
	privKeyAttribute := pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY)

	for _, obj := range objs {
		attributes := []*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_CLASS, name),
		}

		attributes, err := h.PKCS11Ctx.GetAttributeValue(h.session, obj, attributes)
		if err != nil {
			go h.LogOutAndCloseSession()
			return 0, 0, errors.New("error while performing PKCS search of RSA keys and getting object attributes: " + err.Error())
		}
		if bytes.Equal(attributes[0].Value, pubKeyAttribute.Value) {
			pubKeyFounded = true
			rsaPublicKeyHandle = obj
		} else if bytes.Equal(attributes[0].Value, privKeyAttribute.Value) {
			privKeyFounded = false
			rsaPrivateKeyHandle = obj
		}
	}

	if !(pubKeyFounded && privKeyFounded) {
		return 0, 0, errors.New("problems with RSA keys with name [" + name + "]. Probably duplicated object labels in the PKCS token.")
	}

	return rsaPrivateKeyHandle, rsaPublicKeyHandle, nil
}

func (h *DynamicPKCSHandle) Close() error {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.PKCS11Ctx.Destroy()
	h.closed = true
	return nil
}
