//
// to run:
//   example1_run.sh  (to compile all stubs, and java code)
//   example1_server.sh &
//   node example1_node_client.js
//
var foolib = require('./node/foolib-node');

var svc = foolib.SampleServiceClient("http://localhost:9009");

var errHandler = function(res) {
    console.log("ERR: code=" + res.code + " msg=" + res.message);
};

var doAdd = function(callback) {
    svc.Add(3, 31, callback, errHandler);
};

var start = new Date().getTime();
var count = 10000;
var callback = function(res) {
    count--;
    if (count > 0) {
        doAdd(callback);
    }
    else {
        var elapsed = new Date().getTime() - start;
        console.log("Elapsed: " + elapsed);
    }
};

// uncomment to run client benchmark
doAdd(callback);

svc.Add(2, 15, function(res) {
    console.log("Add Result: " + res);
}, errHandler);

var params = { "name" : "Jane" };
svc.Create(params, function(res) {
    console.log("Create Result: " + JSON.stringify(res));
}, errHandler);

svc.StoreName("bob dobbs", function(res) {
    console.log("Store Name: " + res);
}, errHandler);

svc.Say_Hi(function (res) {
    console.log("Say_Hi: " + res);
}, errHandler);

