package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	geohash "github.com/mmcloughlin/geohash"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type dump900AircraftStruct struct {
	Now      float64 `json:"now"`
	Messages int     `json:"messages"`
	Aircraft []struct {
		Hex            string        `json:"hex,omitempty"`
		Flight         string        `json:"flight,omitempty"`
		AltBaro        int           `json:"alt_baro,omitempty"`
		AltGeom        int           `json:"alt_geom,omitempty"`
		Altitude       int           `json:"altitude,omitempty"`
		Gs             float64       `json:"gs,omitempty"`
		Ias            int           `json:"ias,omitempty"`
		Tas            int           `json:"tas,omitempty"`
		Mach           float64       `json:"mach,omitempty"`
		Track          float64       `json:"track,omitempty"`
		TrackRate      float64       `json:"track_rate,omitempty"`
		Roll           float64       `json:"roll,omitempty"`
		MagHeading     float64       `json:"mag_heading,omitempty"`
		BaroRate       int           `json:"baro_rate,omitempty"`
		GeomRate       int           `json:"geom_rate,omitempty"`
		Squawk         string        `json:"squawk,omitempty"`
		Category       string        `json:"category,omitempty"`
		NavQnh         float64       `json:"nav_qnh,omitempty"`
		NavAltitudeMcp int           `json:"nav_altitude_mcp,omitempty"`
		Lat            float64       `json:"lat,omitempty"`
		Lon            float64       `json:"lon,omitempty"`
		Nic            int           `json:"nic,omitempty"`
		Rc             int           `json:"rc,omitempty"`
		SeenPos        float64       `json:"seen_pos,omitempty"`
		Version        int           `json:"version,omitempty"`
		NacP           int           `json:"nac_p,omitempty"`
		NacV           int           `json:"nac_v,omitempty"`
		Sil            int           `json:"sil,omitempty"`
		SilType        string        `json:"sil_type,omitempty"`
		Mlat           []interface{} `json:"mlat,omitempty"`
		Tisb           []interface{} `json:"tisb,omitempty"`
		Messages       int           `json:"messages,omitempty"`
		Seen           float64       `json:"seen,omitempty"`
		Rssi           float64       `json:"rssi,omitempty"`
		Emergency      string        `json:"emergency,omitempty"`
		NavHeading     float64       `json:"nav_heading,omitempty"`
		NicBaro        int           `json:"nic_baro,omitempty"`
		Gva            int           `json:"gva,omitempty"`
		Sda            int           `json:"sda,omitempty"`
	} `json:"aircraft"`
}

var addr = flag.String("listen-address", ":9467",
	"The address to listen on for HTTP requests.")

func main() {
	flag.Parse()

	trackingACcount := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "dump1090_ac_tracking_now_count",
			Help: "Count of tracking ACs in realtime",
		})
	prometheus.MustRegister(trackingACcount)

	catchedMessages := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "dump1090_messages_catched_count",
			Help: "Total messages catched with dump1090",
		})
	prometheus.MustRegister(catchedMessages)

	var trackingACs *prometheus.GaugeVec
	trackingACs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dump1090_ac_tracking_flights",
			Help: "Tracking ACs",
		},
		[]string{"flight", "geohash", "lat", "long", "altitude", "squawk"},
	)
	go func() {
		for {
			aircraftFile, err := os.Open("/run/dump1090-fa/aircraft.json")
			if err != nil {
				fmt.Println(err)
			}

			defer aircraftFile.Close()

			byteValue, _ := ioutil.ReadAll(aircraftFile)

			var aircrafts dump900AircraftStruct

			err = json.Unmarshal(byteValue, &aircrafts)
			if err != nil {
				log.Println("Error unmarshaling aircraft.json")
			}

			trackingACs.Reset()
			trackingACcount.Set(0)

			for _, aircraft := range aircrafts.Aircraft {
				if len(aircraft.Flight) > 0 {
					trackingACcount.Inc()

					reg, err := regexp.Compile("[^a-zA-Z0-9]+")
					if err != nil {
						log.Fatal(err)
					}
					aircraftFlight := reg.ReplaceAllString(aircraft.Flight, "")
					geoHash := geohash.Encode(aircraft.Lat, aircraft.Lon)
					alt := aircraft.AltGeom
					if aircraft.Altitude > 0 {
						alt = aircraft.Altitude
					}

					trackingACs.WithLabelValues(aircraftFlight, geoHash, fmt.Sprintf("%2.6f", aircraft.Lat), fmt.Sprintf("%2.6f", aircraft.Lon), fmt.Sprintf("%d", alt), aircraft.Squawk).Set(1)
					prometheus.Register(trackingACs)
				}
			}

			catchedMessages.Set(float64(aircrafts.Messages))
			time.Sleep(1000 * time.Millisecond)

			aircraftFile.Close()
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Starting web server at %s\n", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Printf("http.ListenAndServer: %v\n", err)
	}
}
