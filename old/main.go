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

	total := 10000

	// batch := grocksdb.NewWriteBatch()
	// for i := 0; i < total; i++ {
	// 	key := []byte("key" + strconv.Itoa(i))
	// 	value := []byte("value" + strconv.Itoa(i))
	// 	var ts [tsrocksdb.TimestampSize]byte
	// 	binary.LittleEndian.PutUint64(ts[:], uint64(1))
	// 	batch.PutCFWithTS(cfHandle, key, ts[:], value)
	// 	fmt.Println("wrote", string(key), string(value))
	// }
	// err = db.Write(defaultSyncWriteOpts, batch)
	// if err != nil {
	// 	panic(err)
	// }

	// opts := grocksdb.NewDefaultFlushOptions()
	// defer opts.Destroy()

	// err = db.FlushCF(cfHandle, opts)
	// if err != nil {
	// 	panic(err)
	// }
	// batch.Destroy()

	batch := grocksdb.NewWriteBatch()
	defer batch.Destroy()
	count := 2
	for i := 0; i < total; i++ {
		key := []byte("key" + strconv.Itoa(i))
		for j := 1; j <= count; j++ {
			value := []byte("value" + strconv.Itoa(i+j%2))
			var ts [tsrocksdb.TimestampSize]byte
			binary.LittleEndian.PutUint64(ts[:], uint64(j+1))
			batch.PutCFWithTS(cfHandle, key, ts[:], value)
			fmt.Println("wrote", string(key), string(value))
		}
	}

	err = db.Write(defaultSyncWriteOpts, batch)
	if err != nil {
		panic(err)
	}

	/*
		err = db.FlushCF(cfHandle, opts)
		if err != nil {
			panic(err)
		}
	*/
}
