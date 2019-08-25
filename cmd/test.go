package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/rs485"
	"github.com/volkszaehler/mbmd/meters/sunspec"
)

// testCmd represents the scan command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Scan for attached devices",
	Long: `Scan loops over all device ids from 1 to 254 and tries to
read a common value depending on device type.
For RTU devices the common value is most likely the L1 voltage,
for TCP devices it tries to read the SunSpec common block.
If successful the detected device type and device id are displayed.

Scan will ignore the config file and requires adapter configuration using command line.`,
	Run: test,
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func test(cmd *cobra.Command, args []string) {
	// create connection
	adapter := viper.GetString("adapter")
	if adapter == "" {
		log.Fatal("Missing adapter configuration")
	}

	conn := createConnection(adapter, viper.GetInt("baudrate"), viper.GetString("comset"))
	client := conn.ModbusClient()

	// raw log
	if viper.GetBool("raw") {
		conn.Logger(log.New(os.Stderr, "", log.LstdFlags))
	}

	// create devices
	devices := make([]meters.Device, 0)
	if _, ok := conn.(*meters.TCP); ok {
		suns := sunspec.NewDevice("SUNS")
		devices = append(devices, suns)
	} else {
		for t := range rs485.Producers {
			dev, err := rs485.NewDevice(t)
			if err != nil {
				log.Fatal(err)
			}
			devices = append(devices, dev)
		}
	}

	deviceList := make(map[int]meters.Device)
	log.Printf("starting bus scan on %s", adapter)

	// SCAN:
	// loop over all valid slave adresses
	for deviceID := 1; deviceID <= 2; deviceID++ {
		log.Printf("device %d", deviceID)

		// give the bus some time to recover before querying the next device
		time.Sleep(40 * time.Millisecond)
		conn.Slave(uint8(deviceID))

		for _, dev := range devices {
			if err := dev.Initialize(client); err != nil {
				if _, partial := err.(meters.SunSpecPartiallyInitialized); !partial {
					continue // devices
				}
				log.Println(err) // log error but continue
			}

			mr, err := dev.Probe(client)
			log.Printf("%+v", mr)
			if err == nil {
				log.Printf("device %d: %s type device found, %s: %.2f\r\n",
					deviceID,
					dev.Descriptor().Manufacturer,
					mr.Measurement,
					mr.Value,
				)

				deviceList[deviceID] = dev
				// continue SCAN
			}
		}

		log.Printf("device %d: n/a\r\n", deviceID)
	}

	log.Printf("found %d active devices:\r\n", len(deviceList))
	for deviceID, dev := range deviceList {
		desc := dev.Descriptor()

		s := ""
		addDesc(&s, "Model", desc.Model)
		addDesc(&s, "Version", desc.Version)
		addDesc(&s, "Serial", desc.Serial)

		if s != "" {
			s = fmt.Sprintf("(%s)", s)
		}

		log.Printf(
			"* #%d type %s %s",
			deviceID,
			desc.Manufacturer,
			s,
		)
	}

	log.Println("WARNING: This lists only the devices that responded to " +
		"a known probe request. Devices with different " +
		"function code definitions might not be detected.")
}
