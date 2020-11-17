package main

func main() {
L1:
	for x := 0; x < 3; x++ {
	L2:
		for y := 0; y < 5; y++ {
			if y > 2 {
				break L2
			}
			if x > 1 {
				break L1
			}

			print(x, ":", y, " ")
		}

		println()
	}
}
