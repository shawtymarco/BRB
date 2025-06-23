package font

import (
	"github.com/samber/lo"
)

var nexusFonts = map[string]string{
	"A":  "\uF100",
	"B":  "\uF101",
	"C":  "\uF102",
	"D":  "\uF103",
	"E":  "\uF104",
	"F":  "\uF105",
	"G":  "\uF106",
	"H":  "\uF107",
	"I":  "\uF108",
	"J":  "\uF109",
	"K":  "\uF10A",
	"L":  "\uF10B",
	"M":  "\uF10C",
	"N":  "\uF10D",
	"O":  "\uF10E",
	"P":  "\uF10F",
	"Q":  "\uF110",
	"R":  "\uF111",
	"S":  "\uF112",
	"T":  "\uF113",
	"U":  "\uF114",
	"V":  "\uF115",
	"W":  "\uF116",
	"X":  "\uF117",
	"Y":  "\uF118",
	"Z":  "\uF119",
	"1":  "\uF11A",
	"2":  "\uF11B",
	"3":  "\uF11C",
	"4":  "\uF11D",
	"5":  "\uF11E",
	"6":  "\uF11F",
	"7":  "\uF120",
	"8":  "\uF121",
	"9":  "\uF122",
	"0":  "\uF123",
	"!":  "\uF124",
	"@":  "\uF125",
	"#":  "\uF126",
	"$":  "\uF127",
	"%":  "\uF128",
	"^":  "\uF129",
	"&":  "\uF12A",
	"*":  "\uF12B",
	"(":  "\uF12C",
	")":  "\uF12D",
	"-":  "\uF12E",
	"=":  "\uF12F",
	"_":  "\uF130",
	"+":  "\uF131",
	"[":  "\uF132",
	"]":  "\uF133",
	"<":  "\uF134",
	">":  "\uF135",
	"\\": "\uF136",
	"|":  "\uF137",
	";":  "\uF138",
	":":  "\uF139",
	"‘":  "\uF13A",
	"’":  "\uF13B",
	"“":  "\uF13C",
	",":  "\uF13D",
	".":  "\uF13E",
	"/":  "\uF13F",
}

func Transform(msg string) string {
	var res string
	for _, r := range []rune(msg) {
		c := string(r)
		res += lo.If(nexusFonts[c] != "", nexusFonts[c]).Else(c)
	}
	return res
}
