package collector

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/uhthomas/tesla"
)

type VehicleCollector struct {
	ctx     context.Context
	s       *tesla.Service
	expire  time.Duration
	m       *sync.RWMutex
	cond    *sync.Cond
	metrics []prometheus.Metric

	infoDesc,
	nameDesc,
	stateDesc,
	softwareVersionDesc,
	odometerMilesSumDesc,
	insideTempDesc,
	outsideTempDesc,
	batteryRatioDesc,
	batteryUsableRatioDesc,
	batteryIdealMilesDesc,
	batteryEstimatedMilesDesc,
	chargeVoltsDesc,
	chargeAmpsDesc,
	chargeAmpsAvailableDesc *prometheus.Desc
}

func NewVehicleCollector(ctx context.Context, s *tesla.Service, expire time.Duration) *VehicleCollector {
	m := &sync.RWMutex{}
	c := &VehicleCollector{
		ctx:                       ctx,
		s:                         s,
		expire:                    expire,
		m:                         m,
		cond:                      sync.NewCond(m),
		infoDesc:                  prometheus.NewDesc("tesla_vehicle_info", "Tesla vehicle info.", []string{"vin", "id", "vehicle_id"}, nil),
		nameDesc:                  prometheus.NewDesc("tesla_vehicle_name", "Tesla vehicle name.", []string{"vin", "name"}, nil),
		stateDesc:                 prometheus.NewDesc("tesla_vehicle_state", "Tesla vehicle state.", []string{"vin", "state"}, nil),
		softwareVersionDesc:       prometheus.NewDesc("tesla_vehicle_software_version", "Tesla vehicle software version.", []string{"vin", "software_version"}, nil),
		odometerMilesSumDesc:      prometheus.NewDesc("tesla_vehicle_odometer_miles_total", "Tesla vehicle odometer miles.", []string{"vin"}, nil),
		insideTempDesc:            prometheus.NewDesc("tesla_vehicle_inside_temp_celsius", "Tesla vehicle inside temperature.", []string{"vin"}, nil),
		outsideTempDesc:           prometheus.NewDesc("tesla_vehicle_outside_temp_celsius", "Tesla vehicle outside temperature.", []string{"vin"}, nil),
		batteryRatioDesc:          prometheus.NewDesc("tesla_vehicle_battery_ratio", "Tesla vehicle battery ratio.", []string{"vin"}, nil),
		batteryUsableRatioDesc:    prometheus.NewDesc("tesla_vehicle_battery_usable_ratio", "Tesla vehicle battery usable ratio.", []string{"vin"}, nil),
		batteryIdealMilesDesc:     prometheus.NewDesc("tesla_vehicle_battery_ideal_miles", "Tesla vehicle battery ideal miles.", []string{"vin"}, nil),
		batteryEstimatedMilesDesc: prometheus.NewDesc("tesla_vehicle_battery_estimated_miles", "Tesla vehicle battery estimated miles", []string{"vin"}, nil),
		chargeVoltsDesc:           prometheus.NewDesc("tesla_vehicle_charge_volts", "Tesla vehicle charge volts.", []string{"vin"}, nil),
		chargeAmpsDesc:            prometheus.NewDesc("tesla_vehicle_charge_amps", "Tesla vehicle charge amps.", []string{"vin"}, nil),
		chargeAmpsAvailableDesc:   prometheus.NewDesc("tesla_vehicle_charge_amps_available", "Tesla vehicle charge amps available.", []string{"vin"}, nil),
	}
	return c
}

func (c *VehicleCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.infoDesc
	ch <- c.nameDesc
	ch <- c.stateDesc
	ch <- c.softwareVersionDesc
	ch <- c.odometerMilesSumDesc
	ch <- c.insideTempDesc
	ch <- c.outsideTempDesc
	ch <- c.batteryRatioDesc
	ch <- c.batteryUsableRatioDesc
	ch <- c.batteryIdealMilesDesc
	ch <- c.batteryEstimatedMilesDesc
	ch <- c.chargeVoltsDesc
	ch <- c.chargeAmpsDesc
	ch <- c.chargeAmpsAvailableDesc
}

func (c *VehicleCollector) Collect(ch chan<- prometheus.Metric) {
	c.cond.Signal()

	c.m.RLock()
	defer c.m.RUnlock()

	for _, m := range c.metrics {
		ch <- m
	}
}

func (c *VehicleCollector) collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(c.ctx, time.Minute)
	defer cancel()

	vs, err := c.s.Vehicles(ctx)
	if err != nil {
		log.Printf("list vehicles: %v", err)
		return
	}

	for _, v := range vs {
		m := metricMaker{ch: ch, vin: v.VIN}
		m.gauge(c.infoDesc, 1,
			strconv.FormatUint(v.ID, 10),
			strconv.FormatUint(v.VehicleID, 10),
		)
		m.gauge(c.nameDesc, 1, v.DisplayName)
		m.gauge(c.stateDesc, 1, v.State)

		// detailed information is not available for sleeping or in
		// service vehicles.
		if v.State != "online" || v.InService {
			continue
		}

		vv, err := c.s.Vehicle(ctx, v.ID)
		if err != nil {
			log.Printf("get vehicle %d: %v", v.ID, err)
			continue
		}

		m.gauge(c.softwareVersionDesc, 1, vv.VehicleState.CarVersion)
		// really this shouldn't be a gauge, as the value can never
		// decrease.
		m.gauge(c.odometerMilesSumDesc, vv.VehicleState.Odometer)
		m.gauge(c.insideTempDesc, vv.ClimateState.InsideTemp)
		m.gauge(c.outsideTempDesc, vv.ClimateState.OutsideTemp)
		m.gauge(c.batteryRatioDesc, vv.ChargeState.BatteryLevel/100)
		m.gauge(c.batteryUsableRatioDesc, vv.ChargeState.UsableBatteryLevel/100)
		m.gauge(c.batteryIdealMilesDesc, vv.ChargeState.BatteryRange)
		m.gauge(c.batteryEstimatedMilesDesc, vv.ChargeState.EstBatteryRange)
		m.gauge(c.chargeVoltsDesc, float64(vv.ChargeState.ChargerVoltage))
		m.gauge(c.chargeAmpsDesc, float64(vv.ChargeState.ChargerActualCurrent))
		m.gauge(c.chargeAmpsAvailableDesc, float64(vv.ChargeState.ChargerPilotCurrent))
	}
}

func (c *VehicleCollector) Refresh() {
	var last time.Time
	for {
		c.cond.L.Lock()

		for time.Since(last) < c.expire {
			c.cond.Wait()
		}

		cc := make(chan prometheus.Metric)

		go func() {
			defer close(cc)
			c.collect(cc)
		}()

		c.metrics = c.metrics[:0]
		for m := range cc {
			c.metrics = append(c.metrics, m)
		}
		last = time.Now()
		c.cond.L.Unlock()
	}
}

type metricMaker struct {
	ch  chan<- prometheus.Metric
	vin string
}

func (m *metricMaker) gauge(desc *prometheus.Desc, value float64, labelValues ...string) {
	m.ch <- prometheus.MustNewConstMetric(
		desc,
		prometheus.GaugeValue,
		value,
		append([]string{m.vin}, labelValues...)...,
	)
}
