package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/takjn/GRDisplayApp/pkg/parser"
	"go.bug.st/serial"
	"go.bug.st/serial.v1/enumerator"
)

// the device VID and PID
const (
	vid = "1f00"
	pid = "2012"
)

func findDevice() (string, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return "", err
	}
	if len(ports) == 0 {
		return "", fmt.Errorf("no serial ports found")
	}
	for _, port := range ports {
		if port.IsUSB && port.VID == vid && port.PID == pid {
			return port.Name, nil
		}
	}

	return "", fmt.Errorf("no device found")
}

func main() {
	var path, filename, device string
	var interval int
	flag.StringVar(&device, "d", "", "device port")
	flag.StringVar(&path, "p", "", "output path")
	flag.StringVar(&filename, "f", "test", "filename")
	flag.IntVar(&interval, "i", 1000, "jpeg interval (ms)")
	flag.Parse()

	binaryMode := true
	if path != "" {
		binaryMode = false
		filename = filename + "%d.jpg"
		path = filepath.Join(path, filename)
	}

	// Find a Mbed CDC
	if device == "" {
		name, err := findDevice()
		if err != nil {
			log.Fatalf("failed to find device: %v", err)
			os.Exit(1)
		}
		device = name
	}

	// Open the CDC
	port, err := serial.Open(device, &serial.Mode{})
	if err != nil {
		log.Fatalf("failed to open device: %v", err)
		os.Exit(1)
	}
	defer port.Close()

	go func() {
		i := 0
		for {
			data, err := parser.GetData(port)
			if err != nil {
				return
			}
			if binaryMode {
				os.Stdout.Write(data)
			} else {
				f := fmt.Sprintf(path, i)
				err = ioutil.WriteFile(f, data, 0666)
				if err != nil {
					log.Fatalf("failed to write file: %v", err)
					break
				}
				log.Println(f, ",", len(data))
				time.Sleep(time.Millisecond * time.Duration(interval))
				i++
			}
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
