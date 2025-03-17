const axios = require('axios');
const dotenv = require('dotenv');

dotenv.config();

// Create axios instances for each backend service with base URLs from environment variables
const accountServiceurl = process.env.CHOREO_ACCOUNTS_SERVICEURL;
const accountChoreoapikey = process.env.CHOREO_ACCOUNTS_CHOREOAPIKEY;
const accountsClient = axios.create({
  baseURL: accountServiceurl,
  headers: {
    'Content-Type': 'application/json',
    'Choreo-API-Key': `${accountChoreoapikey}`
  }
});

// sample nodeJS code snippet
const billServiceurl = process.env.CHOREO_BILL_PARSER_SERVICEURL;
const billChoreoapikey = process.env.CHOREO_BILL_PARSER_CHOREOAPIKEY;
const billParserClient = axios.create({
  baseURL: billServiceurl,
  headers: {
    'Content-Type': 'application/json',
    'Choreo-API-Key': `${billChoreoapikey}`
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
