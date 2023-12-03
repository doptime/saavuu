package data

import (
	"log"
	"math/rand"

	"github.com/bits-and-blooms/bloom/v3"
)

func (db *Ctx[k, v]) BuildKeysBloomFilter(capacity int, falsePosition float64) (err error) {
	var keys []string
	if keys, err = db.Rds.Keys(db.Ctx, db.Key).Result(); err != nil {
		return err
	}
	if capacity <= 0 || falsePosition <= 0 || falsePosition >= 1 {
		db.BloomKeys = bloom.NewWithEstimates(uint(len(keys))+uint(rand.Uint32()%1000+1000), 0.0000001+rand.Float64()/10000000)
	} else {
		db.BloomKeys = bloom.NewWithEstimates(uint(capacity), falsePosition)
	}
	//if type of k is string, then AddString is faster than Add
	for _, it := range keys {
		db.BloomKeys.AddString(it)
	}
	return nil
}
func (db *Ctx[k, v]) TestBloomKey(key k) (exist bool) {
	var (
		keyStr string
		err    error
	)
	if db.BloomKeys == nil {
		log.Fatal("BloomKeys is nil, please BuildKeysBloomFilter first")
	}
	if keyStr, err = db.toKeyStr(key); err != nil {
		log.Fatalf("TestKey -> toKeyStr error: %v", err.Error())
	}
	return db.BloomKeys.TestString(keyStr)
}
func (db *Ctx[k, v]) AddBloomKey(key k) (err error) {
	var (
		keyStr string
	)
	if db.BloomKeys == nil {
		log.Fatal("BloomKeys is nil, please BuildKeysBloomFilter first")
	}
	if keyStr, err = db.toKeyStr(key); err != nil {
		return err
	}
	db.BloomKeys.AddString(keyStr)
	return nil
}
