package main

import (
	"fmt"
	"math"
)

func getSign(lng float64) float64 {
	if lng < 0 {
		return -1
	}else{
		return 1
	}
}

func adjustLng(lng float64) float64 {
	if math.Abs(lng) < 3.141592653589793 {
		return lng
	}else{
		return lng - (getSign(lng) * 6.283185307179586)
	}
}

// 台湾二度分带(TM2)坐标格式转twd97/wgs84，参考https://www.sunriver.com.tw/taiwanmap/grid_tm2_convert.php
func getWgsPoint(x, y float64) (lng, lat float64) {
	x = x - 250000
	con := y / 0.9999 / 6378137
	phi := con
	maxIter := 6
	for i := 0; true; i++ {
		deltaPhi := ((con + 0.0025146070728447817*math.Sin(2.0*phi) - 0.000002639046620230819*math.Sin(4.0*phi) + 3.418046136775059e-9*math.Sin(6.0*phi)) / 0.9983242984445848) - phi
		phi += deltaPhi
		if math.Abs(deltaPhi) <= 1e-10 {
			break
		}
		if i >= maxIter {
			return 0, 0
		}
	}

	if math.Abs(phi) < 1.5707963267948966 {
		var sin_phi=math.Sin(phi)
		var cos_phi=math.Cos(phi)
		var tan_phi = math.Tan(phi)
		var c = 0.006739496775478856 * math.Pow(cos_phi,2)
		var cs = math.Pow(c,2)
		var t = math.Pow(tan_phi,2)
		var ts = math.Pow(t,2)
		con = 1.0 - 0.006694380022900686 * math.Pow(sin_phi,2)
		var n = 6378137 / math.Sqrt(con)
		var r = n * (1.0 - 0.006694380022900686) / con
		var d = x / (n * 0.9999)
		var ds = math.Pow(d,2)
		lat = phi - (n * tan_phi * ds / r) * (0.5 - ds / 24.0 * (5.0 + 3.0 * t + 10.0 * c - 4.0 * cs - 9.0 * 0.006739496775478856 - ds / 30.0 * (61.0 + 90.0 * t + 298.0 * c + 45.0 * ts - 252.0 * 0.006739496775478856 - 3.0 * cs)))
		lng = adjustLng(2.111848394913139 + (d * (1.0 - ds / 6.0 * (1.0 + 2.0 * t + c - ds / 20.0 * (5.0 - 2.0 * c + 28.0 * t - 3.0 * cs + 8.0 * 0.006739496775478856 + 24.0 * ts))) / cos_phi))
	} else {
		lat = 1.5707963267948966 * getSign(y)
		lng = 2.111848394913139
	}
	return lng * 57.29577951308232, lat * 57.29577951308232
}

func main() {
	//getPoint(x, y)
	lng, lat := getWgsPoint(179312.12, 2548798.25) // 臺南市新化區忠孝路2號 國立新化高中
	fmt.Println(lng, lat)
}
