import fetch from "node-fetch";

class Orderbook {
    constructor() {
        this.ORDERBOOK_HOST = process.env.ORDERBOOK_HOST || "http://localhost:8080/"
        const API_KEY = process.env.xx || "abcdef12345"
        this.headers = {
            "X-API-Key": `Bearer ${API_KEY}`,
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

            // error
            if (response.status >= 400) {
                const message = await response.text()
                console.error(message)
                return null
            }

            // Check if the response has a JSON body
            const contentType = response.headers.get('Content-Type');
            if (contentType && contentType.includes('application/json')) {
                // Parse and work with the JSON response
                return await response.json();
            } else {
                // Response does not have a JSON body, it might be empty or in a different format
                return response.status;
            }
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
        const swapId = '10000000-0000-0000-0000-000000000001';
        const body = {
            "amountIn": String(size),
            "symbol": symbol,
            "side": side
        }
        return await this.apiCall('POST', `/lh/v1/begin_auction/${swapId}`, body)
    }
    async createOrder(symbol, side, price, size, cOId) {
        const body = {
            price: String(price),
            size: String(size),
            side: side,
            symbol: symbol,
            ClientOrderId: cOId,
        }
        const jsonResponse = await this.apiCall('POST', '/api/v1/order', body)
        return jsonResponse !== null && jsonResponse.orderId.length > 0;
    }
    async cancelOrders() {
        return await this.apiCall('DELETE', `/api/v1/orders`)
    }

}

export default Orderbook
