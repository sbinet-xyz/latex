// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

import (
	"fmt"
	"testing"
)

func TestTTFBackend(t *testing.T) {
	be := NewTTFBackend()
	for _, sym := range []string{
		//	"A",
		//	"B",
		//	"a",
		//	"g",
		//	"z",
		//	"Z",
		//	"I",
		//	"T",
		//	"i",
		//	"t",
		`\sigma`,
	} {
		for _, math := range []bool{
			true,
			//	false,
		} {
			for _, descr := range []Font{
				{Name: "regular", Size: 12, Type: "regular"},
				{Name: "regular", Size: 10, Type: "regular"},
				//		{Name: "regular", Size: 12, Type: "it"},
				//		{Name: "regular", Size: 10, Type: "it"},
			} {
				t.Run(fmt.Sprintf("%s-math=%v-%s-%g-%s", sym, math, descr.Name, descr.Size, descr.Type), func(t *testing.T) {
					me := be.Metrics(sym, descr, 72, math)
					//t.Fatalf("metrics= %#v", me)
					fmt.Printf("[%s:%g:%s][%s]: adv=%g h=%g, w=%g. x:%g,%g, y=%g,%g, ice:%g, slanted:%v\n",
						descr.Name, descr.Size, descr.Type,
						sym,
						me.Advance, me.Height, me.Width,
						me.XMin, me.XMax, me.YMin, me.YMax,
						me.Iceberg, me.Slanted,
					)
				})
			}
		}
	}
	t.Fatalf("boo")
}
