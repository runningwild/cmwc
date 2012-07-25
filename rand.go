package main

import (
  "fmt"
  "runningwild/rand/core"
  "math/big"
)

func main() {
  var a, b, r uint64
  b = uint64(1 << 31)
  r = uint64(32)
  var p *big.Int
  for i := 0; i < 10; i++ {
    a, p = core.GenerateRandomParams(b, r)
    fmt.Printf("CMWC(%v, %v, %v) -> %v\n", a, b, r, p)
    // p2 := core.Check(c, int(p.Int64()+1234))
    // fmt.Printf("Period: %v %v\n", p, p2)
  }
}
