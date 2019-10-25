package main

func recurse(n int) int {
	if n <= 0 {
		return 0
	}
	return n + recurse(n-1)
}

func main() {
	println(recurse(2000000))
}
