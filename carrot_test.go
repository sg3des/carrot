package carrot

import (
	"log"
	"testing"
	"time"
)

var item *Users

func init() {
	log.SetFlags(log.Lshortfile)
	Open("db")
}

func BenchmarkWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var u = &Users{Name: "name2", Number: i}
		u.Write()
	}

}

func BenchmarkPause(b *testing.B) {
	time.Sleep(15 * time.Second)
	item.ClearCache()
}

func BenchmarkReadFromDisk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		item.Read(i)
	}
}

func BenchmarkReadFromCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		item.Read(1)
	}
}
