// go-qrcode
// Copyright 2014 Tom Harwood

package reedsolomon

import (
	"fmt"
	"log"
	"os"

	bitset "github.com/skip2/go-qrcode/bitset"
)

// gfPoly is a polynomial over GF(2^8).
type gfPoly struct {
	// The ith value is the coefficient of the ith degree of x.
	// term[0]*(x^0) + term[1]*(x^1) + term[2]*(x^2) ...
	term []gfElement
}

// newGFPolyFromData returns |data| as a polynomial over GF(2^8).
//
// Each data byte becomes the coefficient of an x term.
//
// For an n byte input the polynomial is:
// data[n-1]*(x^n-1) + data[n-2]*(x^n-2) ... + data[0]*(x^0).
func newGFPolyFromData(data *bitset.Bitset) gfPoly {
	numTotalBytes := data.Len() / 8
	if data.Len()%8 != 0 {
		numTotalBytes++
	}

	result := gfPoly{term: make([]gfElement, numTotalBytes)}

	i := numTotalBytes - 1
	for j := 0; j < data.Len(); j += 8 {
		result.term[i] = gfElement(data.ByteAt(j))
		i--
	}

	return result
}

// newGFPolyMonomial returns term*(x^degree).
func newGFPolyMonomial(term gfElement, degree int) gfPoly {
	if term == gfZero {
		return gfPoly{}
	}

	result := gfPoly{term: make([]gfElement, degree+1)}
	result.term[degree] = term

	return result
}

func (e gfPoly) data(numTerms int) []byte {
	result := make([]byte, numTerms)

	i := numTerms - len(e.term)
	for j := len(e.term) - 1; j >= 0; j-- {
		result[i] = byte(e.term[j])
		i++
	}

	return result
}

// numTerms returns the number of
func (e gfPoly) numTerms() int {
	return len(e.term)
}

// gfPolyMultiply returns a * b.
func gfPolyMultiply(a, b gfPoly) gfPoly {
	perrcheck(nil, "gfPolyMultiply - 1")
	numATerms := a.numTerms()

	perrcheck(nil, "gfPolyMultiply - 2")
	numBTerms := b.numTerms()

	perrcheck(nil, "gfPolyMultiply - 3")
	result := gfPoly{term: make([]gfElement, numATerms+numBTerms)}

	perrcheck(nil, "gfPolyMultiply - 4")
	for i := 0; i < numATerms; i++ {
		perrcheck(nil, "gfPolyMultiply - 5")
		index := i
		for j := 0; j < numBTerms; j++ {
			jIndex := j
			perrcheck(nil, "gfPolyMultiply - 6")
			abaz := fmt.Sprint(a.term)
			perrcheck(nil, abaz)

			foo := fmt.Sprint(a.term[index])
			perrcheck(nil, foo)

			perrcheck(nil, "gfPolyMultiply - VOR")
			baz := fmt.Sprint(b.term)
			perrcheck(nil, baz)

			bindexTerm := b.term[jIndex]
			perrcheck(nil, "gfPolyMultiply - NACH")

			bar := fmt.Sprint(bindexTerm)
			perrcheck(nil, bar)

			if a.term[index] != 0 && b.term[jIndex] != 0 {
				perrcheck(nil, "gfPolyMultiply - 7")
				newEl := make([]gfElement, index+jIndex+1)
				perrcheck(nil, "gfPolyMultiply - 7.5")
				monomial := gfPoly{term: newEl}
				perrcheck(nil, "gfPolyMultiply - 8")
				monomial.term[index+jIndex] = gfMultiply(a.term[index], b.term[jIndex])
				perrcheck(nil, "gfPolyMultiply - 9")
				result = gfPolyAdd(result, monomial)
			}
		}
	}
	perrcheck(nil, "gfPolyMultiply - 10")
	return result.normalised()
}

// gfPolyRemainder return the remainder of numerator / denominator.
func gfPolyRemainder(numerator, denominator gfPoly) gfPoly {
	perrcheck(nil, "poly - 1")
	if denominator.equals(gfPoly{}) {
		log.Panicln("Remainder by zero")
	}
	perrcheck(nil, "poly - 2")
	remainder := numerator
	perrcheck(nil, "poly - 3")
	for remainder.numTerms() >= denominator.numTerms() {
		perrcheck(nil, "poly - 4")
		degree := remainder.numTerms() - denominator.numTerms()
		perrcheck(nil, "poly - 5")
		coefficient := gfDivide(remainder.term[remainder.numTerms()-1],
			denominator.term[denominator.numTerms()-1])

	  perrcheck(nil, "poly - 6")
		divisor := gfPolyMultiply(denominator,
			newGFPolyMonomial(coefficient, degree))

		perrcheck(nil, "poly - 7")
		remainder = gfPolyAdd(remainder, divisor)
	}
	perrcheck(nil, "poly - 8")
	return remainder.normalised()
}

// gfPolyAdd returns a + b.
func gfPolyAdd(a, b gfPoly) gfPoly {
	numATerms := a.numTerms()
	numBTerms := b.numTerms()

	numTerms := numATerms
	if numBTerms > numTerms {
		numTerms = numBTerms
	}

	result := gfPoly{term: make([]gfElement, numTerms)}

	for i := 0; i < numTerms; i++ {
		switch {
		case numATerms > i && numBTerms > i:
			result.term[i] = gfAdd(a.term[i], b.term[i])
		case numATerms > i:
			result.term[i] = a.term[i]
		default:
			result.term[i] = b.term[i]
		}
	}

	return result.normalised()
}

func (e gfPoly) normalised() gfPoly {
	numTerms := e.numTerms()
	maxNonzeroTerm := numTerms - 1

	for i := numTerms - 1; i >= 0; i-- {
		if e.term[i] != 0 {
			break
		}

		maxNonzeroTerm = i - 1
	}

	if maxNonzeroTerm < 0 {
		return gfPoly{}
	} else if maxNonzeroTerm < numTerms-1 {
		e.term = e.term[0 : maxNonzeroTerm+1]
	}

	return e
}

func (e gfPoly) string(useIndexForm bool) string {
	var str string
	numTerms := e.numTerms()

	for i := numTerms - 1; i >= 0; i-- {
		if e.term[i] > 0 {
			if len(str) > 0 {
				str += " + "
			}

			if !useIndexForm {
				str += fmt.Sprintf("%dx^%d", e.term[i], i)
			} else {
				str += fmt.Sprintf("a^%dx^%d", gfLogTable[e.term[i]], i)
			}
		}
	}

	if len(str) == 0 {
		str = "0"
	}

	return str
}

// equals returns true if e == other.
func (e gfPoly) equals(other gfPoly) bool {
	var minecPoly *gfPoly
	var maxecPoly *gfPoly

	if e.numTerms() > other.numTerms() {
		minecPoly = &other
		maxecPoly = &e
	} else {
		minecPoly = &e
		maxecPoly = &other
	}

	numMinTerms := minecPoly.numTerms()
	numMaxTerms := maxecPoly.numTerms()

	for i := 0; i < numMinTerms; i++ {
		if e.term[i] != other.term[i] {
			return false
		}
	}

	for i := numMinTerms; i < numMaxTerms; i++ {
		if maxecPoly.term[i] != 0 {
			return false
		}
	}

	return true
}

func perrcheck(err error, str string) {
	f, e := os.OpenFile("gf_poly.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if e != nil {
		log.Fatalf("error opening file: %v", err)
		// fmt.Println(str, err)
		// os.Exit(1)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(str, err)

}
