package base

import(
	"time"
)

var(
	s_Quotiet int
	s_Remainder int
	s_Seed int
)

type (
	MRadomLCG struct{
		mSeed int
	}

	IMRandomLCG interface {
		Init()
		setSeed(int)
		getSeed() int
		rand() int
		RandI(int, int) int
	}
)

func (this *MRadomLCG) Init(){
	s_Quotiet = INT_MAX / 16807
	s_Remainder = INT_MAX % 16807
	s_Seed = 1376312589

	this.setSeed(generateSeed())
}

func (this *MRadomLCG) setSeed(s int)  {
	this.mSeed = s
}

func (this *MRadomLCG) getSeed() int {
	return  this.mSeed
}

func (this *MRadomLCG) rand() int{
	if this.mSeed <= s_Quotiet {
		this.mSeed = (this.mSeed * 16807) % INT_MAX
	}else{
		var high_part = this.mSeed / s_Quotiet
		var low_part = this.mSeed % s_Quotiet

		var test = (16807 * low_part) - (s_Remainder * high_part)

		if test > 0 {
			this.mSeed = test
		}else{
			this.mSeed = test + INT_MAX
		}
	}
	return  this.mSeed
}

func (this *MRadomLCG) RandI(i int, n int)int {
	if i > n {
		Assert(false, "MRandomGenerator::randi: inverted range")
	}

	return  int(i + this.rand() % (n - i  + 1))
}

func generateSeed() int{
	s_Seed = int(time.Now().Unix())
	s_Seed = (s_Seed * 0x015a4e35) + 1
	s_Seed = (s_Seed>>16)&0x7fff;
	return  s_Seed
}

var (
	pRadomMgr *MRadomLCG
)

func RANDOMMGR() *MRadomLCG{
	if pRadomMgr == nil {
		pRadomMgr = new(MRadomLCG)
		pRadomMgr.Init()
	}
	return  pRadomMgr
}