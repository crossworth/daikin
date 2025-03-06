package daikin

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
)

// Device represent a Daikin device found on the network.
type Device struct {
	Port     int    `json:"port"`
	APN      string `json:"apn"`
	Hostname string `json:"hostname"`
	IP       net.IP `json:"ip"`
}

func (d Device) String() string {
	return fmt.Sprintf("Device Hostname=%s IP=%s Port=%d APN=%s", d.Hostname, d.IP.String(), d.Port, d.APN)
}

// DiscoveryDevices discovery all the Daikin devices on the network.
// It will respect the context timeout/cancellation and accept a timeout as parameter as well.
func DiscoveryDevices(ctx context.Context, timeout time.Duration) ([]Device, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, fmt.Errorf("could not initialize mDNS resolver: %v", err)
	}
	var (
		entriesChan = make(chan *zeroconf.ServiceEntry)
		errChannel  = make(chan error, 1)
	)
	go func() {
		errChannel <- resolver.Browse(ctx, "_iota._tcp", "local.", entriesChan)
	}()

	var devices []Device
	for {
		select {
		case err := <-errChannel:
			if err != nil {
				return nil, fmt.Errorf("discovery error: %w", err)
			}
		case entry, ok := <-entriesChan:
			if !ok {
				return devices, nil
			}
			if dev, valid := processEntry(entry); valid {
				devices = append(devices, dev)
			}
		case <-ctx.Done():
			for entry := range entriesChan {
				if dev, valid := processEntry(entry); valid {
					devices = append(devices, dev)
				}
			}
			return devices, nil
		}
	}
}

// processEntry process a single [*zeroconf.ServiceEntry].
func processEntry(entry *zeroconf.ServiceEntry) (Device, bool) {
	if entry.Instance != `Daikin\ Smart\ AC` || len(entry.AddrIPv4) == 0 {
		return Device{}, false
	}
	var apn string
	for _, t := range entry.Text {
		if strings.HasPrefix(t, "apn=") {
			apn = strings.TrimPrefix(t, "apn=")
			break
		}
	}
	return Device{
		Port:     entry.Port,
		APN:      apn,
		Hostname: entry.HostName,
		IP:       entry.AddrIPv4[0],
	}, true
}

// ConvertAPNToThingAPN converts the APN received on the DiscoveryDevices for the ThingAPN name.
func ConvertAPNToThingAPN(apn string) string {
	apn = strings.ReplaceAll(apn, "DAIKIN:", "")
	if len(apn) < 12 {
		return apn
	}
	var sb strings.Builder
	sb.WriteString("DAIKIN")
	sb.WriteString(apn[4:6])
	sb.WriteString(apn[2:4])
	sb.WriteString(apn[0:2])
	return sb.String()
}
