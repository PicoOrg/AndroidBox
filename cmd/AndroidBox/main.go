package main

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
