async function makeMarket(ob, price) {
    // cancell all orders
    if (!await ob.cancelOrders()) {
        return false;
    }

    const depth = 3
    // SELL
    let indx = 0;
    for (let i = 0; i < depth; ++i) {
        let fact = i + 1;
        indx++;
        if (!await ob.createOrder('ETH-USD', "sell", price + fact, 10 * fact, `${indx}0000000-0000-0000-0000-00000000000${indx}`))
            return false;
    }
    // BUY
    for (let i = 0; i < depth; ++i) {
        let fact = i + 1;
        indx++;
        if (!await ob.createOrder('ETH-USD', "buy", price + fact, 10 * fact, `${indx}0000000-0000-0000-0000-00000000000${indx}`))
            return false;
    }
    return true
}

export default makeMarket