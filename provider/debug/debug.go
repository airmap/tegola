//	The debug provider returns features that are helpful for debugging a tile
//	including a box for the tile edges and a point in the middle of the tile
//	with z,x,y values encoded
package debug

import (
	"context"
	"fmt"

	"github.com/airmap/tegola"
	"github.com/airmap/tegola/basic"
	"github.com/airmap/tegola/mvt"
	"github.com/airmap/tegola/mvt/provider"
)

const Name = "debug"

const (
	LayerDebugTileOutline = "debug-tile-outline"
	LayerDebugTileCenter  = "debug-tile-center"
)

func init() {
	provider.Register(Name, NewProvider)
}

//	NewProvider Setups a debug provider. there are not currently any config params supported
func NewProvider(config map[string]interface{}) (mvt.Provider, error) {
	return &Provider{}, nil
}

// Provider provides the debug provider
type Provider struct{}

func (p *Provider) MVTLayer(ctx context.Context, layerName string, tile tegola.Tile, dtags map[string]interface{}) (*mvt.Layer, error) {
	var layer mvt.Layer

	//	get tile bounding box
	ext := tile.BoundingBox()

	xlen := ext.Maxx - ext.Minx
	ylen := ext.Maxy - ext.Miny

	switch layerName {
	case "debug-tile-outline":
		//	debug outlines
		layer = mvt.Layer{
			Name:         LayerDebugTileOutline,
			DontSimplify: true,
		}
		debugOutline := mvt.Feature{
			Tags: map[string]interface{}{
				"type": "debug_outline",
			},
			Geometry: &basic.Line{ //	tile outline
				basic.Point{ext.Minx, ext.Miny},
				basic.Point{ext.Maxx, ext.Miny},
				basic.Point{ext.Maxx, ext.Maxy},
				basic.Point{ext.Minx, ext.Maxy},
			},
		}
		layer.AddFeatures(debugOutline)

	case "debug-tile-center":
		//	debug center points
		layer = mvt.Layer{
			Name:         LayerDebugTileCenter,
			DontSimplify: true,
		}
		debugCenter := mvt.Feature{
			Tags: map[string]interface{}{
				"type": "debug_text",
				"zxy":  fmt.Sprintf("Z:%v, X:%v, Y:%v", tile.Z, tile.X, tile.Y),
			},
			Geometry: &basic.Point{ //	middle of the tile
				ext.Minx + (xlen / 2),
				ext.Miny + (ylen / 2),
			},
		}
		layer.AddFeatures(debugCenter)
	}

	return &layer, nil
}

// Layers returns information about the various layers the provider supports
func (p *Provider) Layers() ([]mvt.LayerInfo, error) {
	layers := []Layer{
		{
			name:     "debug-tile-outline",
			geomType: basic.Line{},
			srid:     tegola.WebMercator,
		},
		{
			name:     "debug-tile-center",
			geomType: basic.Point{},
			srid:     tegola.WebMercator,
		},
	}

	var ls []mvt.LayerInfo

	for i := range layers {
		ls = append(ls, layers[i])
	}

	return ls, nil
}
