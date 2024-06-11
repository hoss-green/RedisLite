# RedisLite
A small compact Redis written in go.

### Origin
This started as a small project challenge on codecrafters to test my skills with GO. I'm in the process of polishing, tidying up, and adding in the extra functions. The goals being, in order:

1. Complete the String functions
2. Add more complex types
3. Save to file (*.rdb) read and write functionality for strings

### Current Compatability
#### STRING
| Command      | Implementation | Missing
| APPEND       | Full           | None
| DECR         | 0              | 
| DECRBY       | 0              | 
| GET          | Full           | None
| GETDEL       | 0              | 
| GETEX        | 0              | 
| GETRANGE     | 0              | 
| GETSET       | 0              | 
| INCR         | 0              | 
| INCRBY       | 0              | 
| INCRBYFLOAT  | 0              | 
| LCS          | 0              | 
| MGET         | 0              | 
| MSET         | 0              | 
| SET          | Partial        | XX NX KEEPTTL GET
| SETRANGE     | 0              | 
| STRLEN       | 0              | 
| SUBSTR       | 0              | 
