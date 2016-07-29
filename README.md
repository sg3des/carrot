# CARROT  

## WARNING!! this is just a concept

Carrot is must be a database, relational nosql native database for **Go**, maybe it will be generic... 

Carrot is very very fast, faster than light! Seriously! 

	BenchmarkWrite-8        	 5000000	       374 ns/op	     165 B/op	       1 allocs/op
	BenchmarkReadFromDisk-8 	 2000000	       516 ns/op	      54 B/op	       1 allocs/op
	BenchmarkReadFromCache-8	50000000	       32.4 ns/op	       0 B/op	       0 allocs/op



## IDEA

Database table is a just struct.

Code package is generated for each struct(table), with predefined logic(write, read, search,et.c) and optimal conversation various types of data to bytes.

It should be easy to use during software development.

Structs are convenient, and why not to store structs in database?

Why generic? - otherwise it will input interface{}, parse it with reflect... - it is very slow! 

## Structure

For each struct(table) it creates own directory

For each field it creates own file

There is an id and information for reading it from disk for each item in table, ex: `{id, Field0{offset,length}, Field1{offset,length}}`, this allows quick read item from disk

Carrot is not only a disk based database, it is also may store data in memory(cache). When reading, if item not found in cache, read it from disk, and then store to cache ^^
