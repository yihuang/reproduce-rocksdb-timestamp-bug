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

	version := int64(0)
	itr := db.NewIteratorCF(newTSReadOptions(&version), cfHandle)
	itr.SeekToFirst()
	var pairs []KVPairWithTS
	for ; itr.Valid(); itr.Next() {
		key := moveSliceToBytes(itr.Key())
		value := moveSliceToBytes(itr.Value())

		if binary.LittleEndian.Uint64(itr.Timestamp().Data()) != 0 {
			fmt.Println("skip key", string(key))
			continue
		}

		ts := key[len(key)-tsrocksdb.TimestampSize:]
		key = key[:len(key)-tsrocksdb.TimestampSize]
		pairs = append(pairs, KVPairWithTS{
			Key:       key,
			Value:     value,
			Timestamp: ts,
		})
	}

	defaultSyncWriteOpts := grocksdb.NewDefaultWriteOptions()
	defaultSyncWriteOpts.SetSync(true)
	batch := grocksdb.NewWriteBatch()
	defer batch.Destroy()
	for _, pair := range pairs {
		batch.PutCFWithTS(cfHandle, pair.Key, pair.Timestamp, pair.Value)
		fmt.Printf("fix data: key: %s, ts: %d, value: %s\n", string(pair.Key), binary.LittleEndian.Uint64(pair.Timestamp), string(pair.Value))
	}
	if err := db.Write(defaultSyncWriteOpts, batch); err != nil {
		panic(err)
	}

	/*
		opts := grocksdb.NewDefaultFlushOptions()
		defer opts.Destroy()
		if err := db.FlushCF(cfHandle, opts); err != nil {
			panic(err)
		}
	*/
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
