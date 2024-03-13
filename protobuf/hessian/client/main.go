package main

import "github.com/cloudwego/kitex-benchmark/runner"

func main() {
	runner.Main("hessian", NewHessianClient)
}
