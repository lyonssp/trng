# TRNG

TRNG generates a pool of entropy, similar to /dev/random

## Testing

### Chi Squared Test for Uniform Distribution

The distribution of bit values of the entropy pool has been tested using the
Chi Squared method to determine uniform distribution of bit values

### Compression Testing

While it may be properly statistically distributed, the pool can still be
predictable.  For example, the integer sequence [0,1,2,3,4,5,6,7,8,9] is
uniformly distributed, but wouldn't pass as a random pool of data.  Compression
tools attempt to find patterns in data in order to encode similar blobs of data with
less bytes, so this heuristic measure helps us to gain some confidence in the fact
that our stream of bits is not predictable.
