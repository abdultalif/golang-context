package golangcontext

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	// cara membuat context
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")
	contextG := context.WithValue(contextF, "g", "G")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
	fmt.Println(contextG)

	fmt.Println(contextF.Value("f"))
	fmt.Println(contextF.Value("c"))
	fmt.Println(contextF.Value("b"))

	fmt.Println(contextA.Value("b"))
}

func CreateCounter(ctx context.Context) chan int  {
	destination := make(chan int)
	
	go func() {
		defer close(destination)
		counter := 1
		for  {
			select {
				case <- ctx.Done():
					return
				default:
					destination <- counter
					counter++
					time.Sleep(1 * time.Second)
			}
		} 
	}()
			return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}

	cancel() // membatalkan context

	// sengaja make timer  untuk memastikan goroutine nya mati, karena jika tidak menggunakan timer keburu eksekusi println nya duluan karen goroutine menggunakan async (paralel)
	time.Sleep(3 * time.Second)

	fmt.Println("Total goroutine", runtime.NumGoroutine())
}

func TestContextWithTimeOut(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	// setelah 5 detik context nya selesai
	ctx, cancel := context.WithTimeout(parent, 5 * time.Second)

	defer cancel()

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
	}


	// sengaja make timer  untuk memastikan goroutine nya mati, karena jika tidak menggunakan timer keburu eksekusi println nya duluan karen goroutine menggunakan async (paralel)
	time.Sleep(3 * time.Second)

	fmt.Println("Total goroutine", runtime.NumGoroutine())
}


func TestContextWithDeadLine(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	// context akann selesai sesuai custom deadline
	ctx, cancel := context.WithDeadline(parent, time.Date(
		2025, 12, 31,  // Tahun, Bulan, Hari
		23, 59, 59,    // Jam, Menit, Detik
		0,             // Nanoseconds
		time.Local,      // Timezone (bisa pakai time.Local jika mau WIB)
	))

	defer cancel()

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
	}


	// sengaja make timer  untuk memastikan goroutine nya mati, karena jika tidak menggunakan timer keburu eksekusi println nya duluan karen goroutine menggunakan async (paralel)
	time.Sleep(3 * time.Second)

	fmt.Println("Total goroutine", runtime.NumGoroutine())
}

