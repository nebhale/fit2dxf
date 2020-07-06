package main

import (
	"fmt"
	"github.com/tormoder/fit"
	"github.com/yofu/dxf"
	"math"
	"os"
)

const (
	AltitudeScale   = 50
	MinimumAltitude = 50000
)

func main() {
	in, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer in.Close()

	out, err := os.OpenFile(os.Args[2], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	f, err := fit.Decode(in)
	if err != nil {
		panic(err)
	}

	a, err := f.Activity()
	if err != nil {
		panic(err)
	}

	v := [][]float64{
		{a.Records[0].GetDistanceScaled(), MinimumAltitude},
		{a.Records[0].GetDistanceScaled(), a.Records[0].GetAltitudeScaled() * AltitudeScale},
	}
	for _, r := range a.Records[1:] {
		d := r.GetDistanceScaled()
		a1 := v[len(v)-1][1]
		a2 := r.GetAltitudeScaled() * AltitudeScale

		v = append(v, []float64{math.Sqrt(math.Pow(d, 2) - math.Pow(a2-a1, 2)), a2})
	}
	v = append(v, []float64{v[len(v)-1][0], MinimumAltitude})

	fmt.Printf("%d Points\n", len(v))
	fmt.Printf("%f,%f => %f,%f", v[1][0], v[1][1], v[len(v)-2][0], v[len(v)-2][1])

	d := dxf.NewDrawing()

	if _, err := d.LwPolyline(true, v...); err != nil {
		panic(err)
	}

	if _, err := d.WriteTo(out); err != nil {
		panic(err)
	}
}
