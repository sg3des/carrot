package carrot

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"time"
)

type Users struct {
	id  int
	err error

	Name string
	IP   int
}

var (
	writeUsers     []*Users
	writeUsersLock bool

	indexFile    *os.File
	indexLineLen = 24

	usersNameOffset int64
	usersName       *os.File

	usersIPOffset int64
	usersIP       *os.File
)

type Info struct {
	ID   int
	Name Offset
	IP   Offset
}

type Offset struct {
	Off int64
	Len int
}

func Open() {
	var err error
	usersName, err = os.OpenFile("./db/users/name", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	usersIP, err = os.OpenFile("./db/users/ip", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	indexFile, err = os.OpenFile("./db/users/index", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	parseIndex()

	go keeper()
}

func parseIndex() {
	indexitem := make([]byte, indexLineLen)
	var n int
	var err error
	for err == nil {
		n, err = indexFile.Read(indexitem)
		if n == 0 {
			continue
		}
		if n != indexLineLen || len(indexitem) != indexLineLen {
			log.Println("WTF?", len(indexitem), n, indexLineLen)
			continue
		}

		id := int(binary.LittleEndian.Uint32(indexitem[:4]))
		nameOff := int64(binary.LittleEndian.Uint64(indexitem[4:12]))
		nameLen := int(binary.LittleEndian.Uint16(indexitem[12:14]))
		ipOff := int64(binary.LittleEndian.Uint64(indexitem[14:22]))
		ipLen := int(binary.LittleEndian.Uint16(indexitem[22:24]))

		if Summ := nameOff + int64(nameLen); Summ > usersNameOffset {
			usersNameOffset = Summ
		}
		if Summ := ipOff + int64(ipLen); Summ > usersIPOffset {
			usersIPOffset = Summ
		}

		index.Set(id, Info{ID: id, Name: Offset{nameOff, nameLen}, IP: Offset{ipOff, ipLen}})
	}

	usersName.Seek(usersNameOffset, 0)
	usersIP.Seek(usersIPOffset, 0)
}

func (u *Users) Write(id int) {
	u.id = id
	writeUsersLock = true
	writeUsers = append(writeUsers, u)
	// writeUsers.Set(id, u)
	cacheUsers.Set(id, u)
	writeUsersLock = false
}

func keeper() {
	for {
		time.Sleep(time.Second)
		var reset bool
		for _, u := range writeUsers {
			if u != nil {
				if err := u.write(); err != nil {
					log.Println("carrot: failed write data, reason:", err)
				}
				// log.Println("for write", u, writeUsersLock)
				u = nil
			}
			if !writeUsersLock {
				reset = true
			}
		}
		if reset {
			writeUsers = []*Users{}
		}
		// writeUsers.Lock()
		// for id, u := range writeUsers.items {
		// 	if err := u.write(); err != nil {
		// 		log.Println("failed write")
		// 		continue
		// 	}
		// 	delete(writeUsers.items, id)
		// }
		// writeUsers.Unlock()
	}
}

func (u *Users) nameBytes() []byte {
	return []byte(u.Name)
}

func (u *Users) ipBytes() []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(u.IP))
	return b
}

func (u *Users) write() error {
	// log.Println("write")
	nameOff, _ := usersName.Seek(0, 1)
	nameLen, err := usersName.Write(u.nameBytes())
	if err != nil {
		log.Println("failed write user")
		u.err = err
	}

	ipOff, _ := usersIP.Seek(0, 1)
	ipLen, err := usersIP.Write(u.ipBytes())
	if err != nil {
		log.Println("failed write ip")
		u.err = err
	}

	index.Set(u.id, Info{ID: u.id, Name: Offset{nameOff, nameLen}, IP: Offset{ipOff, ipLen}})

	bid := make([]byte, 4)
	binary.LittleEndian.PutUint32(bid, uint32(u.id))

	bnameOff := make([]byte, 8)
	binary.LittleEndian.PutUint64(bnameOff, uint64(nameOff))

	bnameLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(bnameLen, uint16(nameLen))

	bipOff := make([]byte, 8)
	binary.LittleEndian.PutUint64(bipOff, uint64(ipOff))

	bipLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(bipLen, uint16(ipLen))

	buf := bytes.NewBuffer([]byte{})
	buf.Write(bid)
	buf.Write(bnameOff)
	buf.Write(bnameLen)
	buf.Write(bipOff)
	buf.Write(bipLen)

	_, err = buf.WriteTo(indexFile)
	if err != nil {
		return err
	}

	usersNameOffset = nameOff + int64(nameLen)
	usersIPOffset = ipOff + int64(ipLen)

	return nil
}

func Read(id int) *Users {
	item, ok := cacheUsers.Get(id)
	if ok {
		return item
	}
	// log.Println("return from disk")

	info, ok := index.Get(id)
	if !ok {
		return nil
	}

	var u = new(Users)

	var nameBytes = make([]byte, info.Name.Len)
	usersName.ReadAt(nameBytes, info.Name.Off)
	u.Name = string(nameBytes)

	var ipBytes = make([]byte, info.IP.Len)
	usersIP.ReadAt(ipBytes, info.IP.Off)
	_ip, _ := binary.ReadVarint(bytes.NewReader(ipBytes))
	u.IP = int(_ip)

	cacheUsers.Set(id, u)

	return u
}

func ClearCache() {
	cacheUsers.Truncate()
}
