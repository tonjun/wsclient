const WebSocket = require('ws');

const wss = new WebSocket.Server({ port: 8080 });

wss.on('connection', function connection(ws) {
  console.log("ws connection");

  ws.on('message', function incoming(message) {
    console.log("received: %s", message)

    var obj = JSON.parse(message)
    if (!obj) {
      console.log("unable to parse JSON message")
      return
    }
    console.log("json obj:", obj);
    if (obj.op === 'get-time') {
      console.log("time: %s", new Date());
      ws.send(JSON.stringify({
        'op': 'get-time-response',
        //'data': new Date(),
      }));
    }

  });

});
