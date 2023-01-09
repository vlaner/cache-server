package commands

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Command string

const (
	COMMAND_GET    Command = "GET"
	COMMAND_SET    Command = "SET"
	COMMAND_DEL    Command = "DEL"
	COMMAND_EXPIRE Command = "EXPIRE"
)

var (
	ErrInvalidCmd        = errors.New("invalid command")
	ErrInvalidParameters = errors.New("invalid parameters")
	ErrInvalidExpireTime = errors.New("invalid parameters")
)

type Payload struct {
	Cmd    Command
	Key    string
	Value  string
	Expire time.Duration
}

func Parse(input string) (*Payload, error) {
	input = strings.TrimRight(input, "\r\n")
	cmd := strings.Split(input, " ")

	if len(cmd) < 2 {
		return nil, ErrInvalidCmd
	}
	command := cmd[0]
	key := cmd[1]

	switch Command(strings.ToUpper(command)) {
	case COMMAND_GET:
		return &Payload{
			Cmd: COMMAND_GET,
			Key: key,
		}, nil

	case COMMAND_SET:
		if len(cmd) < 3 {
			return nil, ErrInvalidParameters
		}
		value := strings.Join(cmd[2:], " ")
		return &Payload{
			Cmd:   COMMAND_SET,
			Key:   key,
			Value: value,
		}, nil

	case COMMAND_DEL:
		return &Payload{
			Cmd: COMMAND_DEL,
			Key: key,
		}, nil
	case COMMAND_EXPIRE:
		if len(cmd) < 3 {
			return nil, ErrInvalidParameters
		}
		value := cmd[2]
		num, err := strconv.Atoi(value)
		if err != nil {
			return nil, ErrInvalidExpireTime
		}

		return &Payload{
			Cmd:    COMMAND_EXPIRE,
			Key:    key,
			Expire: time.Second * time.Duration(num),
		}, nil

	default:
		return nil, ErrInvalidCmd
	}
}
