const path = require('path');
const crypto = require('crypto');
const fs = require('fs');

const requestPayload = {
  symbol: "BTCUSD",
  orderType: "limit",
  side: "buy",
  quantity: "1.0",
  price: "50000.00",
  timestamp: "1697813554"
};

const signature = "30440220306c22dc5ab8c650d7bb59934f140d33d089e2ff05bb5de889ccf531a5591dff0220673cd5d2bc2ce5c815917cc680cd2076a8b560d4ab4305531a84791ea937ade7"

// Load the public key from disk
const publicKeyPath = path.resolve(__dirname, 'publicKey.pem');
const publicKey = fs.readFileSync(publicKeyPath, 'utf8');
console.log('publicKey :', publicKey);

const jsonString = JSON.stringify(requestPayload);

// Create a hash of the JSON string
const hash = crypto.createHash('sha256').update(jsonString).digest();

// Create a verify object
const verify = crypto.createVerify('SHA256');

// Update the verify object with the hash
verify.update(hash);

// Verify the signature
const isSignatureValid = verify.verify(publicKey, signature, 'hex');

console.log('Is signature valid?', isSignatureValid);