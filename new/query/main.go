package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"github.com/crypto-org-chain/cronos/versiondb/tsrocksdb"
	"github.com/linxGnu/grocksdb"
)

func newTSReadOptions(version *int64) *grocksdb.ReadOptions {
	var ver uint64
	if version == nil {
		ver = math.MaxUint64
	} else {
		ver = uint64(*version)
	}

	var ts [tsrocksdb.TimestampSize]byte
	binary.LittleEndian.PutUint64(ts[:], ver)

	readOpts := grocksdb.NewDefaultReadOptions()
	readOpts.SetTimestamp(ts[:])
	return readOpts
}

type KVPairWithTS struct {
	Key       []byte
	Value     []byte
	Timestamp []byte
}

func main() {
	dir := os.Args[1]

	db, cfHandle, err := tsrocksdb.OpenVersionDB(dir)
	if err != nil {
		panic(err)
	}

	version := int64(100000)
	itr := db.NewIteratorCF(newTSReadOptions(&version), cfHandle)
	itr.SeekToFirst()
	counter := 0
	for ; itr.Valid(); itr.Next() {
		key := moveSliceToBytes(itr.Key())
		value := moveSliceToBytes(itr.Value())

		if binary.LittleEndian.Uint64(itr.Timestamp().Data()) == 0 {
			// skip 0 timestamp
			continue
		}

		if string(key) != fmt.Sprintf("key-%010d", counter) {
			panic(fmt.Sprintf("wrong key: %s", string(key)))
		}
		if string(value) != fmt.Sprintf("value-%d-%d", counter, counter%1000+20) {
			panic(fmt.Sprintf("wrong value: %s, %s", string(key), string(value)))
		}
		fmt.Println(string(key), string(value))
		counter++
	}
}

func moveSliceToBytes(s *grocksdb.Slice) []byte {
	defer s.Free()
	if !s.Exists() {
		return nil
	}
	v := make([]byte, len(s.Data()))
	copy(v, s.Data())
	return v
}
