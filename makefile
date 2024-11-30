test:
	go test ./... -race

bench:
	go test -count 5 -benchmem -bench Benchmark calendar_benchmark_test.go

bench_cpu:
	go test -count 3 -cpu=1,4,8 -benchmem -bench Benchmark calendar_benchmark_test.go

bench_before:
	go test -count 5 -benchmem -bench Benchmark calendar_benchmark_test.go 2>&1 | tee bench_before.log
bench_after:
	go test -count 5 -benchmem -bench Benchmark calendar_benchmark_test.go 2>&1 | tee bench_after.log
# prepare: go install golang.org/x/perf/cmd/benchstat@latest
bench_diff:
	benchstat bench_before.log bench_before.log

credit:
	gocredits . > CREDITS