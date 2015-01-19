package main

import (
  "fmt"
  "math"
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

func KMeans(points []*Vec, seeds []*Vec, iters uint) []*Vec {
  clusters := seeds
  for i := uint(0); i < iters; i++ {
    assignments := make([][]*Vec, len(clusters))
    // slices are annoying as hell, so let's just be lazy here
    for j := 0; j < len(assignments); j++ {
      assignments[i] = make([]*Vec, 0, len(points))
    }
  }
  return clusters
}

func main() {
  v := Vec{[]float64{1,2,3}}
  fmt.Println(v.add(&v))
  fmt.Println(v.norm())
  fmt.Println(ArgMin([]float64{1,2,3}))
}
