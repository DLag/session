package session

import (
	"bytes"
	"encoding/gob"
	"sync"
	"time"

	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

const DefaultCookieName = "GOSESSION"
const DefaultIPStrict = false
const DefaultTTL = 24 * time.Hour
var DefaultCookie = http.Cookie{
	Name: DefaultCookieName,
	MaxAge: int(DefaultTTL.Seconds()),
	Path: "/",
}

type Manager struct {
	emptyCookie http.Cookie
	ipStrict   bool
	store      Store
	buf        bytes.Buffer
	enc        *gob.Encoder
	dec        *gob.Decoder
	sync.Mutex
}

func DefaultManager(store Store) *Manager {
	return NewManager(DefaultCookie, DefaultIPStrict, store)
}

func NewManager(emptyCookie http.Cookie, ipStrict bool, store Store) *Manager {
	m := &Manager{
		emptyCookie: emptyCookie,
		ipStrict:   ipStrict,
		store:      store,
	}
	m.enc = gob.NewEncoder(&m.buf)
	m.dec = gob.NewDecoder(&m.buf)
	return m
}

func (m *Manager) marshal(v interface{}) ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	m.buf.Reset()
	err := m.enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return m.buf.Bytes(), nil
}

func (m *Manager) unmarshal(data []byte, v interface{}) error {
	m.Lock()
	defer m.Unlock()
	m.buf.Reset()
	m.buf.Write(data)
	return m.dec.Decode(v)
}

func (m *Manager) get(session string) (map[string]interface{}, error) {
	buf, err := m.store.Get(session)
	if err != nil {
		return nil, errors.Wrap(err, "Can't read session data from store")
	}

	values := make(map[string]interface{})
	err = m.unmarshal(buf, &values)
	if err != nil {
		return nil, errors.Wrap(err, "Can't unmarshal session data given by store")
	}

	return values, nil
}

func (m *Manager) Get(session, key string) (interface{}, error) {
	values, err := m.get(session)
	if err != nil {
		return nil, err
	}

	res, _ := values[key]
	return res, nil
}

func (m *Manager) set(session string, v map[string]interface{}) error {
	buf, err := m.marshal(v)
	if err != nil {
		return errors.Wrap(err, "Can't marshal session object")
	}

	err = m.store.Set(session, buf, time.Duration(m.emptyCookie.MaxAge)*time.Second)
	if err != nil {
		return errors.Wrap(err, "Can't write session data to store")
	}
	return nil
}

func (m *Manager) Set(session, key string, v interface{}) error {
	values, err := m.get(session)
	if err != nil {
		return err
	}

	values[key] = v
	return m.set(session, values)
}

func (m *Manager) Delete(session string) error {
	err := m.store.Delete(session)
	if err != nil {
		return errors.Wrap(err, "Can't delete session from store")
	}
	return nil
}

func (m *Manager) Session(w http.ResponseWriter, r *http.Request) *Session {
	c, err := r.Cookie(m.emptyCookie.Name)
	var name string
	if err == nil {
		name = c.Value
	}
	if len(name) != 36 {
		hardening := r.UserAgent()
		if m.ipStrict {
			hardening += strings.SplitN(r.RemoteAddr, ":", 2)[0]
		}
		name = uuid.NewV5(uuid.NewV4(), hardening).String()
	}

	cookie := m.emptyCookie
	cookie.Value = name
	http.SetCookie(w, &cookie)

	s := &Session{
		name:    name,
		manager: m,
	}

	return s
}
