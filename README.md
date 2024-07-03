# KIT - Advanced Data Structures
As part of this course we implemented a bit vector implementation.

## Build instructions

1. To build this project you need >`go 1.22.4`.
2. Then run the `build.sh` script.
   
   This will run `go mod download` and `go build -o bitvector ./cmd/bitvector/main.go`
3. Use the `bitvector` binary for the competition.
4. Profit

This project uses two dependencies which are all unrelated to the algorithm:

- [github.com/urfave/cli/v2](https://github.com/urfave/cli/v2) for the command line interface of the test file generator.
- [github.com/stretchr/testify](https://github.com/stretchr/testify) for unit testing.

This can be validated through the `go.mod` file.


## Documentation

The code is documented using comments, in the `article` directory the LaTeX code for the paper is stored.
If you want to read this code, start with the [main.go](cmd/bitvector/main.go) and go from there.
The implementation of the actual data structure is implemented inside the [interleaved_vector.go](pkg/bit/interleaved_vector.go).
Also interesting are [vector.go](pkg/bit/vector.go) and [make_tables.go](pkg/bit/make_tables.go) which generates the `select` static lookup table.

### Benchmark

As part of the evaluation we created our own benchmark which works by creating test command files with increasing bit vector size.
How to reproduce:
1. Build the generator using the `build_bitvector.sh` script.
    (also build the `bitvector` binary using the `build.sh` script.)
2. Navigate into the `benchmark` directory and run the `generate.sh` script.
This will create directories containing `command.txt` files and `expected.txt` files,
with increasing vector size. Starting from $2^{8}$ to $2^{34}$.
3. Run the `run.sh` script to execute 11 runs of each vector size.
The result will be written to the `results.txt`.
