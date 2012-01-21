
exports.SampleServiceServer = function() {
    var _me = {};

    _me.Create = function(p, _onSuccess, _onError) {
        console.log("Create p=" + JSON.stringify(p));
        _onSuccess({"success":true, "code":999, "note":"node.js here"});
    };

    _me.Add = function(a, b, _onSuccess, _onError) {
        console.log("Add a=" + a + " b=" + b);
        _onSuccess(a+b);
    };

    _me.StoreName = function(name, _onSuccess, _onError) {
        console.log("StoreName name=" + name);
        _onSuccess();
    };

    _me.Say_Hi = function(_onSuccess, _onError) {
        console.log("Say_Hi");
        _onSuccess("hello from node.js");
    };

    _me.getPeople = function(params, _onSuccess, _onError) {
        console.log("getPeople params=" + JSON.stringify(params));
        var people = [
            { "id": 10, "name": "bob dobbs", "email": "bob@example.com" },
            { "id": 20, "name": "edison carter", "email": "ed@example.com" }
        ];
        _onSuccess(people);
    };
    
    return _me;
};

if (require.main === module) {
    var port = 9009;
    var host = "localhost";
    console.log("Starting node server on " + host + ":" + port);

    var foolib = require("./example1/node/foolib-node.js");
    var svc = exports.SampleServiceServer();

    foolib.SampleServiceHttpServer(svc, 1024*1024).listen(port, host);
}