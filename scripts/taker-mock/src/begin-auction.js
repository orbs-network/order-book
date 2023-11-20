
async function beginAuction(ob) {
    const depth = await ob.marketDepth()
    console.log(depth)
    //for (const ask of depth.asks) {
    const ask = depth.asks[0]
    const price = ask[0]
    const size = ask[1]
    const amountInBToken = price * size
    const expectedAmountOut = size
    const res = await ob.beginAuction("ETH-USD", "BUY", amountInBToken)
    console.log(res)
    //}

    return true
}
export default beginAuction