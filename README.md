Emit marker into stdout on receipt of signal. Designed to run in a shell, wrapping another program's output via a pipe.

Loosely modeled after https://github.com/firesock/pipe-marker.

# Usage

## Build and install the binary

Run this from the repository root:

```shell
go install
```

Output should be under your `$GOPATH/bin`.


## Use in a pipe

Replacing with real paths, enter into a terminal window:

```shell
/path/to/noisy-program | /path/to/pipe-marker-go -enabled
```

## Send signals

While running, from another terminal window, send signals to place markers in the stream. 

Supported signals: SIGUSR1, SIGUSR2.

```shell
pid=$(ps -eo pid,comm | grep pipe-marker-go | awk '{print $1}' | head -n1)
kill -s USR1 $pid
```