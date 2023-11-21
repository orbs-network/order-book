
function delay(time) {
    return new Promise(resolve => setTimeout(resolve, time));
}


async function beginAuction(ob) {
    // BUY
    let depth = await ob.marketDepth()
    console.log(depth)
    let sumBToken = 0
    let sumAToken = 0
    for (const ask of depth.asks) {
        const price = parseFloat(ask[0])
        const size = parseFloat(ask[1])
        sumAToken += size
        sumBToken += size * price

    }
    let res = await ob.beginAuction("ETH-USD", "BUY", sumBToken)
    if (!res) {
        console.error('beginAuction BUY failed')
    }
    else if (parseFloat(res.amountOut) === sumAToken) {
        console.log('SUCCESS BUY amount-out IS equal to sum of size', res.amountOut, sumAToken)
    }
    else {
        console.error('BUY amount out is NOT equal not to sum of size', res.amountOut, sumAToken)
        return false;
    }

    // SELL (reset)
    depth = await ob.marketDepth()
    sumAToken = 0
    sumBToken = 0
    for (const bid of depth.bids) {
        const price = parseFloat(bid[0])
        const size = parseFloat(bid[1])
        sumAToken += size
        sumBToken += size * price

    }
    res = await ob.beginAuction("ETH-USD", "SELL", sumAToken)
    if (!res) {
        console.error('beginAuction SELL failed')
    }
    else if (parseFloat(res.amountOut) === sumBToken) {
        console.log('SUCCESS SELL amount-out IS equal to sum of size', res.amountOut, sumBToken)
        return true;
    }
    else {
        console.error('SELL amount out is NOT equal not to sum of size', res.amountOut, sumBToken)
        return false;
    }

    return true
}
export default beginAuction