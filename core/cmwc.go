package core

import (
  "fmt"
  "math/big"
  "math/rand"
  "time"
  "github.com/runningwild/stringz"
)

func Check(c *SlowCMWC, max int) int {
  for i := range c.Q {
    c.Q[i] = uint64(source.Uint32() % c.B)
  }
  seq := make([]byte, max*3)
  for i := range seq {
    seq[i] = byte(c.Int63())
  }
  res := stringz.Find(seq[0:10]).In(seq)
  if len(res) <= 1 {
    return 0
  }
  return res[1]
}

// finds order of b in ab^r+1.  Assumes that b is a power of 2 and that
// a is prime.
func findOrder(a, b, r int64) *big.Int {
  A := big.NewInt(a)
  B := big.NewInt(b)
  Br := big.NewInt(b)
  Br.Exp(Br, big.NewInt(r), nil)
  m := big.NewInt(0)
  m.Mul(A, Br)
  one := big.NewInt(1)
  m.Add(m, one)

  // First check all powers of two
  p := big.NewInt(2)
  for p.Cmp(Br) < 0 {
    v := big.NewInt(0)
    v.Exp(B, p, m)
    if v.Cmp(one) == 0 {
      return p
    }
    p.Mul(p, big.NewInt(2))
  }

  // now check all a*2^n
  p = big.NewInt(a)
  for p.Cmp(m) < 0 {
    v := big.NewInt(0)
    v.Exp(B, p, m)
    if v.Cmp(one) == 0 {
      return p
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
      order := findOrder(A.Int64(), B.Int64(), R.Int64())
      if order != nil {
        if int64(int32(A.Int64())) != A.Int64() {
          fmt.Printf("WHAT!!!!\n")
        }
        return uint64(A.Int64()), order
      }
    }
  }
  return 0, nil
}

func MakeSlowCMWC(a, b, r uint32) *SlowCMWC {
  var c SlowCMWC
  c.Q = make([]uint64, int(r))
  c.C = 5
  c.A = uint64(a)
  c.B = b
  return &c
}

// Doesn't take advantage of mod b == 2^32, but does let us set arbitrary
// values of b for testing purposes.
type SlowCMWC struct {
  Q []uint64
  A uint64
  C uint32
  B uint32
  N int
}

func (c *SlowCMWC) Seed(seed int64) {
  for i := range c.Q {
    c.Q[i] = uint64(seed)
  }
}

func (c *SlowCMWC) Int63() int64 {
  c.N = (c.N + 1) % len(c.Q)
  C := c.C
  xn := uint64(c.B-1) - (c.A*c.Q[c.N] + uint64(C))
  xn = (xn % uint64(c.B))
  cn := (c.A*c.Q[c.N] + uint64(C)) / uint64(c.B)
  c.Q[c.N] = xn
  c.C = uint32(cn)
  return int64(xn)
}
