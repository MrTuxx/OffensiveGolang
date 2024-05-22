package evasion

func DecodeBrainfuck(code string, input string) string {
	const memSize = 30000
	memory := make([]byte, memSize)
	dataPointer := 0
	codePointer := 0
	inputPointer := 0
	output := ""

	for codePointer < len(code) {
		switch code[codePointer] {
		case '>':
			dataPointer++
		case '<':
			dataPointer--
		case '+':
			memory[dataPointer]++
		case '-':
			memory[dataPointer]--
		case '.':
			output += string(memory[dataPointer])
		case ',':
			if inputPointer < len(input) {
				memory[dataPointer] = input[inputPointer]
				inputPointer++
			}
		case '[':
			if memory[dataPointer] == 0 {
				loop := 1
				for loop > 0 {
					codePointer++
					if code[codePointer] == '[' {
						loop++
					} else if code[codePointer] == ']' {
						loop--
					}
				}
			}
		case ']':
			loop := 1
			for loop > 0 {
				codePointer--
				if code[codePointer] == '[' {
					loop--
				} else if code[codePointer] == ']' {
					loop++
				}
			}
			codePointer--
		}
		codePointer++
	}
	return output
}
func CodeBrainfuck(input string) string {
	const memSize = 30000
	memory := make([]byte, memSize)
	dataPointer := 0
	output := ""

	for _, char := range input {
		target := byte(char)
		for memory[dataPointer] != target {
			if memory[dataPointer] < target {
				memory[dataPointer]++
				output += "+"
			} else {
				memory[dataPointer]--
				output += "-"
			}
		}
		output += "."
	}

	return output
}
