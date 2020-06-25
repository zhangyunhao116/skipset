# benchstat installation: go get -u golang.org/x/perf/cmd/benchstat

# 100,000 operations.
go test -bench=. -test.run=NOTEST -timeout=30m -benchmem -cpu=8 -count=100 -benchtime=100000x  > 100,000-100.txt
benchstat 100,000-100.txt

# 1,000,000 operations.
#go test -bench=. -test.run=NOTEST -timeout=30m -benchmem -cpu=8 -count=100 -benchtime=1000000x > 1,000,000-100.txt
#benchstat 1,000,000-100.txt