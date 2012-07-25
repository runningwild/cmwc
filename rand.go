// Simple package implementing Complementary Multiply With Carry random number
// generators.  These generators require that B=2^32 and that R is a power of
// two.  This generator is roughly twice as fast as the one in math/rand when
// being used as a rand.Source, and is also gobbable, in case you need to
// store the state of the rng to replay values later.
//
// A Cmwc requires 4*(R+4) bytes of storage.  Sample parameters are given here
// with their period lengths:
//
//     R=1 (lb2_r=0)
//     a = 4253332471, period = 285436310343122944
//     a = 4264505017, period = 286186087213170688
//     a = 4283145667, period = 287437040058892288
//     a = 4285374253, period = 287586597933678592
//     a = 4288278103, period = 287781472008404992
//
//     R=2 (lb2_r=1)
//     a = 4231509211, period = 1219649491575962978152873984
//     a = 4231512247, period = 1219650366643384974749728768
//     a = 4246691611, period = 1224025520438848726534979584
//     a = 4269449707, period = 1230585095009425492969259008
//     a = 4269838273, period = 1230697091533765258994778112
//
//     R=4 (lb2_r=2)
//     a = 4250569903, period = 22599906052433497083007582018186157588859584512
//     a = 4252386973, period = 22609567253690700671348291300966252011734433792
//     a = 4256424451, period = 22631034168850563433623090969454758271966511104
//     a = 4261847173, period = 22659866304433598515911669683172908474765934592
//     a = 4292669383, period = 22823745282129445683379567487218682232341266432
//
//     R=8 (lb2_r=3)
//     a = 4224759397, period = 7.644 x 10^84
//     a = 4250989063, period = 7.691 x 10^84
//     a = 4268111437, period = 7.722 x 10^84
//     a = 4270484551, period = 7.726 x 10^84
//     a = 4285415527, period = 7.753 x 10^84
//
//     R=16 (lb2_r=4)
//     a = 3864648517, period = 8.096 x 10^161
//     a = 4092063091, period = 8.573 x 10^161
//     a = 4116967117, period = 8.625 x 10^161
//     a = 4144055527, period = 8.682 x 10^161
//     a = 4270521133, period = 8.947 x 10^161
//
//     R=32 (lb2_r=5)
//     a = 1992756781, period = 5.597 x 10^315
//     a = 2392853653, period = 6.721 x 10^315
//     a = 2625435811, period = 7.375 x 10^315
//     a = 3372549937, period = 9.473 x 10^315
//     a = 3945340957, period = 1.108 x 10^316
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

// Creates a Cmwc RNG with (A, B, R) = (a, 2^32, 1<<lb2_r)
func MakeCmwc(a, lb2_r uint32) *Cmwc {
  return &Cmwc{core.MakeCMWC32(a, lb2_r)}
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (c *Cmwc) Int63() int64 {
  return c.cmwc.Int63()
}

// Seed uses the provided seed value to initialize the generator to a
// deterministic state.
func (c *Cmwc) Seed(seed int64) {
  c.cmwc.Seed(seed)
}

// Uses crypto.Reader to seed the generator.
func (c *Cmwc) SeedWithDevRand() {
  c.cmwc.SeedWithDevRand()
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
  r = uint64(1 << 5)
  var p *big.Int
  for i := 0; i < 100; i++ {
    a, p = core.GenerateRandomParams(b, r)
    fmt.Printf("CMWC(%v, %v, %v) -> %v\n", a, b, r, p)
    // c := core.MakeSlowCMWC(uint64(a), uint64(b), uint64(r))
    // p2 := core.Check(c, int(p.Int64()+1234))
    // fmt.Printf("Period: %v %v\n", p, p2)
  }
}
