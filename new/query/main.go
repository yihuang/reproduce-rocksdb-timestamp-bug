package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"example.com/m"
	"github.com/linxGnu/grocksdb"
)

func newTSReadOptions(version *int64) *grocksdb.ReadOptions {
	var ver uint64
	if version == nil {
		ver = math.MaxUint64
	} else {
		ver = uint64(*version)
	}

	var ts [m.TimestampSize]byte
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

	db, cfHandle, err := m.OpenDB(dir)
	if err != nil {
		panic(err)
	}

	version := int64(100000)
	readOpts := newTSReadOptions(&version)
	readOpts2 := newTSReadOptions(&version)

	defer func() {
		readOpts.Destroy()
		readOpts2.Destroy()
	}()

	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%010d", i)
		data, err := db.GetCF(readOpts, cfHandle, []byte(key))
		if err != nil {
			panic(err)
		}
		if string(data.Data()) != fmt.Sprintf("value-%d-%d", i, i%1000+20) {
			panic(fmt.Sprintf("wrong value: %s, %s", key, string(data.Data())))
		}
		data.Free()
	}

	itr := db.NewIteratorCF(readOpts2, cfHandle)
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
			panic(fmt.Sprintf("wrong key: %s, %s, %d", string(key), string(value), binary.LittleEndian.Uint64(itr.Timestamp().Data())))
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
