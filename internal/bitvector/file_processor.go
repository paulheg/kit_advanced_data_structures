package bitvector

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/paulheg/kit_advanced_data_structures/pkg/bit"
)

func ProcessFile(input io.Reader, output io.Writer, statOut io.Writer, verbose bool) error {

	var err error

	reader := bufio.NewReader(input)

	noOfCommands := 0
	line, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("number of commands could not be read: %w", err)
	} else {
		noOfCommands, err = strconv.Atoi(line[:len(line)-1])
		if err != nil {
			return errors.New("number of commands not a valid number")
		}
	}

	// init baseVector
	var baseVector bit.Vector
	var vecLen uint64
	for {
		b, err := reader.ReadSlice('\n')
		if err != nil && err != bufio.ErrBufferFull {
			return fmt.Errorf("could not read bitvector: %w", err)
		}

		for _, v := range b {
			if vecLen%bit.SubvectorBits == 0 {
				baseVector = append(baseVector, 0)
			}

			if v == '1' {
				baseVector.Set(vecLen)
			}
			vecLen++
		}
		if err == nil {
			break
		}
	}

	// Move base vector into new structure, we could also skip this step and directly read into the interleaved structure.
	// Since this copying is not needed normally it does not contribute to the recorded runtime
	// The precomputation of our data structure will be done later and will contribute to the runtime
	vec := bit.NewInterleavedVectorNoPrecompute(baseVector)

	// To see the implementation of the interleaved vector go to pkg/bit/interleaved_vector.go

	commandFuncs := make([]CommandFunc, 0)

	// scan commands
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("could not read command: %w", err)
		}
		line = line[:len(line)-1]
		tokens := strings.Split(line, " ")

		command := Command(tokens[0])
		if c, ok := CommandExecutors[command]; ok {
			args := tokens[1:]
			commandFunc, err := c(args)
			if err != nil {
				return fmt.Errorf("could not parse %s with args %v: %s", command, args, err)
			}

			commandFuncs = append(commandFuncs, commandFunc)
		} else {
			return fmt.Errorf("command %s not found", command)
		}
	}

	results := make([]uint64, noOfCommands)

	begin := time.Now()
	// run pre computation which does contribute to the runtime (creating the prev. sums)
	vec.Precompute()

	endPrecompute := time.Now()
	precomputionTime := endPrecompute.Sub(begin)

	// run commands
	for i, commandFunc := range commandFuncs {

		result, err := commandFunc(vec)
		if err != nil {
			return fmt.Errorf("could not run command: %w", err)
		}

		results[i] = result
	}

	// stop timer
	end := time.Now()
	runtime := end.Sub(begin)
	commandTime := end.Sub(endPrecompute)

	for _, v := range results {
		output.Write([]byte(fmt.Sprintf("%d\n", v)))
	}

	overheadFrac := float64(vec.Overhead()) / float64(vec.Size())

	precomputionFac := float64(precomputionTime) / float64(runtime)

	statOut.Write([]byte(fmt.Sprintf("RESULT name=paul_hegenberg time=%d space=%d", runtime.Milliseconds(), vec.Size())))
	if verbose {
		statOut.Write([]byte(fmt.Sprintf(" overhead=%f precompTime=%d precompFac=%f commandTime=%d", overheadFrac, precomputionTime.Milliseconds(), precomputionFac, commandTime.Milliseconds())))
	}

	return nil
}
