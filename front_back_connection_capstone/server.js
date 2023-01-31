const WebSocket = require('ws')

//create a new websocket server
const wss = new WebSocket.Server({port:8000})

//varioous events inside websockets
wss.on('connection', (ws)=> {
    console.log('a new client is connected')

    //listen for more events

    
//senindg objects as string sto output to client
    const user = {
        name: 'Jose',
        age: 25,
        country: 'USA'
    }
    ws.send(JSON.stringify(user))

// send data to clients
    ws.send('hi this is the message from the server')

    ws.on('message', (data) => {
        console.log(data.toString())
    })
})
