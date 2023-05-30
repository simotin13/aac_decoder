package main

import (
	"encoding/binary"
	fmt "fmt"
	os "os"
)

func main() {
	argc := len(os.Args)
	if argc < 3 {
		fmt.Println("Usage: ./aac_decoder <input.aac> <output.aac>")
		os.Exit(-1)
	}
	filepath := os.Args[1]
	fi, _ := os.Stat(filepath)

	f, err := os.Open(filepath)
	if err != nil {
		os.Exit(-1)
	}
	defer f.Close()

	bin := make([]byte, fi.Size())
	f.Read(bin)

	wf, err := os.Create(os.Args[2])
	if err != nil {
		panic(err)
	}

	leftLength := len(bin)
	var adtsHeaderTop uint32 = 0
	var binSize uint32 = 0
	for 0 < leftLength {
		//fmt.Printf("adtsHeaderTop:0x%X, leftLength:0x%X\n", adtsHeaderTop, leftLength)

		// check sync bytes(0xFFF)
		if bin[adtsHeaderTop] != 0xFF || bin[adtsHeaderTop+1]&0xF0 != 0xF0 {
			fmt.Printf("Not ADTS format: 0x%x%x\n", bin[adtsHeaderTop], bin[adtsHeaderTop+1])
			os.Exit(1)
		}

		/*
			mpegVer := (bin[adtsHeaderTop+1] & 0x0F) & 0x80
					if mpegVer == 0 {
						fmt.Println("MPEG Ver: MPEG-4")
					} else {
						fmt.Println("MPEG Ver: MPEG-2")
					}
				layer := (bin[adtsHeaderTop+1] & 0x0F) & 0x60
				if layer != 0 {
					msg := fmt.Sprintf("Layer: %d", layer>>1)
					panic(msg)
				}

				if bin[adtsHeaderTop+1]&0x01 == 0x00 {
					fmt.Println("CRC Protection Absent: Yes")
				} else {
					fmt.Println("CRC Protection Absent: No")
				}
		*/

		by := bin[adtsHeaderTop+2]
		/*
			if (by&0xC0)>>6 == 0 {
				fmt.Println("MPEG-4 Audio Main")
			} else if (by&0xC0)>>6 == 1 {
				fmt.Println("MPEG-4 Audio LC")
			} else if (by&0xC0)>>6 == 2 {
				fmt.Println("MPEG-4 Audio SSR")
			} else {
				fmt.Printf("by: 0x%X\n", (by&0xC0)>>6)
				fmt.Println("MPEG-4 Audio LTP(Reserved)")
			}
			samplingRate := ((by & 0x3C) >> 2)
			switch samplingRate {
			case 0:
				fmt.Println("Sampling Rate: 96000 Hz")
			case 1:
				fmt.Println("Sampling Rate: 88200 Hz")
			case 2:
				fmt.Println("Sampling Rate: 64000 Hz")
			case 3:
				fmt.Println("Sampling Rate: 48000 Hz")
			case 4:
				fmt.Println("Sampling Rate: 44100 Hz")
			case 5:
				fmt.Println("Sampling Rate: 32000 Hz")
			case 6:
				fmt.Println("Sampling Rate: 24000 Hz")
			case 7:
				fmt.Println("Sampling Rate: 22050 Hz")
			case 8:
				fmt.Println("Sampling Rate: 16000 Hz")
			case 9:
				fmt.Println("Sampling Rate: 12000 Hz")
			case 10:
				fmt.Println("Sampling Rate: 11025 Hz")
			case 11:
				fmt.Println("Sampling Rate: 8000 Hz")
			case 12:
				fmt.Println("Sampling Rate: 7350 Hz")
			default:
				fmt.Printf("Sampling Rate: %d Hz\n", samplingRate)
			}
		*/
		/*
			pvtStream := (by & 0x02) >> 1
			if pvtStream == 1 {
				fmt.Println("PrivateSteam: Decoding")
			} else {
				fmt.Println("PrivateSteam: Encoding")
			}

		*/
		channel := (by & 0x01) << 2
		by = bin[adtsHeaderTop+3]

		channel |= (by & 0xC0) >> 6
		/*
			switch channel {
			case 0:
				fmt.Println("Channel: Defined in AOT Specifc Config")
			case 1:
				fmt.Println("Channel: 1 channel: front-center")
			case 2:
				fmt.Println("Channel: 2 channels: front-left, front-right")
			case 3:
				fmt.Println("Channel: 3 channels: front-center, front-left, front-right")
			case 4:
				fmt.Println("Channel: 4 channels: front-center, front-left, front-right, back-center")
			case 5:
				fmt.Println("Channel: 5 channels: front-center, front-left, front-right, back-left, back-right")
			case 6:
				fmt.Println("Channel: 6 channels: front-center, front-left, front-right, back-left, back-right, LFE-channel")
			case 7:
				fmt.Println("Channel: 8 channels: front-center, front-left, front-right, side-left, side-right, back-left, back-right, LFE-channel")
			default:
				fmt.Printf("Channel: Reserved(%d)\n", channel)
			}
		*/

		/*
			org := (by & 0x20) >> 5
				if org == 1 {
					fmt.Println("Original: Decoding")
				} else {
					fmt.Println("Original: Encoding")
				}
				home := (by & 0x10) >> 4
				if home == 1 {
					fmt.Println("Home: Decoding")
				} else {
					fmt.Println("Home: Encoding")
				}
				copyRightStream := (by & 0x08) >> 3
				if copyRightStream == 1 {
					fmt.Println("CopyRight: Decoding")
				} else {
					fmt.Println("CopyRight: Encoding")
				}
				copyRightStart := (by & 0x04) >> 2
				if copyRightStart == 1 {
					fmt.Println("CopyRightStart: Decoding")
				} else {
					fmt.Println("CopyRightStart: Encoding")
				}
		*/
		var frameLength uint16 = 0
		frameLength = uint16((by & 0x03)) << 11
		frameLength |= uint16(bin[adtsHeaderTop+4]) << 3
		frameLength |= uint16((bin[adtsHeaderTop+5] & 0xE0) >> 5)
		/*
			fmt.Printf("FrameLength: %d(bytes)\n", frameLength)

			bufUsage := (bin[adtsHeaderTop+5] & 0x1F) << 6
			bufUsage |= (bin[adtsHeaderTop+6] & 0xFC) >> 2
			fmt.Printf("BufUsage: %d\n", bufUsage)
		*/
		// Raw Data Block
		/*
			rdb := bin[adtsHeaderTop+6] & 0x03
				fmt.Printf("RDB: %d\n", rdb)
		*/
		var start uint32 = adtsHeaderTop
		var end uint32 = adtsHeaderTop + uint32(frameLength)
		binSize += end - start
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, uint32(frameLength))
		wf.Write(bytes)
		wf.Write(bin[start:end])

		adtsHeaderTop += uint32(frameLength)
		leftLength -= int(frameLength)
	}
	os.Exit(0)
}
