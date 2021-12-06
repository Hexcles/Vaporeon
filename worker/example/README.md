# Minimum example for worker library

## Instructions

Assuming you have a recent version of Linux kernel and systemd (tested on Linux 5.10 LTS
and systemd 294), run the following commands in this directory:

```bash
go build .
sudo systemd-run -p "Delegate=yes" -p "User=$USER" -P --wait ./example
```

Explanation of the `systemd-run` command:

*   `systemd-run` creates a transient systemd service so we don't need to install a
    service unit.
*   `sudo`: We need root to invoke the system instance of systemd which can provide us
    with the cgroups we need. Note that the example program is **not** run with root.
*   `-p "Delegate=yes"`: Ask systemd to delegate all available cgroup controllers to us.
    Note this does **not** actually set up cgroups in any way.
*   `-p "User=$USER"`: Run the service as the current user instead of root.
*   `-P`: Print the stdout/stderr of the command.
*   `--wait`: Wait until the command exits.

Read the outputs carefully and you can see it runs a series of commands in a new namespace
which produces some interesting results.

Now if you don't see any cgroup warnings and have more than 128 MiB available memory, try
the following command:

```bash
sudo systemd-run -p "Delegate=yes" -p "User=$USER" -P --wait ./example tail /dev/zero
```

This time the command should fail quickly.

If you have GNU time, you can try the following to find out what happened:

```bash
sudo systemd-run -p "Delegate=yes" -p "User=$USER" -P --wait ./example /usr/bin/time -v tail /dev/zero
```

`tail /dev/zero` tries to get to the end of an infinite stream and runs out of memory.
Thanks to the worker library, we have set the memory limit to 128 MiB and OOM killer
kicked in.

### No systemd?

We don't really need systemd. The actual requirements are:

* Pure cgroupv2 (unified hierarchy) is used and mounted with `nsdelegate`.
* The program needs to start in a cgroup hierarchy it can control, namely:
    * It needs to have write access to the hierarchy.
    * There is no other process in the hierarchy.
    * memory, io, cpu controllers are available in the hierarchy.
