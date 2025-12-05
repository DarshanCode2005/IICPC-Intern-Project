package problems

import "xcodeengine/model"

var problemSet = []model.Problem{
	{
		ID:          "sum-two-numbers",
		Title:       "Sum Two Numbers",
		Description: "### Task\nRead two integers and output their sum.\n\n### Notes\n- Input fits in 32-bit signed integer.\n- Output should include a newline.",
		InputFormat: "Two integers A and B separated by space.",
		Constraints: "`0 ≤ A, B ≤ 10^9`",
		TestCases: []model.TestCase{
			{Name: "Sample #1", Input: "2 3\n", ExpectedOutput: "5\n"},
			{Name: "Sample #2", Input: "100 250\n", ExpectedOutput: "350\n"},
		},
	},
	{
		ID:          "fizzbuzz",
		Title:       "FizzBuzz",
		Description: "### Task\nPrint numbers from 1 to N.\n- Multiples of 3 => `Fizz`\n- Multiples of 5 => `Buzz`\n- Multiples of 15 => `FizzBuzz`",
		InputFormat: "Single integer N.",
		Constraints: "`1 ≤ N ≤ 10^3`",
		TestCases: []model.TestCase{
			{
				Name:           "Sample #1",
				Input:          "5\n",
				ExpectedOutput: "1\n2\nFizz\n4\nBuzz\n",
			},
			{
				Name:           "Sample #2",
				Input:          "15\n",
				ExpectedOutput: "1\n2\nFizz\n4\nBuzz\nFizz\n7\n8\nFizz\nBuzz\n11\nFizz\n13\n14\nFizzBuzz\n",
			},
		},
	},
	{
		ID:          "balanced-brackets",
		Title:       "Balanced Brackets",
		Description: "### Task\nGiven a string of brackets, determine if the sequence is balanced.",
		InputFormat: "A single string containing characters `()[]{}` only.",
		Constraints: "`1 ≤ length ≤ 10^5`",
		TestCases: []model.TestCase{
			{Name: "Sample #1", Input: "{}[]()\n", ExpectedOutput: "YES\n"},
			{Name: "Sample #2", Input: "{[}]\n", ExpectedOutput: "NO\n"},
		},
	},
	{
		ID:          "two-sum",
		Title:       "Two Sum",
		Description: "### Task\nGiven an array and a target, determine if any pair sums to the target.",
		InputFormat: "First line: N and target. Second line: N integers.",
		Constraints: "`2 ≤ N ≤ 10^5` (values fit in 32-bit signed int)",
		TestCases: []model.TestCase{
			{Name: "Sample #1", Input: "4 9\n2 7 11 15\n", ExpectedOutput: "YES\n"},
			{Name: "Sample #2", Input: "3 10\n1 2 3\n", ExpectedOutput: "NO\n"},
		},
	},
	{
		ID:          "matrix-trace",
		Title:       "Matrix Trace",
		Description: "### Task\nCompute the trace of an `N x N` matrix (sum of diagonal elements).",
		InputFormat: "First line: N. Next N lines: N integers each.",
		Constraints: "`1 ≤ N ≤ 200`",
		TestCases: []model.TestCase{
			{
				Name:           "Sample #1",
				Input:          "3\n1 2 3\n4 5 6\n7 8 9\n",
				ExpectedOutput: "15\n",
			},
			{
				Name:           "Sample #2",
				Input:          "2\n10 1\n1 10\n",
				ExpectedOutput: "20\n",
			},
		},
	},
}

func ListProblems() []model.Problem {
	return problemSet
}

func GetProblem(id string) (*model.Problem, bool) {
	for _, p := range problemSet {
		if p.ID == id {
			cp := p
			return &cp, true
		}
	}
	return nil, false
}
