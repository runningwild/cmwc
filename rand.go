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
//     a = 4287701437, period = 1235845798012744974758576128 
//     a = 4289233213, period = 1236287302385405139166953472 
//     a = 4289388097, period = 1236331944658985020888711168 
//     a = 4290242227, period = 1236578130870167482440613888 
//     a = 4292893333, period = 1237342260149765542355402752 
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
//
//     R=64 (lb2_r=6)
//     a = 2632744897, period = 6.647 x 10^623
//     a = 2103873433, period = 1.062 x 10^624
//
//     R=128 (lb2_r=7)
//     a = 2307597247, period = 4.707 x 10^1239
//     a = 1410304387, period = 2.301 x 10^1240
//
//     R=256 (lb2_r=8)
//     a = 289655593,  period = 4.937 x 10^2472
//     a = 3918333073, period = 1.670 x 10^2473
//
//     R=512 (lb2_r=9)
//     a = 2428121623, period = 4.514 x 10^4939
package cmwc

import (
  "bytes"
  "encoding/binary"
  "fmt"
  "github.com/runningwild/cmwc/core"
  "math/big"
)

type Cmwc struct {
  cmwc *core.CMWC32
}

// Creates a Cmwc RNG with (A, B, R) = (a, 2^32, 1<<lb2_r)
func MakeCmwc(a, lb2_r uint32) *Cmwc {
  return &Cmwc{core.MakeCMWC32(a, lb2_r)}
}

// Creates a Cmwc RNG with a very long period
func MakeGoodCmwc() *Cmwc {
  return &Cmwc{core.MakeCMWC32(3945340957, 5)}
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

func (c *Cmwc) OverwriteWith(c2 *Cmwc) {
  if len(c.cmwc.Q) != len(c2.cmwc.Q) {
    panic("Cannot overwrite with a generator with a different size.")
  }
  c.cmwc.A = c2.cmwc.A
  c.cmwc.C = c2.cmwc.C
  c.cmwc.N = c2.cmwc.N
  for i := range c.cmwc.Q {
    c.cmwc.Q[i] = c2.cmwc.Q[i]
  }
}

func main() {
  var a, b, r uint64
  b = uint64(1 << 32)
  r = uint64(1 << 1)
  var p *big.Int
  for i := 0; i < 10000; i++ {
    a, p = core.GenerateRandomParams(b, r)
    // fmt.Printf("CMWC(%v, %v, %v) -> %v\n", a, b, r, p)
    fmt.Printf("%v %v\n", p, a)
    // c := core.MakeSlowCMWC(uint64(a), uint64(b), uint64(r))
    // p2 := core.Check(c, int(p.Int64()+1234))
    // fmt.Printf("Period: %v %v\n", p, p2)
  }
}
