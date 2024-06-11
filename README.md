# RedisLite
A small compact Redis written in go.

### Origin
This started as a small project challenge on codecrafters to test my skills with GO. I'm in the process of polishing, tidying up, and adding in the extra functions. 

### Roadmap
1. Complete the String functions

### Current Compatability
#### STRING
| Command      | Implementation | Missing
|--------------|----------------|-----------------------
| APPEND       | 100            | None
| DECR         | 0              | 
| DECRBY       | 0              | 
| GET          | 100            | None
| GETDEL       | 0              | 
| GETEX        | 0              | 
| GETRANGE     | 100            | None
| GETSET       | 0              | 
| INCR         | 0              | 
| INCRBY       | 0              | 
| INCRBYFLOAT  | 0              | 
| LCS          | 0              | 
| MGET         | 0              | 
| MSET         | 0              | 
| SET          | 85             | XX NX KEEPTTL GET
| SETRANGE     | 0              | 
| STRLEN       | 90             | 
| SUBSTR       | 100            | None

#### Notes (open todo or unfinished)
| Command      | Notes  
|--------------|-----------------------
| DECR         | 
| DECRBY       | 
| GETDEL       | 
| GETEX        | 
| GETSET       | 
| INCR         | 
| INCRBY       | 
| INCRBYFLOAT  | 
| LCS          | 
| MGET         | 
| MSET         | 
| SET          | Missing Parameter options 
| SETRANGE     | 
| STRLEN       | Does not return an error if the type is wrong 
