package collector

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/uhthomas/tesla"
)

func descs(c *VehicleCollector) []*prometheus.Desc {
	return []*prometheus.Desc{
		c.infoDesc,
		c.nameDesc,
		c.stateDesc,
		c.softwareVersionDesc,
		c.odometerMilesSumDesc,
		c.insideTempDesc,
		c.outsideTempDesc,
		c.batteryRatioDesc,
		c.batteryUsableRatioDesc,
		c.batteryIdealMilesDesc,
		c.batteryEstimatedMilesDesc,
		c.chargeVoltsDesc,
		c.chargeAmpsDesc,
		c.chargeAmpsAvailableDesc,
	}
}

func TestNewVehicleCollector(t *testing.T) {
	ctx, s := context.Background(), &tesla.Service{}
	cc := NewVehicleCollector(ctx, s)
	if cc.ctx != ctx {
		t.Fatalf("cc.ctx != ctx")
	}
	if cc.s != s {
		t.Fatal("cc.s != s (tesla service)")
	}
	for i, desc := range descs(cc) {
		if desc == nil {
			t.Fatalf("desc %#v (%d) is nil", desc, i)
		}
	}
}

func TestVehicleCollector_Describe(t *testing.T) {
	c := NewVehicleCollector(context.Background(), nil)

	ch := make(chan *prometheus.Desc)
	go func() {
		defer close(ch)
		c.Describe(ch)
	}()

	m := map[*prometheus.Desc]struct{}{}
	for desc := range ch {
		m[desc] = struct{}{}
	}

	for _, desc := range descs(c) {
		if _, ok := m[desc]; !ok {
			t.Fatalf("missing desc")
		}
	}
}
