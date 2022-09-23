package pkg

import "testing"

func TestByline_String(t *testing.T) {
	bl := Byline{
		Organization: "英国《金融时报》",
		Authors: []Authors{
			{
				Names: []string{"詹姆斯•波利提"},
				Place: "华盛顿",
			},
			{
				Names: []string{"阿利斯泰尔•格雷", "理查德•亨德森"},
				Place: "纽约报道",
			},
		},
	}

	t.Logf("%s", bl)
}
