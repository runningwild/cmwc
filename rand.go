package cmwc

import (
  "fmt"
  "math/big"
  "bytes"
  "encoding/binary"
  "github.com/runningwild/cmwc/core"
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
func (c *Cmwc) GobEncode() ([]byte, error) {
  buf := bytes.NewBuffer(make([]byte, 4*(len(c.cmwc.Q)+4))[0:0])
  binary.Write(buf, binary.LittleEndian, c.cmwc.A)
  binary.Write(buf, binary.LittleEndian, c.cmwc.C)
  binary.Write(buf, binary.LittleEndian, c.cmwc.N)
  binary.Write(buf, binary.LittleEndian, uint32(len(c.cmwc.Q)))
  err := binary.Write(buf, binary.LittleEndian, c.cmwc.Q)
  if err != nil {
    return nil, err
  }
  return buf.Bytes(), nil
}
func (c *Cmwc) GobDecode(data []byte) error {
  c.cmwc = &core.CMWC32{}
  buf := bytes.NewBuffer(data)
  binary.Read(buf, binary.LittleEndian, &c.cmwc.A)
  binary.Read(buf, binary.LittleEndian, &c.cmwc.C)
  binary.Read(buf, binary.LittleEndian, &c.cmwc.N)
  var length uint32
  err := binary.Read(buf, binary.LittleEndian, &length)
  if err != nil {
    return nil
  }
  c.cmwc.Q = make([]uint32, length)
  err = binary.Read(buf, binary.LittleEndian, c.cmwc.Q)
  c.cmwc.R_mask = uint32(len(c.cmwc.Q)) - 1
  return err
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
