import "@nomicfoundation/hardhat-toolbox";
import "@openzeppelin/hardhat-upgrades";
import "@primitivefi/hardhat-dodoc";
import * as dotenv from "dotenv";
import { HardhatUserConfig } from "hardhat/config";

dotenv.config();

const config: HardhatUserConfig = {
  solidity: {
    compilers: [
      {
        version: "0.8.19",
        settings: {
          optimizer: {
            enabled: true,
            runs: 115,
          },
          outputSelection: {
            "*": {
              "*": ["storageLayout"],
            },
          },
        },
      },
    ],
  },
  networks: {
    'othentic-blade': {
      url: 'http://localhost/api/eth-rpc'
    },
    root: {
      url: process.env.ROOT_RPC || "",
      accounts: process.env.PRIVATE_KEY !== undefined ? [process.env.PRIVATE_KEY] : [],
    },
    rootTest: {
      url: process.env.ROOT_TEST_RPC || "",
      accounts: process.env.PRIVATE_KEY !== undefined ? [process.env.PRIVATE_KEY] : [],
    },
    child: {
      url: process.env.CHILD_RPC || "",
      accounts: process.env.PRIVATE_KEY !== undefined ? [process.env.PRIVATE_KEY] : [],
    },
    childTest: {
      url: process.env.CHILD_TEST_RPC || "",
      accounts: process.env.PRIVATE_KEY !== undefined ? [process.env.PRIVATE_KEY] : [],
    },
    hardhat: {
      // allow impersonation of smart contracts without modifying balance
      gasPrice: 0,
      hardfork: "berlin",
    },
  },
  gasReporter: {
    enabled: (process.env.REPORT_GAS as unknown as boolean) || false,
    currency: "USD",
  },
  etherscan: {
    apiKey: {
      'othentic-blade': 'empty'
    },
    customChains: [
      {
        network: "othentic-blade",
        chainId: 51001,
        urls: {
          apiURL: "http://localhost/api",
          browserURL: "http://localhost"
        }
      }
    ]
  },
  mocha: {
    timeout: 100000000,
  },
  dodoc: {
    // uncomment to stop docs from autogenerating each compile
    // runOnCompile: false,
    exclude: ["mocks", "openzeppelin/contracts", "openzeppelin/contracts-upgradeable"],
  },
};

export default config;