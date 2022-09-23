package pkg

import "strings"

type Bilingual struct {
	CN string `json:"cn"`
	EN string `json:"en"`
}

type RawBody struct {
	BodyCN string `json:"bodyCn" db:"body_cn"`
	BodyEN string `json:"bodyEn" db:"body_en"`
}

func (b RawBody) HasEN() bool {
	return b.BodyEN != ""
}

func (b RawBody) HasCN() bool {
	return b.BodyCN != ""
}

func (b RawBody) IsEmpty() bool {
	return b.BodyCN == "" && b.BodyEN == ""
}

func (b RawBody) splitCN() []string {
	return strings.Split(b.BodyCN, "\r\n")
}

func (b RawBody) splitCNWithTranslator() ([]string, string) {
	paras := b.splitCN()

	l := len(paras)

	t := extractTranslator(paras[l-1])

	if t != "" {
		paras = paras[:l-1]
	}

	return paras, t
}

func (b RawBody) splitEN() []string {
	if b.BodyEN == "" {
		return make([]string, 0)
	}
	return strings.Split(b.BodyEN, "\r\n")
}

func extractTranslator(s string) string {
	if strings.Contains(s, "译者/") {
		return strings.TrimSuffix(strings.TrimPrefix(s, "<p>"), "</p>")
	}

	return ""
}
