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
| APPEND          | 100            | 
| DECR            | 100            | 
| DECRBY          | 100            | 
| GET             | 100            |
| GETDEL          | 100            | 
| GETEX           | 0              | 100 
| GETRANGE        | 100            | 
| GETSET          | 100            | 
| INCR            | 100            | 
| INCRBY          | 100            | 
| INCRBYFLOAT     | 0              | 100 
| LCS             | 0              | 100 
| MGET            | 100            | 
| MSET            | 100            | 
| SET             | 85             | XX NX KEEPTTL GET
| SETRANGE        | 0              | 100
| STRLEN          | 100            | 
| SUBSTR          | 100            | 
|                 |                |
| __TOTAL__ (18)  | 76.95%         |

#### Notes (open todo or unfinished)
| Command      | Notes  
|--------------|-----------------------
| GETEX        | 
| INCRBYFLOAT  | 
| LCS          | 
| SET          | Missing Parameter options 
| SETRANGE     | 
