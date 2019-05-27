package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fogleman/gg"
	geojson "github.com/paulmach/go.geojson"
)

func main() {

	file, err := os.Open("2_5467644889959236843.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	data := make([]byte, 64)
	raw := ""
	for {
		r, err := file.Read(data)
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
		raw = raw + string(data[:r])
	}

	rawGeometryJSON := []byte(raw)

	fc1, err := geojson.UnmarshalFeatureCollection(rawGeometryJSON)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	dc := gg.NewContext(1366, 1024)
	dc.SetHexColor("fff")

	dc.InvertY()
	dc.Scale(4, 4)
	for i := 0; i < len(fc1.Features); i++ {
		for r := 0; r < len(fc1.Features[i].Geometry.MultiPolygon); r++ {
			for k := 0; k < len(fc1.Features[i].Geometry.MultiPolygon[r]); k++ {
				dc.MoveTo(fc1.Features[i].Geometry.MultiPolygon[r][k][0][0], fc1.Features[i].Geometry.MultiPolygon[r][k][0][1])
				for j := 1; j < len(fc1.Features[i].Geometry.MultiPolygon[r][k]); j++ {
					dc.LineTo(fc1.Features[i].Geometry.MultiPolygon[r][k][j][0], fc1.Features[i].Geometry.MultiPolygon[r][k][j][1])
				}
			}
		}
	}

	dc.SetHexColor("f00")
	dc.Fill()
	dc.SavePNG("out.png")

}
