import fetch from "node-fetch";

const ORDERBOOK_HOST = process.env.ORDERBOOK_HOST || "http://localhost:8080/"
const API_KEY = process.env.API_KEY || "abcdef12345"
const url = `${ORDERBOOK_HOST}/api/v1/orderbook/ETH-USD?limit=20`


const headers = {
    "X-API-Key": `Bearer ${API_KEY}`,
    "Content-Type": "application/json", // Assuming you are expecting JSON in response
};

// Make the GET request using the node-fetch library
fetch(url, {
    method: 'GET',
    headers: headers,
})
    .then(response => {
        // Check if the request was successful (status code 200)
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        // Parse and work with the JSON response
        return response.json();
    })
    .then(jsonResponse => {
        console.log("Response JSON:", jsonResponse);
        // Handle the JSON response here
    })
    .catch(error => {
        // Handle any errors that occurred during the fetch
        console.error("Fetch error:", error);
    });