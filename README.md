```
$ ./run.sh
```

* it use `nix` to build depencencies.

* after run `run.sh`, you might need to manually run query several times to reproduce the error:

  ```
  $ nix develop -c go run ./query/main.go /tmp/versiondb | wc -l
  ```

  until you get something like:

  ```
  ...
  key-0000004511 value-4511-531
  key-0000004512 value-4512-532
  panic: wrong key: key-0000009033, value-9033-53, 53
  
  goroutine 1 [running]:
  main.main()
  	$HOME/src/reproduce-rocksdb-timestamp-bug/new/query/main.go:70 +0x560
  exit status 2
  ```
