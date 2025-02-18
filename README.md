1. Generate data using old version, keep it in WAL.

    ```
    cd old
    nix develop -c go run ./main.go /tmp/versiondb
    ```

2. Try to query data with new version, don't see it.

   ```
   cd new
   nix develop -c go run ./query/main.go /tmp/versiondb
   ```

3. Fix data

    ```
    cd new
    nix develop -c go run ./fix/main.go /tmp/versiondb
    ```

4. Query again, data is fixed.

    ```
    cd new
    nix develop -c go run ./query/main.go /tmp/versiondb
    ```

    