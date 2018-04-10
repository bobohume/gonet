package base

import (
	"fmt"
)

const(
	BUILD_NO = "1,5,1,1"
)

var(
	g_Version CVersion
)

type (
	CVersion struct{
		m_sBuildVer	int64
		m_bInit		bool
	}

	ICVersion interface {
		IsAcceptableBuildVersion(string) bool
		Init()
	}
)

func (this *CVersion)Init(){
	var bv1, bv2, bv3, bv4 int
	fmt.Sscanf(BUILD_NO, "%d,%d,%d,%d", &bv1, &bv2, &bv3, &bv4);
	this.m_sBuildVer = int64(bv1 * 1000*1000*1000 + bv2 * 1000*1000 + bv3 * 1000 + bv4);
}

func (this *CVersion)IsAcceptableBuildVersion(version string) bool{
	var v [4]int
	var sClient int64
	fmt.Sscanf(version, "%d,%d,%d,%d",&v[0], &v[1], &v[2], &v[3])
	sClient = int64(v[0] * 1000*1000*1000 + v[1] * 1000*1000 + v[2] * 1000 + v[3]);
	return  sClient >= this.m_sBuildVer
}

func CVERSION() *CVersion{
	if (g_Version.m_bInit == false){
		g_Version.Init();
	}
	return &g_Version
}

