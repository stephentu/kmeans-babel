package main

import (
  "fmt"
  "math"
  "strings"
  "strconv"
  "os"
  "io/ioutil"
)

// this horrible language has no generics

type Vec struct {
  pts []float64
}

func (v *Vec) dot(that *Vec) float64 {
  var acc float64 = 0.0
  for i, el := range v.pts {
    acc += el * that.pts[i]
  }
  return acc
}

func (v *Vec) add(that *Vec) Vec {
  s := make([]float64, len(v.pts))
  for i:= 0; i < len(s); i++ {
    s[i] = v.pts[i] + that.pts[i]
  }
  return Vec{s}
}

func (v *Vec) sub(that *Vec) Vec {
  s := make([]float64, len(v.pts))
  for i:= 0; i < len(s); i++ {
    s[i] = v.pts[i] - that.pts[i]
  }
  return Vec{s}
}

func (v *Vec) mult(a float64) Vec {
  s := make([]float64, len(v.pts))
  for i:= 0; i < len(s); i++ {
    s[i] = v.pts[i] * a
  }
  return Vec{s}
}

func (v *Vec) norm() float64 {
  return math.Sqrt(v.dot(v))
}

func ArgMin(v []float64) uint {
  val := v[0]
  var idx uint = 0
  for i := uint(1); i < uint(len(v)); i++ {
    if v[i] < val {
      val = v[i]
      idx = i
    }
  }
  return idx
}

func Dists(v *Vec, centers []*Vec) []float64 {
  // ugh no functional operators
  d := make([]float64, len(centers))
  for i, el := range centers {
    vv := v.sub(el)
    d[i] = vv.norm()
  }
  return d
}

func Which(v *Vec, centers []*Vec) uint {
  return ArgMin(Dists(v, centers));
}

func Centroid(pts []*Vec) *Vec {
  v := *pts[0]
  for _, elem := range pts[1:] {
    v = v.add(elem)
  }
  v = v.mult(1. / float64(len(pts)))
  return &v
}

func KMeans(points []*Vec, seeds []*Vec, iters uint) []*Vec {
  clusters := seeds
  for i := uint(0); i < iters; i++ {
    assignments := make([][]*Vec, len(clusters))
    for idx := 0; idx < len(clusters); idx++ {
      assignments[idx] = make([]*Vec, 0)
    }
    for _, elem := range points {
      idx := Which(elem, clusters)
      assignments[idx] = append(assignments[idx], elem)
    }
    for idx, pts := range assignments {
      if len(pts) > 0 {
        clusters[idx] = Centroid(pts)
      }
    }
  }
  return clusters
}

func Parse(s []string) []*Vec {
  v := make([]*Vec, 0, len(s))
  for _, line := range s {
    toks := strings.Split(strings.TrimSpace(line), " ")
    pts := make([]float64, 0, len(toks))
    for _, tok := range toks {
      f, _ := strconv.ParseFloat(tok, 64)
      // ignore error
      pts = append(pts, f)
    }
    v = append(v, &Vec{pts})
  }
  return v
}

func main() {
  pointsFile := os.Args[1]
  seedsFile := os.Args[2]

  pointsRaw, _ := ioutil.ReadFile(pointsFile)
  seedsRaw, _ := ioutil.ReadFile(seedsFile)

  points := Parse(strings.Split(strings.TrimSpace(string(pointsRaw)), "\n"))
  seeds := Parse(strings.Split(strings.TrimSpace(string(seedsRaw)), "\n"))

  for _, elem := range KMeans(points, seeds, 20) {
    fmt.Println(elem)
  }
}
