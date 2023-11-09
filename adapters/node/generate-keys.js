const crypto = require('crypto');
const path = require('path');
const fs = require('fs');

// Generate a secp256k1 key pair
const keyPair = crypto.generateKeyPairSync('ec', {
  namedCurve: 'secp256k1'  // Name of the curve
});

// Export the private key as a PEM-formatted string
const privateKey = keyPair.privateKey.export({
  type: 'sec1',
  format: 'pem',
});

// Export the public key as a PEM-formatted string
const publicKey = keyPair.publicKey.export({
  type: 'spki',
  format: 'pem',
});

// Specify the paths where the keys will be saved
const privateKeyPath = path.resolve(__dirname, 'privateKey.pem')
const publicKeyPath = path.resolve(__dirname, 'publicKey.pem')

// Write the keys to disk
fs.writeFileSync(privateKeyPath, privateKey);
fs.writeFileSync(publicKeyPath, publicKey);