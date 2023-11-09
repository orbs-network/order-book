const path = require('path');
const crypto = require('crypto');
const fs = require('fs');

// Load the private key from disk
const privateKeyPath = path.resolve(__dirname, 'privateKey.pem')
const privateKey = fs.readFileSync(privateKeyPath, 'utf8');

const requestPayload = {
  symbol: "BTCUSD",
  orderType: "limit",
  side: "buy",
  quantity: "1.0",
  price: "50000.00",
  timestamp: "1697813554"
};

const jsonString = JSON.stringify(requestPayload);

// Create a hash of the JSON string
const hash = crypto.createHash('sha256').update(jsonString).digest();

// Create a signing object
const sign = crypto.createSign('SHA256');

// Update the signing object with the hash
sign.update(hash);

// Generate the signature using the private key
const signature = sign.sign(privateKey, 'hex');

console.log('signature :', signature);