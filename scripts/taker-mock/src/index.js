import Orderbook from './order-book.js';
import beginAuctionTest from './begin-auction.js';
import abortAuctionTest from './abort-auction.js';
import * as dotenv from 'dotenv';

console.log('------------------- taker-mock started')
dotenv.config()
console.log('ORDERBOOK_HOST', process.env.ORDERBOOK_HOST);
console.log('PUB_KEY', process.env.PUB_KEY);

async function main() {
    const ob = new Orderbook()
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