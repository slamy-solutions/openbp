package pkcs

import (
	"errors"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var _default_seal_password = "12345678"

type Sealer struct {
	savedPassword string

	pkcsProvider PKCS

	lock        sync.Mutex
	updatesChan chan struct{}
}

func NewSealer(pkcsProvider PKCS) *Sealer {
	return &Sealer{
		pkcsProvider:  pkcsProvider,
		savedPassword: "",
		updatesChan:   make(chan struct{}),
		lock:          sync.Mutex{},
	}
}

func (s *Sealer) Start() {
	s.lock.Lock()
	defer s.lock.Unlock()

	go func() {
	mainloop:
		for {
			select {
			case _, ok := <-s.updatesChan:
				if !ok {
					break mainloop
				}
			case <-time.After(time.Hour):
			}

			// Reload
			s.lock.Lock()
			defer s.lock.Unlock()

			if s.savedPassword != "" {
				err := s.pkcsProvider.EnsureSessionAndLogIn(s.savedPassword)
				if err != nil {
					if err == ErrPKCSBadLoginPassword {
						s.savedPassword = ""
						log.Info("[PKCS sealer] Tried unseal with saved password. Bad password. Forgetting it.")
					}
					log.Error("[PKCS sealer] Error while unsealing pkcs: " + err.Error())
				} else {
					log.Info("[PKCS sealer] Successfully created/recreated pkcs session and logged in.")
				}
			}
		}

		log.Info("[PKCS sealer] Successfully closed update worker.")
	}()

	log.Info("[PKCS sealer] Started worker. Submitting default password for unsealing.")

	s.savedPassword = _default_seal_password
	s.updatesChan <- struct{}{}
}

func (s *Sealer) Unseal(newPassword string) (bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	//hash the password to make it consistent length and not to save original password in memory
	//sha := sha256.Sum256([]byte(newPassword))
	//passwordHash := fmt.Sprintf("%x", sha)

	err := s.pkcsProvider.EnsureSessionAndLogIn(s.savedPassword)
	if err != nil {
		if err == ErrPKCSBadLoginPassword {
			return false, nil
		}
		log.Error("[PKCS sealer] Error while unsealing pkcs: " + err.Error())
		return false, errors.New("error while ensuring session and loggin in to the PKCS: " + err.Error())
	}

	log.Info("[PKCS sealer] Successfully created/recreated pkcs session and logged in.")
	s.savedPassword = newPassword
	return true, nil
}

func (s *Sealer) Seal() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.savedPassword = ""
	err := s.pkcsProvider.LogOutAndCloseSession()
	if err != nil {
		log.Error("[PKCS sealer] Error while sealing: " + err.Error())
	}
}

func (s *Sealer) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()

	log.Info("[PKCS sealer] Raised closing event.")

	s.savedPassword = ""
	close(s.updatesChan)
}
