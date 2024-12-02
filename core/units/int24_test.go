package units

import "testing"

func TestToUint32(t *testing.T) {
	tests := []struct {
		name     string
		input    Int24
		expected uint32
	}{
		{
			name:     "Zero value",
			input:    Int24{0, 0, 0},
			expected: 0,
		},
		{
			name:     "Max value",
			input:    Int24{0x7F, 0xFF, 0xFF},
			expected: 0x7FFFFF,
		},
		{
			name:     "Random value",
			input:    Int24{0x12, 0x34, 0x56},
			expected: 0x123456,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToUint32()
			if result != tt.expected {
				t.Errorf("ToUint32() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestIsNegative(t *testing.T) {
	tests := []struct {
		name     string
		input    Int24
		expected bool
	}{
		{
			name:     "Positive number",
			input:    Int24{0x00, 0x00, 0x01},
			expected: false,
		},
		{
			name:     "Negative number",
			input:    Int24{0x80, 0x00, 0x00},
			expected: true,
		},
		{
			name:     "Zero",
			input:    Int24{0x00, 0x00, 0x00},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsNegative()
			if result != tt.expected {
				t.Errorf("IsNegative() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToInt32(t *testing.T) {
	tests := []struct {
		name     string
		input    Int24
		expected int32
	}{
		{
			name:     "Positive number",
			input:    Int24{0x00, 0x00, 0x01},
			expected: 1,
		},
		{
			name:     "Negative number",
			input:    Int24{0x80, 0x00, 0x00},
			expected: -8388608, // MinInt24
		},
		{
			name:     "Zero",
			input:    Int24{0x00, 0x00, 0x00},
			expected: 0,
		},
		{
			name:     "Max positive",
			input:    Int24{0x7F, 0xFF, 0xFF},
			expected: 8388607, // MaxInt24
		},
		{
			name:     "Random positive",
			input:    Int24{0x12, 0x34, 0x56},
			expected: 1193046,
		},
		{
			name:     "Random negative",
			input:    Int24{0xFF, 0xFF, 0xFF},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToInt32()
			if result != tt.expected {
				t.Errorf("ToInt32() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a        Int24
		b        Int24
		expected Int24
	}{
		{
			name:     "Simple addition",
			a:        Int24{0x00, 0x00, 0x01},
			b:        Int24{0x00, 0x00, 0x02},
			expected: Int24{0x00, 0x00, 0x03},
		},
		{
			name:     "Addition with carry",
			a:        Int24{0x00, 0xFF, 0xFF},
			b:        Int24{0x00, 0x00, 0x01},
			expected: Int24{0x01, 0x00, 0x00},
		},
		{
			name:     "Addition with overflow",
			a:        Int24{0x7F, 0xFF, 0xFF},
			b:        Int24{0x00, 0x00, 0x01},
			expected: Int24{0x00, 0x00, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Add(tt.b)
			if result != tt.expected {
				t.Errorf("Add() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		name     string
		a        Int24
		b        Int24
		expected Int24
	}{
		{
			name:     "Simple subtraction",
			a:        Int24{0x00, 0x00, 0x03},
			b:        Int24{0x00, 0x00, 0x01},
			expected: Int24{0x00, 0x00, 0x02},
		},
		{
			name:     "Subtraction with borrow",
			a:        Int24{0x01, 0x00, 0x00},
			b:        Int24{0x00, 0x00, 0x01},
			expected: Int24{0x00, 0xFF, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Sub(tt.b)
			if result != tt.expected {
				t.Errorf("Sub() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		name     string
		a        Int24
		b        Int24
		expected Int24
	}{
		{
			name:     "Simple multiplication",
			a:        Int24{0x00, 0x00, 0x02},
			b:        Int24{0x00, 0x00, 0x03},
			expected: Int24{0x00, 0x00, 0x06},
		},
		{
			name:     "Multiplication with carry",
			a:        Int24{0x00, 0x01, 0x00},
			b:        Int24{0x00, 0x00, 0x02},
			expected: Int24{0x00, 0x02, 0x00},
		},
		{
			name:     "Large multiplication within bounds",
			a:        Int24{0x00, 0xFF, 0xFF},
			b:        Int24{0x00, 0x00, 0x02},
			expected: Int24{0x01, 0xFF, 0xFE},
		},
		{
			name:     "Multiplication with overflow",
			a:        Int24{0x7F, 0xFF, 0xFF},
			b:        Int24{0x00, 0x00, 0x02},
			expected: Int24{0xFF, 0xFF, 0xFE},
		},
		{
			name:     "Multiplication by zero",
			a:        Int24{0x12, 0x34, 0x56},
			b:        Int24{0x00, 0x00, 0x00},
			expected: Int24{0x00, 0x00, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Mul(tt.b)
			if result != tt.expected {
				t.Errorf("Mul() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBitwiseOperations(t *testing.T) {
	a := Int24{0x12, 0x34, 0x56}
	b := Int24{0x0F, 0x0F, 0x0F}

	t.Run("AND", func(t *testing.T) {
		expected := Int24{0x02, 0x04, 0x06}
		result := a.And(b)
		if result != expected {
			t.Errorf("And() = %v, want %v", result, expected)
		}
	})

	t.Run("OR", func(t *testing.T) {
		expected := Int24{0x1F, 0x3F, 0x5F}
		result := a.Or(b)
		if result != expected {
			t.Errorf("Or() = %v, want %v", result, expected)
		}
	})

	t.Run("XOR", func(t *testing.T) {
		expected := Int24{0x1D, 0x3B, 0x59}
		result := a.Xor(b)
		if result != expected {
			t.Errorf("Xor() = %v, want %v", result, expected)
		}
	})

	t.Run("NOT", func(t *testing.T) {
		expected := Int24{0xED, 0xCB, 0xA9}
		result := a.Not()
		if result != expected {
			t.Errorf("Not() = %v, want %v", result, expected)
		}
	})
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		a        Int24
		b        Int24
		expected int
	}{
		{
			name:     "Equal values",
			a:        Int24{0x12, 0x34, 0x56},
			b:        Int24{0x12, 0x34, 0x56},
			expected: 0,
		},
		{
			name:     "First greater",
			a:        Int24{0x12, 0x34, 0x57},
			b:        Int24{0x12, 0x34, 0x56},
			expected: 1,
		},
		{
			name:     "Second greater",
			a:        Int24{0x12, 0x34, 0x56},
			b:        Int24{0x12, 0x34, 0x57},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Compare(tt.b)
			if result != tt.expected {
				t.Errorf("Compare() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name            string
		input           Int24
		wantDecUnsigned string
		wantDecSigned   string
		wantHex         string
		wantBin         string
	}{
		{
			name:            "Positive number",
			input:           Int24{0x12, 0x34, 0x56},
			wantDecUnsigned: "1193046",
			wantDecSigned:   "1193046",
			wantHex:         "12 34 56",
			wantBin:         "00010010 00110100 01010110",
		},
		{
			name:            "Negative number",
			input:           Int24{0x80, 0x00, 0x01},
			wantDecUnsigned: "8388609",
			wantDecSigned:   "-8388607",
			wantHex:         "80 00 01",
			wantBin:         "10000000 00000000 00000001",
		},
		{
			name:            "Zero",
			input:           Int24{0x00, 0x00, 0x00},
			wantDecUnsigned: "0",
			wantDecSigned:   "0",
			wantHex:         "00 00 00",
			wantBin:         "00000000 00000000 00000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.StringDecUnsigned(); got != tt.wantDecUnsigned {
				t.Errorf("StringDecUnsigned() = %v, want %v", got, tt.wantDecUnsigned)
			}
			if got := tt.input.StringDecSigned(); got != tt.wantDecSigned {
				t.Errorf("StringDecSigned() = %v, want %v", got, tt.wantDecSigned)
			}
			if got := tt.input.StringHex(); got != tt.wantHex {
				t.Errorf("StringHex() = %v, want %v", got, tt.wantHex)
			}
			if got := tt.input.StringBin(); got != tt.wantBin {
				t.Errorf("StringBin() = %v, want %v", got, tt.wantBin)
			}
		})
	}
}
