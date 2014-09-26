// abstraction for keyboard input

package keys

import (
	//"fmt"
)

type Key struct {
	Code int
	Flavr int
	Free, Shifted rune
	Name string
}

func New(code int) *Key {
	switch code {
	
	// handle control keys
	case Cesc, Cf1, Cf2, Cf3, Cf4, Cf5, Cf6, Cf7, Cf8, Cf9, Cf10, Cdel, Cbksp,
		 Ctab, Ccaps, Center, Clshft, Cslash, Crshft, Clctrl, Clalt, Cspace,
		 Cralt, Crctrl, Chome, Cup, Cpgup, Cleft, Cright, Cend, Cdown, Cpgdn:
		 
		return &Key{Code: code,
		 			Flavr: ControlKey,
		 			Name: Names[code]}
	
	// handle rune keys
	case Rbktik, R1, R2, R3, R4, R5, R6, R7, R8, R9, R0, Rhyph, Req,
		 Rq, Rw, Re, Rr, Rt, Ry, Ru, Ri, Ro, Rp, Rlbrak, Rrbrak, Rbslsh,
		 Ra, Rs, Rd, Rf, Rg, Rh, Rj, Rk, Rl, Rsmcln, Rquot,
		 Rz, Rx, Rc, Rv, Rb, Rn, Rm, Rcomma, Rdot:
		
		return &Key{Code: code,
					Flavr: RuneKey,
					Free: Runes[code][0],
					Shifted: Runes[code][1],
					Name: Names[code]}
	
	// return unknown key with just code for unrecognized key codes
	default:
		return &Key{Code: code,
					Flavr: UnknownKey}
	}
}

const (

	// key flavrs
	RuneKey 	int	= iota
	ControlKey	int	= iota
	UnknownKey	int	= iota
	
	// keys enumerated by x keycodes, C is control key, R is rune key
	Cesc	int	= 9	// row 0
	Cf1		int	= 67
	Cf2		int	= 68
	Cf3		int	= 69
	Cf4		int	= 70
	Cf5		int	= 71
	Cf6		int	= 72
	Cf7		int	= 73
	Cf8		int	= 74
	Cf9		int	= 75
	Cf10	int	= 76
	Cdel	int	= 119
		
	Rbktik	int	= 49	// row 1
	R1		int	= 10
	R2		int	= 11
	R3		int	= 12
	R4		int	= 13
	R5		int	= 14
	R6		int	= 15
	R7		int	= 16
	R8		int	= 17
	R9		int	= 18
	R0		int	= 19
	Rhyph	int	= 20
	Req		int	= 21
	Cbksp	int	= 22
	
	Ctab	int	= 23	// row 2
	Rq		int	= 24
	Rw		int	= 25
	Re		int	= 26
	Rr		int	= 27
	Rt		int	= 28
	Ry		int	= 29
	Ru		int	= 30
	Ri		int	= 31
	Ro		int	= 32
	Rp		int	= 33
	Rlbrak	int	= 34
	Rrbrak	int	= 35
	Rbslsh	int	= 51
	
	Ccaps	int	= 66	// row 3
	Ra		int	= 38
	Rs		int	= 39
	Rd		int	= 40
	Rf		int	= 41
	Rg		int	= 42
	Rh		int	= 43
	Rj		int	= 44
	Rk		int	= 45
	Rl		int	= 46
	Rsmcln	int	= 47
	Rquot	int	= 48
	Center	int	= 36
	
	Clshft	int	= 50	// row 4
	Rz		int	= 52
	Rx		int	= 53
	Rc		int	= 54
	Rv		int	= 55
	Rb		int	= 56
	Rn		int	= 57
	Rm		int	= 58
	Rcomma	int	= 59
	Rdot	int	= 60
	Cslash	int	= 61
	Crshft	int	= 62
	
	Clctrl	int	= 37	// row 5
	Clalt	int	= 64
	Cspace	int	= 65
	Cralt	int	= 108
	Crctrl	int	= 105
	
	Chome	int	= 110	// arrows
	Cup		int	= 111
	Cpgup	int	= 112
	Cleft	int	= 113
	Cright	int	= 114
	Cend	int	= 115
	Cdown	int	= 116
	Cpgdn	int	= 117
)

var (

	// map from x keycode to free and shifted runes
	Runes = map[int][2]rune{
		10: {'1', '!'},
		11: {'2', '@'},
		12: {'3', '#'},
		13: {'4', '$'},
		14: {'5', '%'},
		15: {'6', '^'},
		16: {'7', '&'},
		17: {'8', '*'},
		18: {'9', '#'},
		19: {'0', ')'},
		20: {'-', '_'},
		21: {'=', '+'},
		24: {'q', 'Q'},
		25: {'w', 'W'},
		26: {'e', 'E'},
		27: {'r', 'R'},
		28: {'t', 'T'},
		29: {'y', 'Y'},
		30: {'u', 'U'},
		31: {'i', 'I'},
		32: {'o', 'O'},
		33: {'p', 'P'},
		34: {'[', '{'},
		35: {']', '}'},
		38: {'a', 'A'},
		39: {'s', 'S'},
		40: {'d', 'D'},
		41: {'f', 'F'},
		42: {'g', 'G'},
		43: {'h', 'H'},
		44: {'j', 'J'},
		45: {'k', 'K'},
		46: {'l', 'L'},
		47: {';', ':'},
		48: {'\'', '"'},
		49: {'`', '~'},
		51: {'\\', '|'},
		52: {'z', 'Z'},
		53: {'x', 'X'},
		54: {'c', 'C'},
		55: {'v', 'V'},
		56: {'b', 'B'},
		57: {'n', 'N'},
		58: {'m', 'M'},
		59: {',', '<'},
		60: {'.', '>'},
		61: {'/', '?'},
	}
	
	Names = map[int]string{
		9:		"esc",			// row 0
		67:		"F1",
		68:		"F2",
		69: 	"F3",
		70: 	"F4",
		71: 	"F5",
		72: 	"F6",
		73: 	"F7",
		74: 	"F8",
		75: 	"F9",
		76:		"F10",
		119:	"delete",
		
		49:		"`",			// row 1
		10:		"1",
		11:		"2",
		12:		"3",
		13:		"4",
		14:		"5",
		15:		"6",
		16:		"7",
		17:		"8",
		18:		"9",
		19:		"0",
		20:		"-",
		21:		"=",
		22:		"backspace",
	
		23:		"tab",			// row 2
		24:		"q",
		25:		"w",
		26:		"e",
		27:		"r",
		28:		"t",
		29:		"y",
		30:		"u",
		31:		"i",
		32:		"o",
		33:		"p",
		34:		"[",
		35:		"]",
		51:		"\\",
	
		66:		"caps lock",	// row 3
		38:		"a",
		39:		"s",
		40:		"d",
		41:		"f",
		42:		"g",
		43:		"h",
		44:		"j",
		45:		"k",
		46:		"l",
		47:		";",
		48:		"'",
		36:		"enter",
	
		50:		"left shift",	// row 4
		52:		"z",
		53:		"x",
		54:		"c",
		55:		"v",
		56:		"b",
		57:		"n",
		58:		"m",
		59:		",",
		60:		".",
		61:		"/",
		62:		"right shift",
	
		37:		"left control",	// row 5
		64:		"left alt",
		65:		"space",
		108:	"right alt",
		105:	"right control",
	
		110:	"home",			// arrows
		111:	"up arrow",
		112:	"page up",
		113:	"left arrow",
		114:	"right arrow",
		115:	"end",
		116:	"down arrow",
		117:	"page down",
	}
)
