# benchstat installation: go get -u golang.org/x/perf/cmd/benchstat

# 100,000 operations.
go test -bench=. -run=NOTEST -benchmem -cpu=8 -count=20 -benchtime=100000x  > 10t.txt

# 1,000,000 operations.
go test -bench=. -run=NOTEST -benchmem -cpu=8 -count=20 -benchtime=1000000x > 100t.txt

#benchstat 1,000,000-100.txt
#benchstat 100,000-100.txt