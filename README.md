# Distributed Go

## About
This repository follows along the content discussed in Travis Jeffrey's Distributed Services with Go.

## Concepts

### Protobuf 
A language Google designed to structure data that flows from different components. It is used to define a standardized schema amongst consumers and producers in a microservice architecture. Not to mention, it is performant compared to JSON.


Benefits:
- Type safety
- Prevents schema violations
- Language agnostic (can compile into different languages)
- Supports versioning and backwards compatibility
- Prevents the need for version checks
- Handles encoding and decoding
- Allows you to define plugins

### Logs
Logs are append-only and are sorted by insertion time typically.

### Write-Ahead Logs (WAL)
A log that denotes all the changes that you want made to something (i.e., a database, replicas of DB's, clusters, etc.) before applying those changes. 


### Segments
A portion of a log. This is a way of optimizing the memory usage when you have a basic log that only appends. By splitting the log into segments, you can free up disk space by removing old segments that are "processed." There is a special segment called the *active segment* which is the segment that is currently being written to.

Segments consist of store files and index files. 

#### Store File
A file that contains all the recorded data

#### Index file
A file that contains the index for a given record.

Index files consist of two things usually, a offset and the stored position of that record. You use the offset to read the recorded # of bytes from the stored position to fetch the entire record.