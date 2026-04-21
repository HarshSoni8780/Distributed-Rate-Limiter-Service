## Rate Limiter: Fixed Window Summary

**Core Concept** The Fixed Window algorithm divides time into 1-minute blocks. It uses Redis as a centralized counter to track how many requests a user makes within that specific block.

**Key Components**

- **Struct:** Bundles the Redis client and the limit into one object.
    
- **Factory Function:** A constructor that returns a pointer to the struct for efficient memory use.
    
- **Window Logic:** Uses integer division of the Unix timestamp by 60 to group all seconds in a minute into one key.
    

**Logic Flow**

1. **Unique Key:** Creates a key using the UserID and the current window number.
    
2. **Atomic Increment:** Uses the Redis `INCR` command to ensure the count is accurate even with high concurrent traffic.
    
3. **Self-Cleaning:** Checks if `count == 1` and sets a 60-second expiration. This prevents memory leaks by deleting old keys.
    
4. **Decision:** Compares the current count against the defined limit.
    

**Advantages**

- **Distributed:** Works across multiple server instances via Redis.
    
- **Speed:** Constant time complexity O(1).
    
- **Efficiency:** Minimal CPU and memory overhead.
    

**Disadvantages**

- **Burst Vulnerability:** Allows double the limit if requests are sent at the edge of two windows.
    

**Example Trace** Limit = 2

- **Request 1:** Count becomes 1. Timer set. **Allowed**.
    
- **Request 2:** Count becomes 2. **Allowed**.
    
- **Request 3:** Count becomes 3. **Blocked** (Rate Limited).