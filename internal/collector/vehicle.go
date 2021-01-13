package collector

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/uhthomas/tesla_exporter/internal/tesla"
)

type VehicleCollector struct {
	ctx context.Context
	c   *tesla.Client
	infoDesc,
	insideTempDesc,
	outsideTempDesc *prometheus.Desc
}

func NewVehicleCollector(ctx context.Context, c *tesla.Client) *VehicleCollector {
	return &VehicleCollector{
		ctx: ctx,
		c:   c,
		infoDesc: prometheus.NewDesc("tesla_vehicle_info", "Tesla vehicle info.", []string{
			"id", "vehicle_id", "vin", "name", "color", "state",
		}, nil),
		insideTempDesc:  prometheus.NewDesc("tesla_vehicle_inside_temp_celsius", "Tesla vehicle inside temperature.", []string{"vin"}, nil),
		outsideTempDesc: prometheus.NewDesc("tesla_vehicle_outside_temp_celsius", "Tesla vehicle outside temperature.", []string{"vin"}, nil),
	}
}

func (c *VehicleCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.infoDesc
	ch <- c.insideTempDesc
	ch <- c.outsideTempDesc
}

func (c *VehicleCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()

	vs, err := c.c.Vehicles(ctx)
	if err != nil {
		panic(err)
	}

	for _, v := range vs {
		ch <- prometheus.MustNewConstMetric(c.infoDesc, prometheus.GaugeValue, 1,
			strconv.FormatUint(v.ID, 10),
			strconv.FormatUint(v.VehicleID, 10),
			v.VIN,
			v.DisplayName,
			v.VehicleConfig.ExteriorColor,
			v.State,
		)
		ch <- prometheus.MustNewConstMetric(c.insideTempDesc, prometheus.GaugeValue, v.ClimateState.InsideTemp, v.VIN)
		ch <- prometheus.MustNewConstMetric(c.outsideTempDesc, prometheus.GaugeValue, v.ClimateState.OutsideTemp, v.VIN)
	}
}
