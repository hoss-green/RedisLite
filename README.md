# RedisLite
A small compact Redis written in go.

### Origin
This started as a small project challenge on codecrafters to test my skills with GO. I'm in the process of polishing, tidying up, and adding in the extra functions. 

### Roadmap
1. Complete the String functions
2. Branch and create an event loop to give options on execution styles
3. Fix the memory leak with not deleted expired keys

### Current Compatability
#### STRING
| Command         | Implementation | Missing
|-----------------|----------------|-----------------------
| APPEND          | 100            | 
| DECR            | 100            | 
| DECRBY          | 100            | 
| GET             | 100            |
| GETDEL          | 100            | 
| GETEX           | 100            | 
| GETRANGE        | 100            | 
| GETSET          | 100            | 
| INCR            | 100            | 
| INCRBY          | 100            | 
| INCRBYFLOAT     | 0              | 100 
| LCS             | 0              | 100 
| MGET            | 100            | 
| MSET            | 100            | 
| SET             | 100            | 
| SETRANGE        | 0              | 100
| STRLEN          | 100            | 
| SUBSTR          | 100            | 
|                 |                |
| __TOTAL__ (18)  | 83.3%          |

#### Notes (open todo or unfinished)
| Command      | Notes  
|--------------|-----------------------
| GETEX        | 
| INCRBYFLOAT  | 
| LCS          | 
| SETRANGE     | 
