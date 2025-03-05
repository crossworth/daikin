package daikin

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDiscoveryDevices(t *testing.T) {
	t.SkipNow()
	devices, err := DiscoveryDevices(context.Background(), 10*time.Second)
	require.NoError(t, err)
	for _, device := range devices {
		fmt.Println(device.String())
		// Device Hostname=DAIKINXXAXX.local. IP=192.168.0.XX Port=80 APN=DAIKIN:XXAXXXXXXXXC
	}
}
