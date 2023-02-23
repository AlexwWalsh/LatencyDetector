const express = require('express');
const http = require('http');

const app = express();

//Confirm cors settings for correct header info
const cors=require("cors");
const corsOptions ={
   origin:'*', 
   credentials:true,            //access-control-allow-credentials:true
   optionSuccessStatus:200,
}

app.use(cors(corsOptions)) // Use this after the variable declaration

app.get('/', (req, res) => {
  res.send('Welcome to the Node.js REST server!');
});

//Make server listen on port 3000
const server = app.listen(3000, () => {
  console.log('Node.js REST server is running on port 3000');
});

//Endpoint that forwards the request to Gorilla Mux server on 8080, returns result back to client
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
  }
);