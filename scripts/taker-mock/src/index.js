import Orderbook from './order-book.js';
import beginAuctionTest from './begin-auction.js';
import makeMarket from './make-market.js';
import * as dotenv from 'dotenv';

console.log('------------------- taker-mock started')
dotenv.config()
console.log('ORDERBOOK_HOST', process.env.ORDERBOOK_HOST);

async function main() {
    const ob = new Orderbook()
    if (! await makeMarket(ob, 2000)) {
        console.log('FAILED to make the market')
        process.exit(1)
    }
    console.log('SUCCESS make the market')
    if (!beginAuctionTest(ob)) {
        console.log("beginAuctionTest failed")
        process.exit(1)
    }

    // if (!abortAuctionTest(ob)) {
    //     console.log("abortAuctionTest failed")
    //     process.exit(1)
    // }
}

main()