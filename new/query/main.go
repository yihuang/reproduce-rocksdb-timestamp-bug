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

	version := int64(10000)
	itr := db.NewIteratorCF(newTSReadOptions(&version), cfHandle)
	itr.SeekToFirst()
	count := 0
	for ; itr.Valid(); itr.Next() {
		key := moveSliceToBytes(itr.Key())
		value := moveSliceToBytes(itr.Value())
		ts := binary.LittleEndian.Uint64(itr.Timestamp().Data())

		if ts == 0 {
			// skip 0 timestamp
			fmt.Println("skip", string(key), string(value), ts)
			continue
		}
		count++
		fmt.Println(string(key), string(value))
	}
	fmt.Println("total", count)
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
