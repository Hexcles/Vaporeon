## Instructions

1.  Run `make` to build all binaries.
2.  In a terminal, `cd server`, and run:
    ```bash
    sudo systemd-run -p "Delegate=yes" -p "User=$USER" -P --wait ./server
    ```
    (You could also use `nobody` instead of `$USER` to further lock down permissions.)

    If you don't need cgroup, you can also just run `./server`. See [this doc][1] for more
    on systemd.
3.  In another terminal, `cd client`.
    1.  `./client` by default runs as admin.
    2.  Check out `./client -h`.
    3.  `./client -cert_file ../certs/client2_cert.pem -key_file ../certs/client2_key.pem`
        runs as a guest.
    4.  Guests can't stream/query/kill jobs launched by the admin. They can't shut down
	the server, either.
    5.  Admins can do anything :).

[1]: worker/example#no-systemd
