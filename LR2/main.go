package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

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

	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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

	fc, err := geojson.UnmarshalFeatureCollection(rawGeometryJSON)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("/")
		fmt.Println(r.URL)
		fmt.Println(r.URL.Host)
		srs := strings.Split(r.URL.Path, "/")
		fmt.Println(srs)
		srs[3] = strings.Trim(srs[3], ".png")
		fmt.Println(srs)

		const S = 100
		dc := gg.NewContext(1366, 1024)
		for i := 0; i < len(fc.Features); i++ {

			z, _ := strconv.ParseFloat(srs[1], 64)
			x, _ := strconv.ParseFloat(srs[2], 64)
			y, _ := strconv.ParseFloat(srs[3], 64)
			l := 0.0
			if z != 0 {
				l = (512 / (z * 2))
			}

			if fc.Features[i].Properties["layer"] == "second" && z < 5 {
				continue
			}
			for r := 0; r < len(fc.Features[i].Geometry.MultiPolygon); r++ {
				for k := 0; k < len(fc.Features[i].Geometry.MultiPolygon[r]); k++ {
					dc.MoveTo(fc.Features[i].Geometry.MultiPolygon[r][k][0][0], fc.Features[i].Geometry.MultiPolygon[r][k][0][1])
					dc.Push()

					dc.ScaleAbout(z+1, z+1, x*l, y*l)
					for j := 1; j < len(fc.Features[i].Geometry.MultiPolygon[r][k]); j++ {
						dc.LineTo((fc.Features[i].Geometry.MultiPolygon[r][k][j][0]+256)/2, (fc.Features[i].Geometry.MultiPolygon[r][k][j][1]*(-1)+256)/2)
					}
				}
			}

			dc.SetLineWidth(10)
			fmt.Println(fc.Features[i].Properties["color"])
			switch fc.Features[i].Properties["color"] {
			case "green":
				dc.SetRGBA255(91, 255, 15, 255)
				dc.StrokePreserve()
				dc.SetRGBA255(91, 155, 15, 255)

			case "orange":
				dc.SetRGBA255(255, 184, 5, 255)
				dc.StrokePreserve()
				dc.SetRGBA255(200, 184, 5, 255)
			default:
				dc.SetRGBA255(255, 255, 255, 255)
				dc.StrokePreserve()
				dc.SetRGBA255(0, 0, 0, 255)
			}

			dc.Fill()
			dc.Pop()
		}
		dc.SavePNG("out.png")

		file.Close()

		photo, err := os.Open("out.png")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer photo.Close()

		stat, err = photo.Stat()
		if err != nil {
			return
		}
		data = make([]byte, stat.Size())
		for {
			_, err := photo.Read(data)
			if err == io.EOF { // если конец файла
				break // выходим из цикла
			}
		}
		w.Write(data)
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":8081", nil)
}
