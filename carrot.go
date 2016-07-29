package carrot

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

type Info struct {
	ID     int
	Name   Offset
	Number Offset
}

type Offset struct {
	Off int64
	Len int
}

type usersInfo struct {
	write     []*Users
	writeLock bool

	indexFile *os.File
	indexLen  int

	index *indexMap
	cache *usersMap

	lastid int

	fields usersfields
}

type usersfields struct {
	Name       *os.File
	NameOffset int64

	Number       *os.File
	NumberOffset int64
}

type Users struct {
	ID      int
	Written bool
	Err     error

	Name   string
	Number int
}

var (
	users = usersInfo{
		indexLen: 25,
		index:    &indexMap{Items: make(map[int]Info)},
		cache:    &usersMap{Items: make(map[int]*Users)},
	}
)

func Open(openpath string) error {
	// users
	err := os.MkdirAll(path.Join(openpath, "users"), 0755)
	if err != nil {
		return nil
	}

	users.fields.Name, err = os.OpenFile(path.Join(openpath, "users", "name"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil
	}

	users.fields.Number, err = os.OpenFile(path.Join(openpath, "users", "number"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil
	}

	users.indexFile, err = os.OpenFile(path.Join(openpath, "users", "index"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil
	}

	parseIndex()
	go keeper()

	return nil
}

func Close() {
	for {
		if len(users.write) == 0 {
			break
		}
	}

	users.fields.Name.Close()
	users.fields.Number.Close()
	users.indexFile.Close()
}

func parseIndex() {
	indexitem := make([]byte, users.indexLen)
	var n int
	var err error
	for err == nil {
		n, err = users.indexFile.Read(indexitem)
		if n == 0 {
			continue
		}
		if n != users.indexLen || len(indexitem) != users.indexLen {
			log.Println("carrot: warning! failed parse index:", indexitem)
			continue
		}

		id := int(binary.LittleEndian.Uint32(indexitem[:4]))
		nameOff := int64(binary.LittleEndian.Uint64(indexitem[4:12]))
		nameLen := int(binary.LittleEndian.Uint16(indexitem[12:14]))
		numberOff := int64(binary.LittleEndian.Uint64(indexitem[14:22]))
		numberLen := int(binary.LittleEndian.Uint16(indexitem[22:24]))

		if Summ := nameOff + int64(nameLen); Summ > users.fields.NameOffset {
			users.fields.NameOffset = Summ
		}
		if Summ := numberOff + int64(numberLen); Summ > users.fields.NumberOffset {
			users.fields.NumberOffset = Summ
		}
		if id > users.lastid {
			users.lastid = id
		}

		users.index.Set(id, Info{ID: id, Name: Offset{nameOff, nameLen}, Number: Offset{numberOff, numberLen}})
	}

	users.fields.Name.Seek(users.fields.NameOffset, 0)
	users.fields.Number.Seek(users.fields.NumberOffset, 0)
}

func (u *Users) Write() {
	if u.ID == 0 {
		users.lastid++
		u.ID = users.lastid
	}
	id := u.ID
	users.writeLock = true
	users.write = append(users.write, u)
	users.cache.Set(id, u)
	users.writeLock = false
}

func keeper() {
	var reset bool
	for {
		time.Sleep(time.Second)
		reset = false
		for id, u := range users.write {
			if u != nil {
				if err := u.write(); err != nil {
					log.Println("carrot: failed write data, reason:", err)
				}
				users.cache.Set(u.ID, u)
				// log.Println("for write", u, writeUsersLock)
				users.write[id] = nil
			}
			if !users.writeLock {
				reset = true
			}
		}

		//fully clean writeUsers
		if reset && !users.writeLock {
			users.write = []*Users{}
		}
	}
}

func sMarshal(v string) []byte {
	return []byte(v)
}

func sUnmarshal(v []byte) string {
	return string(v)
}

func iMarshal(v int) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v))
	return b
}

func iUnmarshal(v []byte) int {
	return int(binary.LittleEndian.Uint64(v))
}

//write data to disk
func (u *Users) write() error {
	NameOff, _ := users.fields.Name.Seek(0, 1)
	NameLen, err := users.fields.Name.Write(sMarshal(u.Name))
	if err != nil {
		log.Println("failed write user")
		u.Err = err
	}

	NumberOff, _ := users.fields.Number.Seek(0, 1)
	NumberLen, err := users.fields.Number.Write(iMarshal(u.Number))
	if err != nil {
		log.Println("failed write number")
		u.Err = err
	}

	users.index.Set(u.ID, Info{ID: u.ID, Name: Offset{NameOff, NameLen}, Number: Offset{NumberOff, NumberLen}})

	bID := make([]byte, 4)
	binary.LittleEndian.PutUint32(bID, uint32(u.ID))

	bNameOff := make([]byte, 8)
	binary.LittleEndian.PutUint64(bNameOff, uint64(NameOff))

	bNameLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(bNameLen, uint16(NameLen))

	bNumberOff := make([]byte, 8)
	binary.LittleEndian.PutUint64(bNumberOff, uint64(NumberOff))

	bNumberLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(bNumberLen, uint16(NumberLen))

	buf := bytes.NewBuffer([]byte{})
	buf.Write(bID)
	buf.Write(bNameOff)
	buf.Write(bNameLen)
	buf.Write(bNumberOff)
	buf.Write(bNumberLen)
	buf.Write(make([]byte, 1))

	_, err = buf.WriteTo(users.indexFile)
	if err != nil {
		return err
	}

	u.Written = true

	users.fields.NameOffset = NameOff + int64(NameLen)
	users.fields.NumberOffset = NumberOff + int64(NumberLen)

	return nil
}

func (u *Users) Read(id int) error {
	if u == nil {
		u = new(Users)
	}

	item, ok := users.cache.Get(id)
	if ok {
		*u = *item
		return nil
	}

	info, ok := users.index.Get(id)
	if !ok {
		return errors.New(fmt.Sprintf("item by id %d not found", id))
	}

	var bName = make([]byte, info.Name.Len)
	users.fields.Name.ReadAt(bName, info.Name.Off)
	u.Name = sUnmarshal(bName)

	var bNumber = make([]byte, info.Number.Len)
	users.fields.Number.ReadAt(bNumber, info.Number.Off)
	u.Number = iUnmarshal(bNumber)

	users.cache.Set(id, u)
	return nil
}

func (u *Users) ClearCache() {
	users.cache.Truncate()
}
