# Orderbook

Orderbook is used to stored all the orders and orderFragments and also responsible for 
streaming updates to its subscriber. Instead using the cache, we store the data in files which 
prevents data loss when restarting and decreases memory usage.  

## Tables 

**Status**

The `status` table maps the orderID to its current status.
We don't need to encode the status as the order status is of type uint8.

**OrderFragment**

The `orderFragments` table maps the orderID to the actual orderFragment we have.
We encode the fragment in JSON format and store in the db. 
For now,  the data will not be encrypted nor synchronized.

**Order**

The `order` table maps the orderID to the actual order.
We encode the order in JSON format and store in the db. 
For now,  the data will not be encrypted nor synchronized.

**Atom**

The `atom` table maps the orderID to the atom message needed in the atomic swap.
 
