import Orderbook from './order-book.js';
import beginAuction from './begin-auction.js';
import * as dotenv from 'dotenv';

console.log('------------------- taker-mock started')
dotenv.config()
console.log('ORDERBOOK_HOST', process.env.ORDERBOOK_HOST);
console.log('PUB_KEY', process.env.PUB_KEY);

async function main() {
    const ob = new Orderbook()
    if (!beginAuction(ob)) {
        return console.log("beginAuction failed")
    }
}

main()