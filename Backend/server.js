const express = require('express');
const http = require('http');
const WebSocket = require('ws');
const bodyParser = require('body-parser');

const app = express();
const wss = new WebSocket.Server({ port: 5000 });

let ipAddresses = []; // your list of IP addresses
const cors=require("cors");
const corsOptions ={
   origin:'*', 
   credentials:true,            //access-control-allow-credentials:true
   optionSuccessStatus:200,
}

app.use(cors(corsOptions)) // Use this after the variable declaration

app.use(bodyParser.json())

app.get('/', (req, res) => {
  res.send('Welcome to the Node.js REST server!');
});

//Make server listen on port 3000
//Change the below IP address to your machine's IPv4 IP Address OR Localhost if you do not need VM connection
//Example alternative:
const server = app.listen(3000, '192.168.1.64', () => {
  console.log('Node.js REST server is running on port 3000');
});
// const server = app.listen(3000, 'localhost', () => {
//  console.log('Node.js REST server is running on port 3000');
// });


app.post('/server', (req, res) => {
  const newData = req.body
  // console.log("req.body", req.body)
  ipAddresses.push(newData);
  console.log("IP addresses: ", ipAddresses)
  wss.clients.forEach(client => {
    client.send(JSON.stringify(ipAddresses))
  });
  res.send("Node server recieved IP Address")
})


app.get('/grabMuxInfo/:IpAddress/:Endpoint', (req, res) => {
  const { IpAddress } = req.params;
  console.log(req.params.Endpoint)
  // console.log(IpAddress)
  const options = {
    hostname: IpAddress,
    // hostname: "localhost",
    port: 8080,
    path: '/' + req.params.Endpoint,
    method: 'GET'
  };

  const gorrilaReq = http.request(options, (gorrilaRes) => {
    let data = '';
    gorrilaRes.on('data', (chunk) => {
      data += chunk;
    });

    gorrilaRes.on('end', () => {
      res.send(data)
    });
  });

  gorrilaReq.on('error', (error) => {
    console.error(error);
  });

  gorrilaReq.end();
});

wss.on('connection', (ws) => {
  ws.send(JSON.stringify(ipAddresses));
});
