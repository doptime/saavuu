package data

import (
	"math/rand"

	"github.com/bits-and-blooms/bloom/v3"
)

func (db *Ctx[k, v]) BuildKeysBloomFilter(capacity int, falsePosition float64) (err error) {
	var keys []string
	if keys, err = db.Rds.Keys(db.Ctx, db.Key).Result(); err != nil {
		return err
	}
	if capacity <= 0 || falsePosition <= 0 || falsePosition >= 1 {
		db.KeysBloom = bloom.NewWithEstimates(uint(len(keys))+uint(rand.Uint32()%1000+1000), 0.0000001+rand.Float64()/10000000)
	} else {
		db.KeysBloom = bloom.NewWithEstimates(uint(capacity), falsePosition)
	}
	//if type of k is string, then AddString is faster than Add
	for _, it := range keys {
		db.KeysBloom.AddString(it)
	}
	return nil
}
func (db *Ctx[k, v]) TestBloomKey(key k) (exist bool, err error) {
	var keyStr string
	if err = nil; db.KeysBloom == nil {
		err = db.BuildKeysBloomFilter(-1, -1.0)
	}
	if err != nil {
		return false, err
	}
	if keyStr, err = db.toKeyStr(key); err != nil {
		return false, err
	}
	return db.KeysBloom.TestString(keyStr), nil
}
func (db *Ctx[k, v]) AddBloomKey(key k) (err error) {
	var keyStr string
	if err = nil; db.KeysBloom == nil {
		err = db.BuildKeysBloomFilter(-1, -1.0)
	}
	if err != nil {
		return err
	}
	if keyStr, err = db.toKeyStr(key); err != nil {
		return err
	}
	db.KeysBloom.AddString(keyStr)
	return nil
}
