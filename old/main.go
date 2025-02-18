package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"

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

	version := uint64(100)
	var ts [tsrocksdb.TimestampSize]byte
	binary.LittleEndian.PutUint64(ts[:], version)

	batch := grocksdb.NewWriteBatch()
	defer batch.Destroy()
	for i := 0; i < 10000; i++ {
		key := []byte("key" + strconv.Itoa(i))
		value := []byte("value" + strconv.Itoa(i))
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
