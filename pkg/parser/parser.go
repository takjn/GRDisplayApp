package parser

import (
	"fmt"
	"io"
)

var (
	dataBuff = []byte{}
	workBuff = make([]byte, 1024*1024)
)

// GetData gets data from the given reader.
func GetData(port io.Reader) ([]byte, error) {
	var data []byte // max: 256kb
	var dataSize int

	idx := 0
	step := 0
	for step < 3 {
		// Read data
		if n, err := port.Read(workBuff); err != nil {
			return nil, fmt.Errorf("failed to read: %v", err)
		} else {
			dataBuff = append(dataBuff, workBuff[0:n]...)
		}

		// step 0: find the header: FF FF AA 55
		if step == 0 {
			if len(dataBuff) < 4 {
				continue
			}

			for idx = 0; idx < len(dataBuff)-4; idx++ {
				if dataBuff[idx] == 0xFF {
					idx++
					if dataBuff[idx] == 0xFF {
						idx++
						if dataBuff[idx] == 0xAA {
							idx++
							if dataBuff[idx] == 0x55 {
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
				dataBuff = dataBuff[idx+4:]
				idx = 0
			} else {
				dataBuff = dataBuff[idx:]
				idx = 0
			}
		}

		// step 1: get the data size
		if step == 1 {
			if len(dataBuff) < 4 {
				continue
			}

			// get data size (little endian)
			dataSize = int(dataBuff[idx])
			dataSize += int(dataBuff[idx+1]) * 0x100
			dataSize += int(dataBuff[idx+2]) * 0x10000
			dataBuff = dataBuff[4:]
			idx = 0

			step = 2
		}

		// step 2: get the data
		if step == 2 {
			if len(dataBuff) > 0 {
				if dataSize >= len(data)+len(dataBuff) {
					data = append(data, dataBuff...)
					dataBuff = []byte{}
				} else {
					remain := dataSize - len(data)
					data = append(data, dataBuff[:remain]...)
					dataBuff = dataBuff[remain:]
				}
				idx = 0
				if len(data) == dataSize {
					break
				}
			}
		}
	}

	return data, nil
}
