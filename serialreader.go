package main

import (
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"github.com/wcharczuk/go-chart"
	"os"
	"bufio"
)

func draw_graph(x, y []float64, filename, xlabel, ylabel string) {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      xlabel,
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:      ylabel,
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: x,
				YValues: y,
			},
		},
		Height: 1000,
		Width:  2000,
	}
	file, err := os.Create(filename)
	fmt.Println(err)
	writer := bufio.NewWriter(file)
	graph.Render(chart.PNG, writer)
	writer.Flush()
	file.Close()
}

func reprint_graph(counts []int, count int) {
	x := make([]float64, 511)
	y := make([]float64, 511)
	for i := 1; i < 512; i++ {
			x[i-1] = float64(i)
			y[i-1] = 512.0 * float64(counts[i]) / float64(count)
	}
	xlabel := "Bin Number"
	ylabel := "Portion"
	filename := "chart.png"
	draw_graph(x, y, filename, xlabel, ylabel)
}

func grapher(counts []int, count *int){
	reader := bufio.NewReader(os.Stdin)
	for {
		if text, _ := reader.ReadString('\n'); len(text) > 0 {
			reprint_graph(counts, *count)
		}
	}
}

func main() {
	fmt.Println("Hello World")
	// Set up options.
	options := serial.OpenOptions{
	PortName: "/dev/ttyUSB0",
	BaudRate: 115200,
	DataBits: 8,
	StopBits: 1,
	MinimumReadSize: 1,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
	fmt.Printf("serial.Open: %v\n", err)
	}

	// Make sure to close it later.
	defer port.Close()
	
	count := 0
	counts := make([]int, 512)
	
	go grapher(counts, &count)
	
	for {
		b := make([]byte, 2)
		port.Read(b)
		idx := int(b[0]) * 256 + int(b[1])
		if idx == int(0x1111) {
			fmt.Print("\nSwapped 0x1111 to 0x100 ")
			idx = 0x100
		} else if idx == 0x100 {
			fmt.Print("\nFinding next byte\n")
			b2 := make([]byte, 1)
			b2[0] = 0
			for ; b2[0] == 0; port.Read(b2) {
				
			}
			idx = idx + int(b2[0])
			//port.Read(b2)
		} else {
			fmt.Print("\nDidn't swap ")
		}
		if idx >= 512 {
			if (idx & 0xFF00) == idx {
				idx = idx >> 8
				fmt.Printf("\nCoverted %x to %x(%v)\n", b, idx, idx);
				counts[idx]++
				count++
			} else {
				fmt.Printf("\nError: idx more than 512: %v, %x\n", idx, b)
			}
		} else {
			fmt.Printf("%x(%v,%x),", b, idx, idx)
			counts[idx]++
			count++
		}
	}
}
