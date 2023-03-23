package qrnfunc

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const bufferSize = 16384 *3

func GenerateQrnValue(filePrefix string) uint64 {

	var ioldIndex byte

	index := uint64(0)
	randomValue := make([]byte, 8)

	for {
		buffer := readUpFile(filePrefix)
		if buffer[0] == ioldIndex {
			writeDnFile(filePrefix, buffer[0])
			time.Sleep(1 * time.Second)
		} else {
			ioldIndex = buffer[0]

			if len(buffer) > 0 {
				for j := 1; j <= bufferSize-1; j++ {
					fmt.Printf("[0x%02x]", buffer[j])
					randomValue[index] = buffer[j]
					index++
					if (index == 8) {
						return uint64(binary.LittleEndian.Uint64(randomValue))
					}
				}
				fmt.Println()

				writeDnFile(filePrefix, buffer[0])
			} else {
				time.Sleep(1 * time.Second)
			}

			fmt.Println(fmt.Sprintf("INDEX: %d", ioldIndex))
		}
	}
}

func readUpFile(filePrefix string) []byte {
	buffer, err := ioutil.ReadFile(filePrefix + "up.ini")
	if err != nil {
		fmt.Println("OPEN-ERROR: up.ini")
		panic(err.Error())
	}

	return buffer
}

func writeDnFile(filePrefix string, data byte) {
	os.Remove(filePrefix + "dn.ini")
	err := ioutil.WriteFile(filePrefix+"dn.ini", []byte{data}, 0644)
	if err != nil {
		fmt.Println("CREATE-ERROR: dn.ini")
		panic(err.Error())
	}
}

