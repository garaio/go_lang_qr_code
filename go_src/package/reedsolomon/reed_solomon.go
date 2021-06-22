// go-qrcode
// Copyright 2014 Tom Harwood

// Package reedsolomon provides error correction encoding for QR Code 2005.
//
// QR Code 2005 uses a Reed-Solomon error correcting code to detect and correct
// errors encountered during decoding.
//
// The generated RS codes are systematic, and consist of the input data with
// error correction bytes appended.
package reedsolomon

import (
	"log"
	"os"

	bitset "github.com/skip2/go-qrcode/bitset"
)

// Encode data for QR Code 2005 using the appropriate Reed-Solomon code.
//
// numECBytes is the number of error correction bytes to append, and is
// determined by the target QR Code's version and error correction level.
//
// ISO/IEC 18004 table 9 specifies the numECBytes required. e.g. a 1-L code has
// numECBytes=7.
func Encode(data *bitset.Bitset, numECBytes int) *bitset.Bitset {
	// Create a polynomial representing |data|.
	//
	// The bytes are interpreted as the sequence of coefficients of a polynomial.
	// The last byte's value becomes the x^0 coefficient, the second to last
	// becomes the x^1 coefficient and so on.
	errcheck(nil, "red - 1")
	ecpoly := newGFPolyFromData(data)
	errcheck(nil, "red - 2")
	ecpoly = gfPolyMultiply(ecpoly, newGFPolyMonomial(gfOne, numECBytes))
	errcheck(nil, "red - 3")
	// Pick the generator polynomial.
	generator := rsGeneratorPoly(numECBytes)
	errcheck(nil, "red - 4")
	// Generate the error correction bytes.
	remainder := gfPolyRemainder(ecpoly, generator)
	errcheck(nil, "red - 5")
	// Combine the data & error correcting bytes.
	// The mathematically correct answer is:
	//
	//	result := gfPolyAdd(ecpoly, remainder).
	//
	// The encoding used by QR Code 2005 is slightly different this result: To
	// preserve the original |data| bit sequence exactly, the data and remainder
	// are combined manually below. This ensures any most significant zero bits
	// are preserved (and not optimised away).
	errcheck(nil, "red - 6")
	result := bitset.Clone(data)
	errcheck(nil, "red - 7")
	result.AppendBytes(remainder.data(numECBytes))

	return result
}

// rsGeneratorPoly returns the Reed-Solomon generator polynomial with |degree|.
//
// The generator polynomial is calculated as:
// (x + a^0)(x + a^1)...(x + a^degree-1)
func rsGeneratorPoly(degree int) gfPoly {
	if degree < 2 {
		log.Panic("degree < 2")
	}

	generator := gfPoly{term: []gfElement{1}}

	for i := 0; i < degree; i++ {
		nextPoly := gfPoly{term: []gfElement{gfExpTable[i], 1}}
		generator = gfPolyMultiply(generator, nextPoly)
	}

	return generator
}

func errcheck(err error, str string) {
	f, e := os.OpenFile("reed_solomon.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if e != nil {
		log.Fatalf("error opening file: %v", err)
		// fmt.Println(str, err)
		// os.Exit(1)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(str, err)

}
