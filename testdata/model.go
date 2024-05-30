package model

type Struct1 struct {
	A, Alt string `csv:"a"`
	B      string `csv:"b"`
	C      string `csv:"c"`
	D      string `csv:"d"`
	E      string `csv:"e"`
	F      string
}

type Struct2 struct {
	F, Alt string `csv:"f"`
	G      string `csv:"g"`
	H      string `csv:"h"`
	I      string `csv:"i"`
	J      string `csv:"j"`
	K      string
}

type Struct3 struct {
	K, Alt string
	L      string
	M      string
	N      string
	O      string
	P      string
}
