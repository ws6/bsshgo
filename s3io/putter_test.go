package s3io

//import (
//	"fmt"
//	"testing"
//)
//
//type GrowPartSizeCase struct {
//	StartPartSize int64
//	EndPartSize   int64
//
//	TimesGrown int
//	TotalSize  int64
//}
//
//func (tc GrowPartSizeCase) Run(t *testing.T) {
//	var ptSize = tc.StartPartSize
//	var ptUploaded int64
//	var nGrow int
//	for i := 0; ptUploaded < tc.TotalSize; i++ {
//		ptUploaded += ptSize
//		if i%2000 == 0 {
//			if growPartSize(i, ptSize, tc.TotalSize) {
//				ptSize = ptSize * 2
//				nGrow++
//			}
//		}
//	}
//
//	if ptSize != tc.EndPartSize {
//		t.Fatalf("Expected final part size of %d but got %d; %d delta", tc.EndPartSize, ptSize, ptSize-tc.EndPartSize)
//	}
//	if nGrow != tc.TimesGrown {
//		t.Fatalf("Expected part size to grow %d times; got %d", tc.TimesGrown, nGrow)
//	}
//}
//
//func TestGrowPartSize(t *testing.T) {
//	cases := []GrowPartSizeCase{
//		{
//			StartPartSize: 16 * mb,
//			EndPartSize:   32 * mb,
//			TotalSize:     180 * gb,
//			TimesGrown:    1,
//		},
//	}
//
//	for _, tc := range cases {
//		t.Run(fmt.Sprintf("%d-%d@%d", tc.StartPartSize, tc.EndPartSize, tc.TotalSize), tc.Run)
//	}
//}
