# Orderbook

Orderbook is used to stored all the orders and orderFragments and also responsible for 
streaming updates to its subscriber. Instead using the cache, we store the data in files which 
prevents data loss when restarting and decreases memory usage.  

