const express = require('express');
const http = require('http');

const app = express();

app.get('/', (req, res) => {
  res.send('Welcome to the Node.js REST server!');
});

const server = app.listen(3000, () => {
  console.log('Node.js REST server is running on port 3000');
});

app.get('/grabMuxInfo', (req, res) => {
  const options = {
    hostname: 'localhost',
    port: 8080,
    path: '/data',
    method: 'GET'
  };

  const gorrilaReq = http.request(options, (gorrilaRes) => {
    let data = '';

    gorrilaRes.on('data', (chunk) => {
      data += chunk;
    });

    gorrilaRes.on('end', () => {
      res.send(data);
    });
  });

  gorrilaReq.on('error', (error) => {
    console.error(error);
  });

  gorrilaReq.end();
});