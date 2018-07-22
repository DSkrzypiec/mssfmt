package stringF

import (
	"fmt"
	"strings"
	"testing"
)

// Regular case with some different separators.
func TestSplitWithSep1(t *testing.T) {
	input1 := "Damian,Skrzypiec(Is)The\tGuy"
	output1 := SplitWithSep(input1)

	if len(output1) != 9 {
		t.Errorf("Expected output1 length = 9 | Got = %d \n", len(output1))
		fmt.Println("Expected: ", input1)
		fmt.Println("Got: ", strings.Join(output1, ""))
	}

	if strings.Join(output1, "") != input1 {
		t.Fail()
		fmt.Println("Expected: ", input1)
		fmt.Println("Got: ", strings.Join(output1, ""))
	}
}

// Starting from separator and ends on separator.
func TestSplitWithSep2(t *testing.T) {
	input2 := ",some words\nis SQL statements("
	output2 := SplitWithSep(input2)

	if len(output2) != 11 {
		t.Errorf("Expected output2 length = 11 | Got = %d \n", len(output2))
		fmt.Println("Expected: ", input2)
		fmt.Println("Got: ", strings.Join(output2, ""))
	}

	if strings.Join(output2, "") != input2 {
		t.Fail()
		fmt.Println("Expected: ", input2)
		fmt.Println("Got: ", strings.Join(output2, ""))
	}
}

// Case with only the separators.
func TestSplitWithSep3(t *testing.T) {
	input3 := ",,,,(((())))"
	output3 := SplitWithSep(input3)

	if len(output3) != 12 {
		t.Errorf("Expected output3 length = 12 | Got = %d \n", len(output3))
		fmt.Println("Expected: ", input3)
		fmt.Println("Got: ", strings.Join(output3, ""))
	}

	if strings.Join(output3, "") != input3 {
		t.Fail()
		fmt.Println("Expected: ", input3)
		fmt.Println("Got: ", strings.Join(output3, ""))
	}
}

// Single string without any sparator
func TestSplitWithSep4(t *testing.T) {
	input4 := "bigstringwithoutanyapusesasldkasldkalsdkalskdlaksd"
	output4 := SplitWithSep(input4)

	if len(output4) != 1 {
		t.Errorf("Expected output4 length = 1 | Got = %d \n", len(output4))
		fmt.Println("Expected: ", input4)
		fmt.Println("Got: ", strings.Join(output4, ""))
	}

	if strings.Join(output4, "") != input4 {
		t.Fail()
		fmt.Println("Expected: ", input4)
		fmt.Println("Got: ", strings.Join(output4, ""))
	}
}

// Benchmark for the common case
func BenchmarkSplitWithSepCommon(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SplitWithSep("\tmax(")
	}
}

// Benchmark for 10 separators
func BenchmarkSplitWithSep10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SplitWithSep("damian,splits,strings(in)some\nwierd\tmanner,testing testing)testing")
	}
}
