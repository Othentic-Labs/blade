const { ethers } = require('ethers');
const fs = require('fs');
const path = require('path');

// Get command-line arguments
const args = process.argv.slice(2);
console.log('Received arguments:', args);

if (args.length < 6) {
  console.error('Usage: node deploy_erc20.js <providerURL> <privateKey> <tokenName> <tokenSymbol> <decimals> <totalSupply>');
  process.exit(1);
}

// Extract arguments
const providerURL = args[0];
const privateKey = args[1];
const tokenName = args[2];
const tokenSymbol = args[3];
const decimals = parseInt(args[4]);
const totalSupply = args[5];

console.log('providerURL:', providerURL);
console.log('privateKey:', privateKey);
console.log('tokenName:', tokenName);
console.log('tokenSymbol:', tokenSymbol);
console.log('decimals:', decimals);
console.log('totalSupply:', totalSupply);

// Read ABI and bytecode from files
const abiPath = path.join(__dirname, 'erc20.abi');
const binPath = path.join(__dirname, 'erc20.bin');

const ERC20_ABI = JSON.parse(fs.readFileSync(abiPath, 'utf8'));
const ERC20_BYTECODE = fs.readFileSync(binPath, 'utf8').trim();

const provider = new ethers.JsonRpcProvider(providerURL);
console.log('Connected to provider:', providerURL);

const wallet = new ethers.Wallet(privateKey, provider);
console.log('Wallet initialized with address:', wallet.address);

(async () => {
  try {
    console.log('Deploying contract with the following parameters:');
    console.log('Token Name:', tokenName);
    console.log('Token Symbol:', tokenSymbol);
    console.log('Decimals:', decimals);
    console.log('Total Supply:', totalSupply);

    const factory = new ethers.ContractFactory(ERC20_ABI, ERC20_BYTECODE, wallet);

    console.log('Contract factory created. Deploying the contract...');
    const contract = await factory.deploy(totalSupply, tokenName, decimals, tokenSymbol);

    console.log('Contract deploying...');

    await contract.waitForDeployment();

    console.log('Contract deployed...' + JSON.stringify(contract));


    console.log('Contract deployed at address:', contract.target);

    // Save the contract address to a file or output it
    fs.writeFileSync('/data/erc20_address.txt', contract.target);
    console.log('Contract address saved to /data/erc20_address.txt');
  } catch (error) {
    console.error('Error deploying ERC20 contract:', error);
  }
})();
