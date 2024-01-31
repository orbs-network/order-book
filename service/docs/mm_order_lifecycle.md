### Redis keys per order

1. clientOId:<ID>:order
2. orderID:<ID>:order
3. <SYMBOL>:<buy/sell>:prices (for storing bid/ask min/max prices)
4. userId:<ID>:openOrders
5. userId:<ID>:filledOrders (only for filled orders)

### Current lifecycle

1. Created, not filled at all, then cancelled -> all keys are deleted
2. Created, locked, then (attempted) cancel -> denied due to pending fill
3. Created, partial filled (so no longer locked), cancelled -> all keys are deleted
4. Created, filled, then (attempted) cancelled -> denied (on fill, order removed from `:prices`, order removed from `:openOrders`, added to `:filledOrders`)

### Future lifecycle

1. Created, not filled at all, then cancelled -> update order `cancelled` true, order removed from `:prices`, order removed from `:openOrders`
2. Created, locked, then (attempted) cancel -> denied due to pending fill
3. Created, partial filled (so no longer locked), cancelled -> update order `cancelled` true, order removed from `:prices`, order removed from `:openOrders`, add to `:filledOrders`
4. Created, filled, then (attempted) cancelled -> denied (on fill, update order `cancelled` true, order removed from `:prices`, order removed from `:openOrders`, added to `:filledOrders`)
