const { ethers } = require("ethers");

// 3. Set up a provider (e.g., using Infura, Alchemy, or a local node)
const provider = new ethers.JsonRpcProvider("http://localhost:8545");

// 4. Define the ERC-20 token contract address and ABI
const tokenAddress = "0xbbe066ff0f613b86a36a53fca229ca9540f23c3e";
const tokenABI = [
    // Only include the balanceOf function from the ABI
    "function balanceOf(address owner) view returns (uint256)"
];

// 5. Create a contract instance
const tokenContract = new ethers.Contract(tokenAddress, tokenABI, provider);

// 6. Define the address you want to check the balance for
const walletAddress = "0x861F89324C60E5158C7f5b38F84E9b2AB76552f7";

// 7. Fetch the balance
async function getBalance() {
    try {
        const balance = await tokenContract.balanceOf(walletAddress);
        console.log(`Balance of ${walletAddress}: ${ethers.formatUnits(balance, 18)} Tokens`);
    } catch (error) {
        console.error("Error fetching balance:", error);
    }
}

// 8. Call the function
getBalance();