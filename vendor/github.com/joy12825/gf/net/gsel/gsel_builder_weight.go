// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/joy12825/gf.

package gsel

type builderWeight struct{}

func NewBuilderWeight() Builder {
	return &builderWeight{}
}

func (*builderWeight) Name() string {
	return "BalancerWeight"
}

func (*builderWeight) Build() Selector {
	return NewSelectorWeight()
}
