package carrot

import (
	"log"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile)
	Open()
}

func BenchmarkWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var u = &Users{Name: "name2", IP: 123837}
		u.Write(i)
	}

}

func BenchmarkPause(b *testing.B) {
	b.Log("some pause, as Go not safety for concurrent write and read, and it`s very BAD! but ")
	time.Sleep(20 * time.Second)
}

func BenchmarkRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Read(i)
	}
}
