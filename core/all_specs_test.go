package core_test

import (
  "testing"
  "github.com/orfjackal/gospec/src/gospec"
)

// List of all specs here
func TestAllSpecs(t *testing.T) {
  r := gospec.NewRunner()
  r.AddSpec(CMWCSpec)
  r.AddSpec(SlowCMWCSpec)
  r.AddSpec(CMWCGobSpec)
  gospec.MainGoTest(r, t)
}
