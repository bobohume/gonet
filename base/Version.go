package base

import (
	"fmt"
)

const (
	BUILD_NO = "1,5,1,1"
)

type (
	Version struct {
		m_sBuildVer int64
	}

	IVersion interface {
		IsAcceptableBuildVersion(string) bool
		Init()
	}
)

func (v *Version) Init() {
	var bv1, bv2, bv3, bv4 int
	fmt.Sscanf(BUILD_NO, "%d,%d,%d,%d", &bv1, &bv2, &bv3, &bv4)
	v.m_sBuildVer = int64(bv1*1000*1000*1000 + bv2*1000*1000 + bv3*1000 + bv4)
}

func (v *Version) IsAcceptableBuildVersion(version string) bool {
	var _v [4]int
	var sClient int64
	fmt.Sscanf(version, "%d,%d,%d,%d", &_v[0], &_v[1], &_v[2], &_v[3])
	sClient = int64(_v[0]*1000*1000*1000 + _v[1]*1000*1000 + _v[2]*1000 + _v[3])
	return sClient >= v.m_sBuildVer
}

var (
	VERSION Version
)

func init() {
	VERSION.Init()
}
