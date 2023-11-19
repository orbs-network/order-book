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
    async apiCall(method, path) {
        const url = `${this.ORDERBOOK_HOST}/api/v1/${path}`
        try {
            const response = await fetch(url, {
                method: method,
                headers: this.headers,
            });

            // Check if the request was successful (status code 200)
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
                return null;
            }

            // Parse and work with the JSON response
            const jsonResponse = await response.json();
            console.log("Response JSON:", jsonResponse);
            return jsonResponse
            // Handle the JSON response here
        } catch (error) {
            // Handle any errors that occurred during the fetch
            console.error("Fetch error:", error);
            return null;
        }
    }
    async marketDepth() {
        return await this.apiCall('GET', 'orderbook/ETH-USD?limit=20')
    }

}

export default Orderbook
