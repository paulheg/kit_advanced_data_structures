package bitvector

import (
	"bufio"
	"fmt"
	"io"
	"math/bits"
	"math/rand"
	"strconv"
	"strings"

	"github.com/paulheg/kit_advanced_data_structures/pkg/bit"
)

func GenerateRandomTestCase(vectorSlices64, commands uint64, commandOut, expectedOut io.Writer) error {

	commandBuffer := bufio.NewWriterSize(commandOut, 1024*1024)
	defer commandBuffer.Flush()
	expectedBuffer := bufio.NewWriterSize(expectedOut, 1024*1024)
	defer expectedBuffer.Flush()

	commandBuffer.Write([]byte(fmt.Sprintf("%d\n", commands)))

	var ones, zeros uint64

	generatorFormatString := "%0" + strconv.FormatUint(bit.SubvectorBits, 10) + "b"

	var vector bit.Vector = make([]bit.Subvector, vectorSlices64)
	for i := 0; i < int(vectorSlices64); i++ {
		vector[i] = bit.Subvector(rand.Uint64())
		ones += uint64(bits.OnesCount64(uint64(vector[i])))
		binary := fmt.Sprintf(generatorFormatString, vector[i])

		bBinary := []byte(binary)

		for i, j := 0, len(bBinary)-1; i < j; i, j = i+1, j-1 {
			bBinary[i], bBinary[j] = bBinary[j], bBinary[i]
		}

		n, err := commandBuffer.Write(bBinary)
		if n != int(bit.SubvectorBits) || err != nil {
			return fmt.Errorf("error writing, n=%d: %w", n, err)
		}
	}

	zeros = vector.Bits() - ones
	commandBuffer.Write([]byte{'\n'})

	fullVector := bit.NewInterleavedVector(vector)

	for i := 0; i < int(commands); i++ {
		fullCommand, expectedResult, err := randomCommandAndResult(fullVector, ones, zeros)
		if err != nil {
			return err
		}

		commandBuffer.Write([]byte(fullCommand + "\n"))
		expectedBuffer.Write([]byte(expectedResult + "\n"))
	}

	return nil
}

func randomCommandAndResult(vec bit.RankSelectVector, ones, zeros uint64) (fullCommand string, result string, err error) {

	command, fullCommand := randomCommand(ones, zeros)

	executor, err := CommandExecutors[command](strings.Split(fullCommand, " ")[1:])
	if err != nil {
		return
	}

	resultValue, err := executor(vec)
	if err != nil {
		return
	}
	result = fmt.Sprintf("%d", resultValue)
	return
}

func randomCommand(ones, zeros uint64) (Command, string) {

	generators := map[Command]func() string{
		Access: func() string {
			return fmt.Sprintf("%s %d", Access, rand.Int63n(int64(ones+zeros)))
		},
		Rank: func() string {
			return fmt.Sprintf("%s %d %d", Rank, rand.Intn(2), rand.Int63n(int64(ones+zeros)))
		},
		Select: func() string {
			alpha := rand.Intn(2)

			position := 0
			if alpha == 0 {
				position = int(rand.Int63n(int64(zeros))) + 1
			} else {
				position = int(rand.Int63n(int64(ones))) + 1
			}

			return fmt.Sprintf("%s %d %d", Select, alpha, position)
		},
	}

	commandNum := rand.Intn(len(Commands))
	command := Commands[commandNum]

	return command, generators[command]()
}
