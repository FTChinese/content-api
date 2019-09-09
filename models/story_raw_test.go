package models

import "testing"

func TestRawStory_BylineCN(t *testing.T) {
	s := RawStory{
		BylineDescCN:   "英国《金融时报》",
		BylineAuthorCN: "詹姆斯•波利提,阿利斯泰尔•格雷;理查德•亨德森",
		BylineStatusCN: "华盛顿,纽约报道",
	}

	bl := s.BylineCN()

	t.Logf("%+v", bl.Authors)
}
