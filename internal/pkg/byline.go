package pkg

import (
	"fmt"
	"strings"
)

type Authors struct {
	Names []string `json:"names"`
	Place string   `json:"place"`
}

// Byline aggregates byline data
// {
//  "organization": "英国《金融时报》",
//  "authors": [
//      {
//          names: ["詹姆斯•波利提"],
//          place: "华盛顿"
//      },
//      {
//          names: ["阿利斯泰尔•格雷", "理查德•亨德森"],
//          place: "纽约报道"
//      }
//  ]
// }
type Byline struct {
	Organization string    `json:"organization"`
	Authors      []Authors `json:"authors"`
}

func (b Byline) String() string {
	var authors []string
	for _, v := range b.Authors {
		authors = append(authors, fmt.Sprintf("%s %s", strings.Join(v.Names, ","), v.Place))
	}

	return fmt.Sprintf("%s %s", b.Organization, strings.Join(authors, ", "))
}
