package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"slices"

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
	dirOut := os.Args[1]

	dbOut, cfHandleOut, err := tsrocksdb.OpenVersionDB(dirOut)
	if err != nil {
		panic(err)
	}

	var ts [tsrocksdb.TimestampSize]byte
	var pairs []KVPairWithTS
	// add intermidiate versions
	for i := 0; i < 10000; i++ {
		key := []byte(fmt.Sprintf("key-%010d", i))
		for j := 0; j < 5; j++ {
			version := j + 10
			binary.LittleEndian.PutUint64(ts[:], uint64(version))
			value := []byte(fmt.Sprintf("value-%d-%d", i, j))
			pairs = append(pairs, KVPairWithTS{
				Key:       key,
				Value:     value,
				Timestamp: slices.Clone(ts[:]),
			})
		}
	}

	// write a pass to make sure the version is updated
	for i := 0; i < 10000; i++ {
		version := i%1000 + 20
		binary.LittleEndian.PutUint64(ts[:], uint64(version))

		key := []byte(fmt.Sprintf("key-%010d", i))
		value := []byte(fmt.Sprintf("value-%d-%d", i, version))

		pairs = append(pairs, KVPairWithTS{
			Key:       key,
			Value:     value,
			Timestamp: slices.Clone(ts[:]),
		})
	}

	defaultSyncWriteOpts := grocksdb.NewDefaultWriteOptions()
	defaultSyncWriteOpts.SetSync(true)

	readOpts := grocksdb.NewDefaultReadOptions()
	defer readOpts.Destroy()

	batch := grocksdb.NewWriteBatch()
	defer batch.Destroy()
	for _, pair := range pairs {
		batch.PutCFWithTS(cfHandleOut, pair.Key, pair.Timestamp, pair.Value)
		fmt.Printf("fix data: key: %s, ts: %d, value: %s\n", string(pair.Key), binary.LittleEndian.Uint64(pair.Timestamp), string(pair.Value))

		// also write the timestamp 0 values
		batch.PutCFWithTS(cfHandleOut, append(pair.Key, pair.Timestamp...), ts[:], pair.Value)
	}
	if err := dbOut.Write(defaultSyncWriteOpts, batch); err != nil {
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
