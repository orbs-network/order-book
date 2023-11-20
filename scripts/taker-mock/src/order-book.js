import fetch from "node-fetch";

class Orderbook {
    constructor() {
        this.ORDERBOOK_HOST = process.env.ORDERBOOK_HOST || "http://localhost:8080/"
        const PUB_KEY = process.env.xx || "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEhqhj8rWPzkghzOZTUCOo/sdkE53sU1coVhaYskKGKrgiUF7lsSmxy46i3j8w7E7KMTfYBpCGAFYiWWARa0KQwg=="
        this.headers = {
            "X-Public-Key": PUB_KEY,
            "Content-Type": "application/json", // Assuming you are expecting JSON in response
        };
    }
    async apiCall(method, path, body) {
        const url = `${this.ORDERBOOK_HOST}${path}`
        try {
            const req = {
                method: method,
                headers: this.headers
            }
            // add body
            if (body) {
                req.body = JSON.stringify(body)
            }
            const response = await fetch(url, req);

            // Parse and work with the JSON response
            return await response.json();
        }
        catch (error) {
            // Handle any errors that occurred during the fetch
            console.error("Fetch error:", error);
            return null;
        }
    }
    async marketDepth() {
        const jsonResponse = await this.apiCall('GET', '/api/v1/orderbook/ETH-USD?limit=20')
        if (jsonResponse.code !== 'OK') {
            //console.log("Response JSON:", jsonResponse.data);
            console.error('failed to get market depth')
            return null
        }
        return jsonResponse.data
    }
    async beginAuction(symbol, side, size) {
        const auctionId = '10000000-0000-0000-0000-000000000001';
        const body = {
            "amountIn": String(size),
            "symbol": symbol,
            "side": side
        }
        return await this.apiCall('POST', `/lh/v1/begin_auction/${auctionId}`, body)
    }

}

export default Orderbook
