async function test(ob, side, extra) {
    console.log(`--- behinAuction ${side} ${extra ? 'expect insufficient' : ''}`)
    const depth = await ob.marketDepth()
    const buySide = side === "BUY"
    let orders = buySide ? depth.asks : depth.bids
    let sumBToken = 0
    let sumAToken = 0
    for (const ask of orders) {
        const price = parseFloat(ask[0])
        const size = parseFloat(ask[1])
        sumAToken += size
        sumBToken += size * price

    }
    let sumToken = buySide ? sumBToken : sumAToken
    const sumTokenOposite = buySide ? sumAToken : sumBToken

    if (extra) {
        sumToken += 1;
    }
    let res = await ob.beginAuction("ETH-USD", side, sumToken)
    // expected error
    if (extra && !res) {
        console.log(`SUCCESS ${side}\texpected null res`)
        return true
    }
    else if (!res) {
        console.error(`beginAuction ${side} failed`)
    }
    else if (parseFloat(res.amountOut) === sumTokenOposite) {
        console.log(`SUCCESS ${side}\tamount-out IS equal to sum of size `, res.amountOut, sumTokenOposite)
    }
    else {
        console.error(`FAIL ${side} amount out is NOT  equal not to sum of size`, res.amountOut, sumTokenOposite)
        return false;
    }
    return true
}
async function beginAuctionTest(ob) {
    if (!await test(ob, "BUY"))
        return false;

    if (!await test(ob, "SELL"))
        return false;

    // insufficiant liquidity
    if (!await test(ob, "BUY", true))
        return false;

    // insufficiant liquidity
    if (!await test(ob, "SELL", true))
        return false;

    return true
}
export default beginAuctionTest