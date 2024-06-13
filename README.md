# RedisLite
A small compact Redis written in go.

### Origin
This started as a small project challenge on codecrafters to test my skills with GO. I'm in the process of polishing, tidying up, and adding in the extra functions. 

### Roadmap
1. Complete the String functions

### Current Compatability
#### STRING
| Command         | Implementation | Missing
|-----------------|----------------|-----------------------
| APPEND          | 100            | None
| DECR            | 90             | 
| DECRBY          | 90             | 
| GET             | 90             | None
| GETDEL          | 100            | 
| GETEX           | 0              | 
| GETRANGE        | 100            | None
| GETSET          | 100            | 
| INCR            | 90             | 
| INCRBY          | 90             | 
| INCRBYFLOAT     | 0              | 
| LCS             | 0              | 
| MGET            | 100            | 
| MSET            | 100            | 
| SET             | 85             | XX NX KEEPTTL GET
| SETRANGE        | 0              | 
| STRLEN          | 90             | 
| SUBSTR          | 100            | None
|                 |                |
| __TOTAL__ (18)  | 74.1%          |

#### Notes (open todo or unfinished)
| Command      | Notes  
|--------------|-----------------------
| DECR         | Allows overflow
| DECRBY       | Allows overflow
| GET          | Does not return an error if the type is wrong 
| GETEX        | 
| INCR         | Allows overflow 
| INCRBY       | Allows overflow
| INCRBYFLOAT  | 
| LCS          | 
| SET          | Missing Parameter options 
| SETRANGE     | 
| STRLEN       | Does not return an error if the type is wrong 
