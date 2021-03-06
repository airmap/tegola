package server_test

import (
	"context"

	"github.com/airmap/tegola"
	"github.com/airmap/tegola/atlas"
	"github.com/airmap/tegola/basic"
	"github.com/airmap/tegola/mvt"
	"github.com/airmap/tegola/server"
)

//	test server config
const (
	httpPort       = ":8080"
	serverVersion  = "0.4.0"
	serverHostName = "tegola.io"
)

var (
	testMapName        = "test-map"
	testMapAttribution = "test attribution"
	testMapCenter      = [3]float64{1.0, 2.0, 3.0}
)

type testMVTProvider struct{}

func (tp *testMVTProvider) MVTLayer(ctx context.Context, layerName string, tile tegola.Tile, tags map[string]interface{}) (*mvt.Layer, error) {
	var layer mvt.Layer

	return &layer, nil
}

func (tp *testMVTProvider) Layers() ([]mvt.LayerInfo, error) {
	return []mvt.LayerInfo{
		layer{
			name:     "test-layer",
			geomType: basic.Polygon{},
			srid:     tegola.WebMercator,
		},
	}, nil
}

var testLayer1 = atlas.Layer{
	Name:              "test-layer",
	ProviderLayerName: "test-layer-1",
	MinZoom:           4,
	MaxZoom:           9,
	Provider:          &testMVTProvider{},
	GeomType:          basic.Point{},
	DefaultTags: map[string]interface{}{
		"foo": "bar",
	},
}

var testLayer2 = atlas.Layer{
	Name:              "test-layer-2-name",
	ProviderLayerName: "test-layer-2-provider-layer-name",
	MinZoom:           10,
	MaxZoom:           20,
	Provider:          &testMVTProvider{},
	GeomType:          basic.Line{},
	DefaultTags: map[string]interface{}{
		"foo": "bar",
	},
}

var testLayer3 = atlas.Layer{
	Name:              "test-layer",
	ProviderLayerName: "test-layer-3",
	MinZoom:           10,
	MaxZoom:           20,
	Provider:          &testMVTProvider{},
	GeomType:          basic.Point{},
	DefaultTags:       map[string]interface{}{},
}

type layer struct {
	name     string
	geomType tegola.Geometry
	srid     int
}

func (l layer) Name() string {
	return l.name
}

func (l layer) GeomType() tegola.Geometry {
	return l.geomType
}

func (l layer) SRID() int {
	return l.srid
}

func init() {
	server.Version = serverVersion
	server.HostName = serverHostName

	testMap := atlas.NewWGS84Map(testMapName)
	testMap.Attribution = testMapAttribution
	testMap.Center = testMapCenter
	testMap.Layers = append(testMap.Layers, []atlas.Layer{
		testLayer1,
		testLayer2,
		testLayer3,
	}...)

	//	register a map with atlas
	atlas.AddMap(testMap)
}
