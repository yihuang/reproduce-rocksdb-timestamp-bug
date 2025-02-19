package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/crypto-org-chain/cronos/versiondb/tsrocksdb"
	"github.com/linxGnu/grocksdb"
)

func main() {
	dir := os.Args[1]
	db, cfHandle, err := tsrocksdb.OpenVersionDB(dir)
	if err != nil {
		panic(err)
	}

	defaultSyncWriteOpts := grocksdb.NewDefaultWriteOptions()
	defaultSyncWriteOpts.SetSync(true)

	var ts [tsrocksdb.TimestampSize]byte

	batch := grocksdb.NewWriteBatch()
	defer batch.Destroy()

	for i := 0; i < 10000; i++ {
		// add intermidiate versions
		key := []byte(fmt.Sprintf("key-%10d", i))
		for j := 0; j < 5; j++ {
			version := j + 10
			binary.LittleEndian.PutUint64(ts[:], uint64(version))

			value := []byte(fmt.Sprintf("value-%d-%d", i, j))
			batch.PutCFWithTS(cfHandle, key, ts[:], value)
			fmt.Println("wrote", string(key), string(value), version)
		}
	}

	// write a pass to make sure the version is updated
	for i := 0; i < 10000; i++ {
		version := i%1000 + 20
		binary.LittleEndian.PutUint64(ts[:], uint64(version))

		key := []byte(fmt.Sprintf("key-%10d", i))
		value := []byte(fmt.Sprintf("value-%d-%d", i, version))
		batch.PutCFWithTS(cfHandle, key, ts[:], value)
		fmt.Println("wrote", string(key), string(value), version)
	}
	err = db.Write(defaultSyncWriteOpts, batch)
	if err != nil {
		panic(err)
	}

	opts := grocksdb.NewDefaultFlushOptions()
	defer opts.Destroy()

	/*
		err = db.FlushCF(cfHandle, opts)
		if err != nil {
			panic(err)
		}
	*/
}
