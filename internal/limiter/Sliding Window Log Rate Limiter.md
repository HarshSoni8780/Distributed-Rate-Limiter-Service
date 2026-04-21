

**Core Concept** The Sliding Window Log treats time as a continuous flow. Instead of counting in rigid blocks, it records every request with a timestamp. This creates a "moving window" that ensures the limit is never exceeded in any 60-second slice of time, solving the boundary burst problem.

**Key Components**

- **Redis Sorted Set (ZSet):** The primary data structure. It uses the timestamp as the "score" to keep requests in chronological order.
    
- **Nanosecond Precision:** Uses UnixNano for timestamps to ensure every request has a unique identity in the set, preventing data overwrites.
    
- **Proactive Cleanup:** Every check starts by deleting old data, which "slides" the window forward.
    

**Logic Flow**

1. **The Slide:** Removes all timestamps from Redis that are older than the window duration (e.g., Now minus 10 seconds).
    
2. **The Count:** Uses the `ZCard` command to see how many valid timestamps remain in the set.
    
3. **The Gatekeeper:** Compares the remaining count to the allowed limit. If the limit is reached, the request is rejected immediately.
    
4. **The Log:** If allowed, the current timestamp is added to the set to count against future requests.
    
5. **The Maintenance:** Sets an expiration on the whole key so Redis clears the memory if the user becomes inactive.
    

**Advantages**

- **Extreme Accuracy:** Mathematically prevents bursts at the edge of a minute.
    
- **Fairness:** Users only have to wait for their specific oldest request to expire, not for a clock reset.
    
- **Flexibility:** Easily scales to different durations (1s, 10s, 1min) just by changing the duration variable.
    

**Disadvantages**

- **Memory Cost:** Stores every request timestamp rather than a single number.
    
- **Performance:** Operations like `ZAdd` and `ZRem` take O(logN) time, which is slightly heavier than simple increments.
    

**Example Trace (Limit 2, Window 10s)**

- **0s:** Request 1 arrives. **Allowed**. (Set: [0s])
    
- **5s:** Request 2 arrives. **Allowed**. (Set: [0s, 5s])
    
- **8s:** Request 3 arrives. Window is 0s to 10s. Count is 2. **Blocked**.
    
- **11s:** Request 4 arrives. Cleanup removes 0s. Count is now 1. **Allowed**.