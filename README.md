# CARROT  

**this is just sketch of package**

Carrot is must be a database, relation nosql native database for **Go**, maybe it will be generic... 

Carrot is very very fast, faster than light! Seriously! 

	BenchmarkWrite-8        	 5000000	       374 ns/op	     165 B/op	       1 allocs/op
	BenchmarkReadFromDisk-8 	 1000000	      1039 ns/op	     171 B/op	       4 allocs/op
	BenchmarkReadFromCache-8	20000000	      92.4 ns/op	      64 B/op	       1 allocs/op


## IDEA

Database table is a just struct.

Generate code package for each struct(table), with predefined logic(write, read, search,et.c) and optimal conversation various types of data to bytes.

It should be easy to use during software development.

Structs is convenient, and why not store structs in database?

Why generic? - otherwise it will input interface{}, parse it with reflect... - it is very slow! 

## Structure

For each struct(table) create own directory

For each field create own file

For each item in table has id and information for read it from disk, ex: {id, Field0{offset,length}, Field1{offset,length}}, this allows quick read item from disk

Carrot is not only disk based database, it is also may store data in memory(cache). When reading, if item not found in cache, read it from disk, and then store to cache ^^