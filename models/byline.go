package models

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
	Organization string   `json:"organization"`
	Authors      []Author `json:"authors"`
}

type Author struct {
	Names []string `json:"name"`
	Place string   `json:"place"`
}
