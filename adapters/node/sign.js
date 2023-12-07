const fs = require('fs');
const path = require('path');
const {
  signTypedData
} = require("@metamask/eth-sig-util");

// Define the path to the private key file
const privateKeyPath = path.join(__dirname, 'privateKey.txt');

// Read the private key from the file
const privateKeyHex = fs.readFileSync(privateKeyPath, 'utf8');


const order = {
  price: "20.99",
  size: "1000",
  symbol: "BTC-ETH",
  side: "sell",
  clientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37"
};

const domain = {
  name: 'orderbook',
  version: '1.0',
  chainId: 1, // Mainnet ID, change accordingly if you're using a different network
  verifyingContract: '0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC' // Replace with your contract address
};

const types = {
  Order: [
    { name: 'price', type: 'string' },
    { name: 'size', type: 'string' },
    { name: 'symbol', type: 'string' },
    { name: 'side', type: 'string' },
    { name: 'clientOrderId', type: 'string' }
  ]
};

const data = {
  types: {
    EIP712Domain: [
      { name: 'name', type: 'string' },
      { name: 'version', type: 'string' },
      { name: 'chainId', type: 'uint256' },
      { name: 'verifyingContract', type: 'address' }
    ],
    ...types
  },
  primaryType: 'Order',
  domain: domain,
  message: order
};

// Sign the EIP-712 structured data
const signature = signTypedData({
  // Remove the 0x prefix if present
  privateKey: Buffer.from(privateKeyHex.startsWith('0x') ? privateKeyHex.slice(2) : privateKeyHex, "hex"),
  data: data,
  version: "V4",
});

console.log('EIP-712 Signature:', signature);