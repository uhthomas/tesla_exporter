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
		c: c,
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

	v, err := c.c.Vehicles(ctx)
	if err != nil {
		panic(err)
	}

	for _, vv := range v {
		ch <- prometheus.MustNewConstMetric(c.infoDesc, prometheus.GaugeValue, 1,
			strconv.FormatUint(vv.ID, 10),
			strconv.FormatUint(vv.VehicleID, 10),
			vv.VIN,
			vv.DisplayName,
			vv.Color,
			vv.State,
		)
	}
}
