package gqserver

import (
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"
)

// TimeFn provides either the real current time or a fake time
type TimeFn func() time.Time

type StateManager interface {
	ParseConfig(string) error
	SetAESKey(string)
	PutUsedRandom([32]byte)
	DelUsedRandom([32]byte)
}

// State type stores the global state of the program
type State struct {
	WebServerAddr  string
	Key            string
	AESKey         []byte
	Now            TimeFn
	SS_LOCAL_HOST  string
	SS_LOCAL_PORT  string
	SS_REMOTE_HOST string
	SS_REMOTE_PORT string
	M              sync.RWMutex
	UsedRandom     map[[32]byte]int
}

// ParseConfig parses the config file into a State variable
func (sta *State) ParseConfig(configPath string) error {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &sta)
	if err != nil {
		return err
	}
	return nil
}

// SetAESKey calculates the SHA256 of the string key
func (sta *State) SetAESKey() {
	h := sha256.New()
	h.Write([]byte(sta.Key))
	sta.AESKey = h.Sum(nil)
}

func (sta *State) PutUsedRandom(random [32]byte) {
	sta.M.Lock()
	sta.UsedRandom[random] = int(sta.Now().Unix())
	sta.M.Unlock()
}

func (sta *State) DelUsedRandom(random [32]byte) {
	sta.M.Lock()
	delete(sta.UsedRandom, random)
	sta.M.Unlock()
}
