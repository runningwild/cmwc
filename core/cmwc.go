package core

import (
  crand "crypto/rand"
  "fmt"
  "github.com/runningwild/stringz"
  "math/big"
  "math/rand"
  "time"
)

func Check(c *SlowCMWC, max int) int {
  for i := range c.Q {
    c.Q[i] = uint64(uint64(source.Uint32()) % c.B)
  }
  seq := make([]byte, max*3)
  for i := range seq {
    seq[i] = byte(c.Next())
  }
  res := stringz.Find(seq[0:10]).In(seq)
  if len(res) <= 1 {
    return 0
  }
  return res[1]
}

// finds order of b in ab^r+1.  Assumes that b is a power of 2 and that
// a is prime.
// 2^guess is a lower bound on the return value
func findOrder(a, b, r, guess int64) *big.Int {
  A := big.NewInt(a)
  B := big.NewInt(b)
  Br := big.NewInt(b)
  Br.Exp(Br, big.NewInt(r), nil)
  m := big.NewInt(0)
  m.Mul(A, Br)
  one := big.NewInt(1)
  m.Add(m, one)

  lower := big.NewInt(0)
  lower.Exp(big.NewInt(2), big.NewInt(guess*32), nil)
  higher := big.NewInt(0)
  higher.Exp(big.NewInt(2), big.NewInt((r+8)*32), nil)
  // First check all powers of two
  p := big.NewInt(2)
  for p.Cmp(Br) < 0 && p.Cmp(higher) < 0 {
    if p.Cmp(lower) > 0 {
      v := big.NewInt(0)
      v.Exp(B, p, m)
      if v.Cmp(one) == 0 {
        return p
      }
    }
    p.Mul(p, big.NewInt(2))
  }

  // now check all a*2^n
  p = big.NewInt(a)
  for p.Cmp(m) < 0 && p.Cmp(higher) < 0 {
    if p.Cmp(lower) > 0 {
      v := big.NewInt(0)
      v.Exp(B, p, m)
      if v.Cmp(one) == 0 {
        return p
      }
    }
    p.Mul(p, big.NewInt(2))
  }
  return nil
}

var source *rand.Rand

func init() {
  source = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Generates parameters for a SlowCMWC PNRG.
func GenerateRandomParams(b, r uint64) (a uint64, period *big.Int) {
  B := big.NewInt(int64(b))
  R := big.NewInt(int64(r))
  v := big.NewInt(0)
  v.Exp(B, R, nil)
  one := big.NewInt(1)
  for i := 0; i < 1000000; i++ {
    A := big.NewInt(0)
    A.Rand(source, B)
    if !A.ProbablyPrime(20) {
      continue
    }
    p := big.NewInt(0)
    p.Mul(A, v)
    p.Add(p, one)
    if p.ProbablyPrime(20) {
      order := findOrder(A.Int64(), B.Int64(), R.Int64(), R.Int64())
      if order != nil {
        if int64(uint32(A.Int64())) != A.Int64() {
          fmt.Printf("WHAT!!!!\n")
        }
        return uint64(A.Int64()), order
      }
    }
  }
  return 0, nil
}

func MakeCMWC32(a, lb2_r uint32) *CMWC32 {
  r := uint32(1) << lb2_r
  c := CMWC32{
    Q:      make([]uint32, int(r)),
    A:      a,
    R_mask: r - 1,
    C:      0,
  }
  for i := range c.Q {
    c.Q[i] = uint32(0)
  }
  return &c
}

type CMWC32 struct {
  Q []uint32
  A uint32
  C uint32
  N uint32

  R_mask uint32
}

func (c *CMWC32) Next() uint32 {
  c.N = (c.N + 1) & c.R_mask
  xp := uint64(c.Q[int(c.N)])
  ax_c := uint64(c.A)*xp + uint64(c.C)
  xn := uint32(0xffffffff - ax_c)
  cn := uint32(ax_c >> 32)
  c.Q[c.N] = xn
  c.C = cn
  return xn
}

func (c *CMWC32) Int63() int64 {
  return int64((uint64(c.Next())<<32 | uint64(c.Next())) & 0x7fffffffffffffff)
}

func (c *CMWC32) Seed(seed int64) {
  var foo uint32 = 0x3596ac35
  for i := 0; i < len(c.Q)-3; i += 4 {
    c.Q[i] = uint32(seed & 0xffffffff)
    c.Q[i+1] = uint32(seed >> 32)
    c.Q[i+2] = c.Q[i] ^ foo
    c.Q[i+3] = c.Q[i+1] ^ foo
    foo = (foo >> 3) | (foo << 29)
  }
  for i := len(c.Q) - len(c.Q)%4; i < len(c.Q); i++ {
    c.Q[i] = foo
    foo = (foo >> 3) | (foo << 29)
  }
  c.C = uint32((seed & 0xffffffff) ^ (seed >> 32))
  c.N = 0
  for i := 0; i < len(c.Q)*10; i++ {
    c.Next()
  }
}

func (c *CMWC32) SeedWithDevRand() {
  buf := make([]byte, len(c.Q)*4)
  crand.Reader.Read(buf)
  for i := range c.Q {
    for j := 0; j < 4; j++ {
      c.Q[i] |= uint32(buf[i*4+j]) << uint(4*j)
    }
  }
  buf = buf[0:4]
  crand.Reader.Read(buf)
  for i := 0; i < 4; i++ {
    c.C |= uint32(buf[i]) << uint(i*8)
  }
}

func MakeSlowCMWC(a, b, r uint64) *SlowCMWC {
  var c SlowCMWC
  c.Q = make([]uint64, int(r))
  c.C = 0
  c.A = a
  c.B = b
  return &c
}

// Doesn't take advantage of mod b == 2^32, but does let us set arbitrary
// values of b for testing purposes.
type SlowCMWC struct {
  Q []uint64
  A uint64
  C uint64
  N int
  B uint64
}

func (c *SlowCMWC) Seed(seed int64) {
  for i := range c.Q {
    c.Q[i] = uint64(seed)
  }
}

func (c *SlowCMWC) Next() uint32 {
  c.N = (c.N + 1) % len(c.Q)
  px := c.Q[c.N]
  C := c.C
  xn := c.B - 1 - (c.A*px + C)
  xn = xn % c.B
  cn := (c.A*px + uint64(C)) / c.B
  c.Q[c.N] = xn
  c.C = cn
  return uint32(xn)
}
