package textprocessor

import (
	"fmt"
	"io"
	"strings"
)

func Run(src io.Reader, dst io.Writer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	var sb strings.Builder
	b := make([]byte, 100)
	n, err := src.Read(b)
	for err != nil {
		if _, err1 := sb.Write(b[:n]); err1 != nil {
			return err1
		}
		n, err = src.Read(b)
	}
	if err != io.EOF {
		return err
	}

	fields := strings.Fields(sb.String())
	var finalFields []string
	for _, field := range fields {
		switch field {
		case "(hex)":
			if len(finalFields) == 0 {
				return fmt.Errorf("Unexpected token (no word to apply to)")
			}
		case "(bin)":
			if len(finalFields) == 0 {
				return fmt.Errorf("Unexpected token (no word to apply to)")
			}
		case "(up)":
			if len(finalFields) == 0 {
				return fmt.Errorf("Unexpected token (no word to apply to)")
			}
		case "(low)":
			if len(finalFields) == 0 {
				return fmt.Errorf("Unexpected token (no word to apply to)")
			}
		case "(cap)":
			if len(finalFields) == 0 {
				return fmt.Errorf("Unexpected token (no word to apply to)")
			}
		default:
		}
	}

	return nil
}
