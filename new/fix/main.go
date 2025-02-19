package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"slices"

	"example.com/m"
	"github.com/linxGnu/grocksdb"
)

type KVPairWithTS struct {
	Key       []byte
	Value     []byte
	Timestamp []byte
}

func main() {
	dir := os.Args[1]

	db, cfHandle, err := m.OpenDB(dir)
	if err != nil {
		panic(err)
	}

	var ts [m.TimestampSize]byte
	var pairs []KVPairWithTS
	// add intermidiate versions
	for i := 0; i < 10000; i++ {
		key := []byte(fmt.Sprintf("key-%010d", i))
		for j := 0; j < 5; j++ {
			version := j + 10
			binary.LittleEndian.PutUint64(ts[:], uint64(version))
			value := []byte(fmt.Sprintf("value-%d-%d", i, version))
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

	// clear ts to 0
	binary.LittleEndian.PutUint64(ts[:], 0)

	batch := grocksdb.NewWriteBatch()
	defer batch.Destroy()
	for _, pair := range pairs {
		batch.PutCFWithTS(cfHandle, pair.Key, pair.Timestamp, pair.Value)
		fmt.Printf("fix data: key: %s, ts: %d, value: %s\n", string(pair.Key), binary.LittleEndian.Uint64(pair.Timestamp), string(pair.Value))

		// also write the timestamp 0 values
		batch.PutCFWithTS(cfHandle, append(pair.Key, pair.Timestamp...), ts[:], pair.Value)
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
