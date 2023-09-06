package core

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"memkv/internal/data_structure"
)

func resetSetStore() {
	setStore = make(map[string]data_structure.Set)
}

func TestEvalSADD(t *testing.T) {
	resetSetStore()
	res, err := Decode(evalSADD([]string{"set", "adele"}))
	assert.Nil(t, err)
	assert.EqualValues(t, 1, res)

	res, err = Decode(evalSADD([]string{"set", "adele", "bob", "chris"}))
	assert.Nil(t, err)
	assert.EqualValues(t, 2, res)
}

func TestEvalSREM(t *testing.T) {
	resetSetStore()
	res, err := Decode(evalSREM([]string{"set", "adele"}))
	assert.Nil(t, err)
	assert.EqualValues(t, 0, res)

	evalSADD([]string{"set", "a", "b", "c"})
	res, err = Decode(evalSREM([]string{"set", "a", "d"}))
	assert.Nil(t, err)
	assert.EqualValues(t, 1, res)
}

func TestEvalSCARD(t *testing.T) {
	resetSetStore()

	evalSADD([]string{"set", "a", "b", "c"})
	res, err := Decode(evalSCARD([]string{"set"}))
	assert.Nil(t, err)
	assert.EqualValues(t, 3, res)
}

func TestEvalSMEMBERS(t *testing.T) {
	resetSetStore()

	evalSADD([]string{"set", "a", "b", "c"})
	res, err := Decode(evalSMEMBERS([]string{"set"}))
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"a", "b", "c"}, res)
}

func TestEvalSMISMEMBER(t *testing.T) {
	resetSetStore()

	evalSADD([]string{"set", "a", "b", "c"})
	res, err := Decode(evalSMISMEMBER([]string{"set", "a", "d"}))
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1, 0}, res)
}

func TestEvalSRAND(t *testing.T) {
	resetSetStore()

	evalSADD([]string{"set", "a", "b", "c"})
	res, err := Decode(evalSRAND([]string{"set", "2"}))

	assert.Nil(t, err)
	m := make(map[string]struct{})
	m["a"] = struct{}{}
	m["b"] = struct{}{}
	m["c"] = struct{}{}
	rd := make(map[string]struct{})
	for _, key := range res.([]interface{}) {
		k := key.(string)
		assert.Contains(t, m, k, "key must be in set")
		assert.NotContains(t, m, rd, "key must be not duplicated")
		rd[k] = struct{}{}
	}
}

func TestEvalSPOP(t *testing.T) {
	resetSetStore()

	evalSADD([]string{"set", "a", "b", "c"})
	res, err := Decode(evalSPOP([]string{"set", "2"}))

	assert.Nil(t, err)
	m := make(map[string]struct{})
	m["a"] = struct{}{}
	m["b"] = struct{}{}
	m["c"] = struct{}{}
	for _, key := range res.([]interface{}) {
		k := key.(string)
		delete(m, k)
	}
	var expected []string
	for k := range m {
		expected = append(expected, k)
	}

	res, err = Decode(evalSMEMBERS([]string{"set"}))
	assert.ElementsMatch(t, expected, res)
}

func TestEvalGEOADD(t *testing.T) {
	delete(zsetStore, "vn")
	res, err := Decode(evalGEOADD([]string{"vn", "10", "20", "p1"}))
	assert.Nil(t, err)
	assert.EqualValues(t, res, 1)

	res, err = Decode(evalGEOADD([]string{"vn", "10", "20", "p1"}))
	assert.Nil(t, err)
	assert.EqualValues(t, res, 0)

	res, err = Decode(evalGEOADD([]string{"vn", "-10", "20", "p1"}))
	assert.Nil(t, err)
	assert.EqualValues(t, res, 1)

	res, err = Decode(evalGEOADD([]string{"vn", "-10", "20", "p2", "-1", "2", "p3"}))
	assert.Nil(t, err)
	assert.EqualValues(t, res, 2)

	zset, exist := zsetStore["vn"]
	assert.True(t, exist)
	assert.EqualValues(t, 3, zset.Len())

	res, err = Decode(evalGEOADD([]string{"vn"}))
	assert.EqualValues(t, "(error) ERR wrong number of arguments for 'GEOADD' command", res)
	res, err = Decode(evalGEOADD([]string{"vn", "-10", "20", "p4", "20"}))
	assert.EqualValues(t, "(error) ERR wrong number of arguments for 'GEOADD' command", res)
}

func TestEvalGEODIST(t *testing.T) {
	delete(zsetStore, "vn")
	evalGEOADD([]string{"vn", "20", "10", "p1"})
	evalGEOADD([]string{"vn", "40", "30", "p2"})
	evalGEOADD([]string{"vn", "10", "85", "p3"})
	evalGEOADD([]string{"vn", "10", "-85", "p4"})
	evalGEOADD([]string{"vn", "180", "20", "p5"})
	evalGEOADD([]string{"vn", "179.9999", "20", "p6"})
	res, err := Decode(evalGEODIST([]string{"vn", "p1", "p2"}))
	assert.Nil(t, err)
	dist, err := strconv.ParseFloat(res.(string), 64)
	assert.Nil(t, err)
	assert.LessOrEqual(t, math.Abs(dist-3041460.716138), 1.0)

	res, err = Decode(evalGEODIST([]string{"vn", "p3", "p4"}))
	assert.Nil(t, err)
	dist, err = strconv.ParseFloat(res.(string), 64)
	assert.Nil(t, err)
	assert.LessOrEqual(t, math.Abs(dist-18908471), 1.0)

	res, err = Decode(evalGEODIST([]string{"vn", "p5", "p6"}))
	assert.Nil(t, err)
	dist, err = strconv.ParseFloat(res.(string), 64)
	assert.Nil(t, err)
	assert.LessOrEqual(t, math.Abs(dist-10.451853), 1.0)

	res, err = Decode(evalGEODIST([]string{"vn", "p1", "p2", "km"}))
	assert.Nil(t, err)
	dist, err = strconv.ParseFloat(res.(string), 64)
	assert.Nil(t, err)
	assert.LessOrEqual(t, math.Abs(dist-3041), 1.0)
}

