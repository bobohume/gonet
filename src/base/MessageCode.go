package base

import (
	"hash/crc32"
	"strings"
)

const(
	CRYPT_BUFFER_SIZE = 0x500
)

var (
	bCryptBufferCreated bool
	CryptBuffer []uint32
)

func PrepareCryptBuffer(){
	dwSeed := uint32(0x00100001)
	if (bCryptBufferCreated == false){
		CryptBuffer = make([]uint32, CRYPT_BUFFER_SIZE)
		for index1 := 0; index1 < 0x100; index1++{
			i := 0
			for index2 := index1; i < 5; index2 += 0x100{
				var temp1, temp2 uint32
				dwSeed = (dwSeed * 125 + 3) % 0x2AAAAB
				temp1  = (dwSeed & 0xFFFF) << 0x10
				dwSeed = (dwSeed * 125 + 3) % 0x2AAAAB
				temp2  = (dwSeed & 0xFFFF)
				CryptBuffer[index2] = (temp1 | temp2)
				i++
			}
		}

		bCryptBufferCreated = true
	}
}

func DecryptName1(strName string) uint32 {
	buf := []byte(strings.ToLower(strName))
	dwSeed1	:= uint32(0x7FED7FED)
	dwSeed2	:= uint32(0xEEEEEEEE)
	var ch uint32

	for _, v := range buf{
		ch = uint32(v)
		dwSeed1 = CryptBuffer[0x100 + ch] ^ (dwSeed1 + dwSeed2)
		dwSeed2 = ch + dwSeed1 + dwSeed2 + (dwSeed2 << 5) + 3
	}
	return dwSeed1;
}

func DecryptName2(strName string) uint32 {
	buf := []byte(strings.ToLower(strName))
	dwSeed1	:= uint32(0x7FED7FED)
	dwSeed2	:= uint32(0xEEEEEEEE)
	var ch uint32

	for _, v := range buf{
		ch = uint32(v)
		dwSeed1 = CryptBuffer[0x200 + ch] ^ (dwSeed1 + dwSeed2)
		dwSeed2 = ch + dwSeed1 + dwSeed2 + (dwSeed2 << 5) + 3
	}

	return dwSeed1;
}

func GetMessageCode(strName string, m1 *uint32, m2 *uint32){
	PrepareCryptBuffer()
	*m1 = DecryptName1(strName)
	*m2 = DecryptName2(strName)
}

func GetMessageCode2(strName string) uint32 {
	PrepareCryptBuffer();
	return DecryptName1(strName);
}

var(

)

func GetMessageCode1(strName string) uint32 {
	return crc32.ChecksumIEEE([]byte(strName))
	//return GetMessageCode2(strName)
	/*h := fnv.New32()
	h.Write([]byte(strName))
	sum := h.Sum(nil)
	return uint32(BytesToInt(sum))*/
}