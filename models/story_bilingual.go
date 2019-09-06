package models

type Bilingual struct {
	CN string `json:"cn"`
	EN string `json:"en"`
}

type SplitBody struct {
	CN []string `json:"cn"`
	EN []string `json:"en"`
}
