## Pending Orders Sequence Diagram

```mermaid
sequenceDiagram
UC ->> DEX: getPrice(not related to us)
DEX ->> LH: getQuote(amountIn)
LH ->> LH: startAuction()
LH ->> OB: getQuote(auctionId, amountIn)
OB ->> LH: return(MM-signed-Order[], amountOut)
OB ->> OB: allowQuoteOrdersToCancel(auctionId)
LH ->> LH: AuctionEnds(winner=OB)
LH ->> UC: requestUserSign(amountOut)
UC ->> LH: return(userPermit)
LH ->> OB: approveOrders(auctionId)
OB ->> OB: checkIfOrderStillInBook()
OB ->> OB: markOrdersAsPending()
OB ->> LH: return(obSignature)
LH ->> WM: sendTx(MM-sig, USER-sig OB-sig, TTL...?)
WM ->> LH: mined()
LH ->> OB: txCompleted(auctionID)
```

## Endpoints (MVP)

### MM

1. addOrder(order) POST - create a new order
2. cancelOrder(orderId) DELETE - cancel an order
3. orders() GET (returns Order[] for a given MM) - get all orders for a given MM

### LH

1. getQuote(amountIn, tokenIn, tokenOut) GET - get a quote for a given amountIn of tokenIn for tokenOut
2. confirmSwap(orderId) POST - order can no longer be cancelled