func TestEvalGeoHash(t *testing.T) {
	delete(zsetStore, "vn")
	evalGEOADD([]string{"vn", "13.361389", "38.115556", "p1"})
	evalGEOADD([]string{"vn", "15.087269", "37.502669", "p2"})
	evalGEOADD([]string{"vn", "100", "80", "p3"})
	evalGEOADD([]string{"vn", "40", "-20", "p4"})
	evalGEOADD([]string{"vn", "-20", "39", "p5"})
	ret, err := Decode(evalGEOHASH([]string{"vn", "p1", "p2", "p3", "p4", "p5", "p6"}))
	expected := []string{"sqc8b49rny0", "sqdtr74hyu0", "ynpp5e9cbc0", "kukqnpp5e90", "ewcvbgsrqn0", ""}
	assert.Nil(t, err)
	assert.ElementsMatch(t, expected, ret)

	ret, err = Decode(evalGEOHASH([]string{"not_exist"}))
	assert.Nil(t, err)
	expected = []string{}
	assert.ElementsMatch(t, expected, ret)
}

func TestSimpleEvalGEOSEARCH(t *testing.T) {
	delete(zsetStore, "nyc")
	evalGEOADD([]string{"nyc", "-73.9733487", "40.7648057", "central park"})
	evalGEOADD([]string{"nyc", "-73.9903085", "40.7362513", "union square"})
	evalGEOADD([]string{"nyc", "-74.0131604", "40.7126674", "wtc one"})
	evalGEOADD([]string{"nyc", "-73.7858139", "40.6428986", "jfk"})
	evalGEOADD([]string{"nyc", "-73.9375699", "40.7498929", "q4"})
	evalGEOADD([]string{"nyc", "-73.9564142", "40.7480973", "4545"})

	ret, err := Decode(evalGEOSEARCH([]string{"nyc", "FROMLONLAT", "-73.9798091", "40.7598464", "3000"}))
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"central park", "4545", "union square"}, ret)

	evalGEOADD([]string{"nyc", "-73.9798091", "40.7598464", "me"})
	ret, err = Decode(evalGEOSEARCH([]string{"nyc", "FROMMEMBER", "me", "3000"}))
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"me", "central park", "4545", "union square"}, ret)
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func TestRandomEvalGEOSEARCH(t *testing.T) {
	delete(zsetStore, "nyc")
	targetLon := -73.9798091
	targetLat := 40.7598464
	for round := 0; round < 10; round++ {
		var expected []string
		radius := randFloat(1000.0, 2000000.0)
		for i := 0; i < 10000; i++ {
			lon := randFloat(-150, 150)
			lat := randFloat(-45, 45)
			name := fmt.Sprintf("%d", i)
			evalGEOADD([]string{"nyc",
				fmt.Sprintf("%f", lon),
				fmt.Sprintf("%f", lat),
				name})
			dist := data_structure.GeohashGetDistance(targetLon, targetLat, lon, lat)
			if dist <= radius {
				expected = append(expected, name)
			}
		}

		ret, err := Decode(evalGEOSEARCH([]string{"nyc", "FROMLONLAT",
			fmt.Sprintf("%f", targetLon),
			fmt.Sprintf("%f", targetLat),
			fmt.Sprintf("%f", radius)}))

		assert.Nil(t, err)
		assert.ElementsMatch(t, expected, ret)
	}
}

func TestEvalGEOPOS(t *testing.T) {
	delete(zsetStore, "nyc")
	evalGEOADD([]string{"nyc", "-73.9733487", "40.7648057", "central park"})
	evalGEOADD([]string{"nyc", "-73.9375699", "40.7498929", "q4"})
	ret, err := Decode(evalGEOPOS([]string{"nyc", "x"}))
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(ret.([]interface{})))
	assert.EqualValues(t, 0, len(ret.([]interface{})[0].([]interface{})))

	ret, err = Decode(evalGEOPOS([]string{"nyc", "central park", "q4"}))
	fmt.Println(ret)
	long1 := ret.([]interface{})[0].([]interface{})[0].(string)
	lat1 := ret.([]interface{})[0].([]interface{})[1].(string)
	assert.EqualValues(t, "-73.973348", long1)
	assert.EqualValues(t, "40.764806", lat1)

	long2 := ret.([]interface{})[1].([]interface{})[0].(string)
	lat2 := ret.([]interface{})[1].([]interface{})[1].(string)
	assert.EqualValues(t, "-73.937573", long2)
	assert.EqualValues(t, "40.749892", lat2)
}
