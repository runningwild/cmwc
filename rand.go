package rand

import (
  "fmt"
  "runningwild/rand/core"
  "math/big"
)

type Cmwc struct {
  cmwc *core.CMWC32
}

func MakeCmwc(a, lb2_r uint32) *Cmwc {
  return &Cmwc{core.MakeCMWC32(a, lb2_r)}
}
func (c *Cmwc) Int63() int64 {
  return c.cmwc.Int63()
}
func (c *Cmwc) Seed(seed int64) {
  c.cmwc.Seed(seed)
}

func main() {
  var a, b, r uint64
  b = uint64(1 << 32)
  r = uint64(1 << 4)
  var p *big.Int
  for i := 0; i < 10; i++ {
    a, p = core.GenerateRandomParams(b, r)
    fmt.Printf("CMWC(%v, %v, %v) -> %v\n", a, b, r, p)
    // c := core.MakeSlowCMWC(uint64(a), uint64(b), uint64(r))
    // p2 := core.Check(c, int(p.Int64()+1234))
    // fmt.Printf("Period: %v %v\n", p, p2)
  }
}
