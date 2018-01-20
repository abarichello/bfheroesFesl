package level

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var useRedis = false

type Level struct {
	DB     *leveldb.DB
	Legacy *LegacyRedis
}

func New(path string, rc *redis.Client) (*Level, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	// lg := newLegacyRedis(rc)
	return &Level{db, nil}, nil
}

func (lvl *Level) Get(ident string, key string) string {
	if useRedis {
		return lvl.Legacy.Get(ident, key)
	}

	val, _ := lvl.DB.Get(
		[]byte(fmt.Sprintf("%s/%s", ident, key)),
		nil,
	)
	return string(val)
}

func (lvl *Level) Set(ident string, key, value string) error {

	if useRedis {
		return lvl.Legacy.Set(ident, key, value)
	}

	err := lvl.DB.Put(
		[]byte(fmt.Sprintf("%s/%s", ident, key)),
		[]byte(value),
		nil,
	)
	return err

}

func (lvl *Level) SetM(ident string, hash map[string]interface{}) error {
	if useRedis {
		return lvl.Legacy.SetM(ident, hash)
	}

	tx, _ := lvl.DB.OpenTransaction()
	for k, v := range hash {
		tx.Put(
			[]byte(fmt.Sprintf("%s/%s", ident, k)),
			[]byte(v.(string)),
			nil,
		)
	}
	tx.Commit()
	return nil
}

type State struct {
	lvl        *Level
	identifier string
}

func (lvl *Level) NewState(identifier string) *State {
	logrus.WithField("op", "NewState").Debugf("redisState:%s", identifier)
	return &State{lvl, fmt.Sprintf("redisState:%s", identifier)}
}

func (st *State) Get(key string) string {
	logrus.WithFields(logrus.Fields{"op": "Get", "meth": "State"}).Debug(key)
	return st.lvl.Get(st.identifier, key)
}

func (st *State) Set(key, value string) error {
	logrus.WithFields(logrus.Fields{"op": "Set", "meth": "State"}).Debug(key, ":", value)

	return st.lvl.Set(st.identifier, key, value)
}

func (st *State) SetM(hash map[string]interface{}) error {
	for k, v := range hash {
		logrus.WithFields(logrus.Fields{"op": "SetM", "meth": "State"}).Debug(k, ":", v)
	}

	return st.lvl.SetM(st.identifier, hash)
}

func (st *State) Delete() error {
	logrus.WithFields(logrus.Fields{"op": "Delete", "meth": "State"}).Debug(st.identifier)

	if useRedis {
		return st.lvl.Legacy.Delete(st.identifier)
	}

	err := st.lvl.DB.Delete([]byte(st.identifier), nil)
	return err
}

type Object struct {
	lvl        *Level
	identifier string
}

func (lvl *Level) NewObject(prefix, ident string) *Object {
	logrus.WithField("op", "NewObject").Debugf("%s:%s", prefix, ident)
	return &Object{lvl, fmt.Sprintf("%s:%s", prefix, ident)}
}

func (obj *Object) Get(key string) string {
	logrus.WithFields(logrus.Fields{"op": "Get", "meth": "Object"}).Debug(key)
	return obj.lvl.Get(obj.identifier, key)
}

func (obj *Object) Set(key, value string) error {
	logrus.WithFields(logrus.Fields{"op": "Set", "meth": "Object"}).Debug(key, ":", value)
	return obj.lvl.Set(obj.identifier, key, value)
}

func (obj *Object) HKeys() []string {
	if useRedis {
		ks := obj.lvl.Legacy.HKeys(obj.identifier)
		logrus.WithFields(logrus.Fields{"op": "HKeys", "meth": "Object"}).Debug(ks)
		return ks
	}

	keys := []string{}
	it := obj.lvl.DB.NewIterator(&util.Range{
		Start: []byte(obj.identifier),
		Limit: []byte(obj.identifier + "/zzz"),
	}, nil)
	for it.Next() {
		keys = append(keys, string(it.Key()))
	}
	logrus.WithFields(logrus.Fields{"op": "HKeys", "meth": "Object"}).Debug(keys)
	return keys
}

func (obj *Object) Delete() error {
	logrus.WithFields(logrus.Fields{"op": "Delete", "meth": "Object"}).Debug(obj.identifier)

	if useRedis {
		return obj.lvl.Legacy.Delete(obj.identifier)
	}

	err := obj.lvl.DB.Delete([]byte(obj.identifier), nil)
	return err
}
