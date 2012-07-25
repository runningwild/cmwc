package core_test

import (
  . "github.com/orfjackal/gospec/src/gospec"
  "github.com/orfjackal/gospec/src/gospec"
  "testing"
  "runningwild/rand/core"
  rrand "runningwild/rand"
  "math/rand"
)

func SlowCMWCSpec(c gospec.Context) {
  c.Specify("SlowCMWC produces sequences with the expected period.", func() {
    c.Specify("(A, B, R, Period) = (3, 16, 2, 96)", func() {
      gen := core.MakeSlowCMWC(3, 16, 2)
      period := core.Check(gen, 500)
      c.Expect(period, Equals, 96)
    })
    c.Specify("(A, B, R, Period) = (13, 16, 2, 416)", func() {
      gen := core.MakeSlowCMWC(13, 16, 2)
      period := core.Check(gen, 500)
      c.Expect(period, Equals, 416)
    })
    c.Specify("(A, B, R, Period) = (37, 256, 2, 128)", func() {
      gen := core.MakeSlowCMWC(37, 256, 2)
      period := core.Check(gen, 500)
      c.Expect(period, Equals, 128)
    })
    c.Specify("(A, B, R, Period) = (103, 256, 2, 52736)", func() {
      gen := core.MakeSlowCMWC(103, 256, 2)
      period := core.Check(gen, 100000)
      c.Expect(period, Equals, 52736)
    })
  })
}

func CMWCSpec(c gospec.Context) {
  c.Specify("CMWC32 produces the same output as SlowCMWC.", func() {
    fast := core.MakeCMWC32(11, 4)
    slow := core.MakeSlowCMWC(11, 1<<32, 1<<4)
    N := 100000
    for i := 0; i < N; i++ {
      f := fast.Next()
      s := slow.Next()
      c.Expect(f, Equals, s)
      if f != s {
        break
      }
    }
  })
}

func BenchmarkCMWC32Next(b *testing.B) {
  b.StopTimer()
  c := core.MakeCMWC32(3278470471, 4)
  b.StartTimer()
  for i := 0; i < b.N; i++ {
    c.Next()
  }
}

func BenchmarkCMWC32AsRandSource(b *testing.B) {
  b.StopTimer()
  c := rrand.MakeCmwc(3278470471, 4)
  r := rand.New(c)
  b.StartTimer()
  for i := 0; i < b.N; i++ {
    r.Int63()
  }
}

func BenchmarkStdRand(b *testing.B) {
  for i := 0; i < b.N; i++ {
    rand.Int63()
  }
}
