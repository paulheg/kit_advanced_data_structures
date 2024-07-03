package bitvector

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/paulheg/kit_advanced_data_structures/pkg/bit"
)

type Command string

const (
	Access Command = "access"
	Rank   Command = "rank"
	Select Command = "select"
)

var Commands []Command = []Command{
	Access, Rank, Select,
}

type CommandFuncGenerator func(args []string) (CommandFunc, error)
type CommandFunc func(vec bit.RankSelectVector) (uint64, error)

var CommandExecutors = map[Command]CommandFuncGenerator{
	Access: func(args []string) (CommandFunc, error) {
		if len(args) != 1 {
			return nil, errors.New("access only accepts one argument")
		}

		position, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("position argument not a valid number: %w", err)
		}

		return func(vec bit.RankSelectVector) (uint64, error) {
			b := vec.Access(position)
			v := uint64(0)
			if b {
				v = 1
			}
			return v, nil
		}, nil
	},
	Rank: func(args []string) (CommandFunc, error) {
		if len(args) != 2 {
			return nil, errors.New("rank only accepts two arguments")
		}

		alpha, err := strconv.ParseBool(args[0])
		if err != nil {
			return nil, fmt.Errorf("alpha argument not valid: %w", err)
		}

		position, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("position argument not valid: %w", err)
		}

		return func(vec bit.RankSelectVector) (uint64, error) {
			return vec.Rank(alpha, position), nil
		}, nil

	},
	Select: func(args []string) (CommandFunc, error) {
		if len(args) != 2 {
			return nil, errors.New("select only accepts two arguments")
		}

		alpha, err := strconv.ParseBool(args[0])
		if err != nil {
			return nil, fmt.Errorf("alpha argument not valid: %w", err)
		}

		position, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("position argument not valid: %w", err)
		}

		return func(vec bit.RankSelectVector) (uint64, error) {
			return vec.Select(alpha, position), nil
		}, nil
	},
}
