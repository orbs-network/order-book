import Orderbook from './orderBook.js';

async function main() {
    const ob = new Orderbook()
    const depth = await ob.marketDepth()
    console.log(depth)
}

main()