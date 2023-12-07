const crypto = require('crypto');
const fs = require('fs');
const path = require('path');

// Generate a random private key
let privKey;
do {
  privKey = crypto.randomBytes(32);
} while (crypto.createECDH('secp256k1').setPrivateKey(privKey, 'hex').getPublicKey() === null);

// Get the ECDH object and set the private key
const ecdh = crypto.createECDH('secp256k1');
ecdh.setPrivateKey(privKey);

// Get the public key in uncompressed format
const pubKey = ecdh.getPublicKey('hex', 'uncompressed');

const privateKeyPath = path.join(__dirname, 'privateKey.txt');
const publicKeyPath = path.join(__dirname, 'publicKey.txt');

// Write private key to disk
fs.writeFileSync(privateKeyPath, privKey.toString('hex'));

// Write public key to disk
fs.writeFileSync(publicKeyPath, pubKey);

console.log(`Private key saved to ${privateKeyPath}`);
console.log(`Public key saved to ${publicKeyPath}`);