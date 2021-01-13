package collector

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/uhthomas/tesla_exporter/internal/tesla"
)

func TestNewVehicleCollector(t *testing.T) {
	ctx, c := context.Background(), &tesla.Client{}
	cc := NewVehicleCollector(ctx, c)
	if cc.ctx != ctx {
		t.Fatalf("cc.ctx != ctx")
	}
	if cc.c != c {
		t.Fatal("cc.c != c (tesla client)")
	}
	for i, desc := range []*prometheus.Desc{
		cc.infoDesc,
		cc.nameDesc,
		cc.stateDesc,
		cc.softwareVersionDesc,
		cc.odometerMilesSumDesc,
		cc.insideTempDesc,
		cc.outsideTempDesc,
		cc.batteryRatioDesc,
		cc.batteryUsableRatioDesc,
		cc.batteryIdealMilesDesc,
		cc.batteryEstimatedMilesDesc,
		cc.chargeVoltsDesc,
		cc.chargeAmpsDesc,
		cc.chargeAmpsAvailableDesc,
	} {
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

	for want := range map[*prometheus.Desc]struct{}{
		c.infoDesc:                  {},
		c.nameDesc:                  {},
		c.stateDesc:                 {},
		c.softwareVersionDesc:       {},
		c.odometerMilesSumDesc:      {},
		c.insideTempDesc:            {},
		c.outsideTempDesc:           {},
		c.batteryRatioDesc:          {},
		c.batteryUsableRatioDesc:    {},
		c.batteryIdealMilesDesc:     {},
		c.batteryEstimatedMilesDesc: {},
		c.chargeVoltsDesc:           {},
		c.chargeAmpsDesc:            {},
		c.chargeAmpsAvailableDesc:   {},
	} {
		if _, ok := m[want]; !ok {
			t.Fatalf("missing desc")
		}
	}
}
