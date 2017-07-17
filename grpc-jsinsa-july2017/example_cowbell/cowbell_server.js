var PROTO_PATH = __dirname + '/../service/service.proto';

var grpc = require('grpc');
var cowbell = grpc.load(PROTO_PATH).cowbell;

var total = 0;

/**
 * Implements the MoreCowbell RPC method.
 */
function moreCowbell(call, callback) {
    var cowbellResponse;

    console.log("moreCowbell qty="+ call.request.qty);

    total += call.request.qty; 

    cowbellResponse = {
        total: total
    };

    callback(null, cowbellResponse);
}

/**
 * Starts an insecure RPC server that receives requests for the Cowbell on port 9090
 */
function main() {
  var server = new grpc.Server();
  server.addService(cowbell.CowbellService.service, {
    moreCowbell: moreCowbell
  });
  server.bind('0.0.0.0:9090', grpc.ServerCredentials.createInsecure());
  server.start();

  console.log("gRPC Server started on 0.0.0.0:9090")
}

main();