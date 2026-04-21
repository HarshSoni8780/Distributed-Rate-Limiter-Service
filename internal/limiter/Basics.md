
## 1) Concurrent Rate Limiter

**Definition:**  
Limits the number of **requests in progress at the same time**, instead of requests per second.
### Why needed

- Some APIs are **slow / resource-heavy**
- Users retry when response is slow → increases load
### Example

- Limit = 20 concurrent
- 21st request → rejected (429) or waits
### Use cases

- Payment APIs (e.g., Stripe)

gees