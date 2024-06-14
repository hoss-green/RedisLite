# RedisLite
A small compact Redis written in go.

### Origin
This started as a small project challenge on codecrafters to test my skills with GO, I completed the initial challenge in 5 days. I'm in the process of polishing, tidying up, and adding in the missing string and some core functions. 

### Roadmap
1. Complete the String functions
2. Add in the core general commands DEL, COPY, EXISTS
3. Branch and create an event loop in branch to give options on execution styles
4. Fix the memory leak with not deleted expired keys

### Current Compatability
#### GENERAL
| Command         | Implementation | Missing
|-----------------|----------------|-----------------------
| DEL             | 0              | 100 
| COPY            | 0              | 100 
| EXISTS          | 0              | 100 
|                 |                |
| __TOTAL__ (3)   | 0%             |

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
| INCRBYFLOAT     | 100            | 
| LCS             | 0              | 100 
| MGET            | 100            | 
| MSET            | 100            | 
| SET             | 100            | 
| SETRANGE        | 0              | 100
| STRLEN          | 100            | 
| SUBSTR          | 100            | 
|                 |                |
| __TOTAL__ (18)  | 88.9%          |

#### Notes (dnf or unfinished)
| Command      | Notes  
|--------------|-----------------------
| DECR         | calculation overflow allowed (dnf)
| DECRBY       | calculation overflow allowed (dnf)
| INCR         | calculation overflow allowed (dnf)
| INCRBY       | calculation overflow allowed (dnf)
| INCRBYFLOAT  | calculation overflow allowed (dnf)
| LCS          | 
| SETRANGE     | 
