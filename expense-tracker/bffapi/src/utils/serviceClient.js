const axios = require('axios');
const dotenv = require('dotenv');

dotenv.config();

// Create axios instances for each backend service with base URLs from environment variables
// sample nodeJS code snippet
const accountsServiceurl = process.env.CHOREO_ACCOUNTS_SERVICEURL;
const accountsChoreoapikey = process.env.CHOREO_ACCOUNTS_CHOREOAPIKEY;

const accountsClient = axios.create({
  baseURL: accountsServiceurl,
  headers: {
    'Content-Type': 'application/json',
    'Choreo-API-Key': `${accountsChoreoapikey}`
  }
});

// sample nodeJS code snippet
const receiptsServiceurl = process.env.CHOREO_RECEIPTS_SERVICEURL;
const receiptsChoreoapikey = process.env.CHOREO_RECEIPTS_CHOREOAPIKEY;

const billParserClient = axios.create({
  baseURL: receiptsServiceurl,
  headers: {
    'Content-Type': 'application/json',
    'Choreo-API-Key': `${receiptsChoreoapikey}`
  }
});

// Add response interceptor to log errors
const addErrorInterceptor = (client) => {
  client.interceptors.response.use(
    response => response,
    error => {
      console.error('API Error:', error.message);
      if (error.response) {
        console.error('Response data:', error.response.data);
        console.error('Response status:', error.response.status);
      }
      return Promise.reject(error);
    }
  );
};

addErrorInterceptor(accountsClient);
addErrorInterceptor(billParserClient);

module.exports = { accountsClient, billParserClient };
