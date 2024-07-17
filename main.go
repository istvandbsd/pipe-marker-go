package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var signalMarkers = map[os.Signal]string{
	syscall.SIGUSR1: "===USR1===",
	syscall.SIGUSR2: "===USR2===",
}

// Reads from stdin and sends output forward.
func reader(send chan<- []byte) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		send <- []byte(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, "encountered while reading stdin:", err)
		if err != nil {
			return
		}
	}
	close(send)

}

// Catches signals on the dedicated channel and inserts marker into stream on another channel
func marker(notif <-chan os.Signal, out chan<- []byte, enabled bool) {
	if enabled {
		for {
			s := <-notif
			sigInstance := s.(syscall.Signal)
			emit := []byte(signalMarkers[sigInstance])
			out <- emit
		}
	}
}

// Writes the potentially aggregated stream to stdout.
func writer(recv <-chan []byte, done chan os.Signal) {
	writer := bufio.NewWriter(os.Stdout)
	for {
		stream, ok := <-recv
		if !ok {
			done <- syscall.SIGINT
			return
		}

		stream = append(stream, '\n')

		if _, err := writer.Write(stream); err != nil {
			fmt.Fprintln(os.Stderr, "write failed:", err)
		}

		if err := writer.Flush(); err != nil {
			fmt.Fprintln(os.Stderr, "could not flush:", err)
		}

	}
}

func main() {

	enabled := flag.Bool("enabled", false, "process signals on/off")
	flag.Parse()

	signals := []os.Signal{syscall.SIGUSR1, syscall.SIGUSR2}
	markerChan := make(chan os.Signal, 2)
	signal.Notify(markerChan, signals...)

	doneChan := make(chan os.Signal)
	dataChan := make(chan []byte)

	go func() {
		reader(chan<- []byte(dataChan))
	}()
	go writer(dataChan, doneChan)
	go marker(markerChan, dataChan, *enabled)

	// wait for the goroutines to finish
	<-doneChan
}
