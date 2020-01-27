package parser

import (
	"fmt"
	"io"
)

func readData(port io.Reader, data []byte) ([]byte, error) {
	buff := make([]byte, 1024*256)
	n, err := port.Read(buff)
	if err != nil {
		return data, err
	}
	data = append(data, buff[0:n]...)
	return data, nil
}

// GetData gets data from the given reader.
func GetData(port io.Reader) ([]byte, error) {
	var data []byte // max: 256kb
	var dataSize int

	buff := []byte{}
	idx := 0
	var err error

	step := 0
	for step < 3 {
		buff, err = readData(port, buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read: %v", err)
		}

		// step 0: find the header: FF FF AA 55
		if step == 0 {
			if len(buff) < 4 {
				continue
			}

			for idx = 0; idx < len(buff)-4; idx++ {
				if buff[idx] == 0xFF {
					idx++
					if buff[idx] == 0xFF {
						idx++
						if buff[idx] == 0xAA {
							idx++
							if buff[idx] == 0x55 {
								idx++
								step = 1
								break
							}
						}
					}
				}
			}

			if step == 1 {
				// remove the header before going to step 2.
				buff = buff[idx+4:]
				idx = 0
			} else {
				buff = buff[idx:]
			}
		}

		// step 1: get the data size
		if step == 1 {
			if len(buff) < 4 {
				continue
			}

			// get data size (little endian)
			dataSize = int(buff[idx])
			dataSize += int(buff[idx+1]) * 0x100
			dataSize += int(buff[idx+2]) * 0x10000
			buff = buff[4:]
			idx = 0

			step = 2
		}

		// step 2: get the data
		if step == 2 {
			if len(buff) > 0 {
				if dataSize >= len(data)+len(buff) {
					data = append(data, buff...)
				} else {
					remain := dataSize - len(data)
					data = append(data, buff[:remain]...)
				}
				buff = []byte{}
				idx = 0
				if len(data) == dataSize {
					step = 3
					break
				}
			}
		}
	}

	return data, nil
}
