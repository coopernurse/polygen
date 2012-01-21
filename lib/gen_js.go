package polygenlib

import (
	"fmt"
)

func StartJsFile(p *Package) *StrBuf {
	b := NewStrBuf("//")
	b.prelude()
	return b
}

func JsFilename(pkg string) string {
	return pkg + ".js"
}

type JsGenerator struct { }

func (g JsGenerator) GenFiles(p *Package) []File {
	b := StartJsFile(p)

	b.f("var %s = {", p.Name)
	b.w(jsPostBoilerplate)
	b.w(jsBoilerplate)

	for i := 0; i < len(p.Interfaces); i++ {
		iface := p.Interfaces[i]
		b.blank()
		b.f(" ,  %s : function(_url) {", iface.Name)
		GenJsClientFunc(iface, b, p.Name)
		b.w("    }")
	}
	b.w("};")
	file := File{JsFilename(p.Name), b.b.Bytes()}
	return []File{ file }
}

type NodeJsGenerator struct { }

func (g NodeJsGenerator) GenFiles(p *Package) []File {
	b := StartJsFile(p)
	b.w("var _http   = require('http');")
	b.w("var _https  = require('https');")
	b.w("var _urlmod = require('url');")
	b.blank()
	b.w("var _util = {")
	b.w(nodePostBoilerplate)
	b.w(nodeBoilerplate)
	b.w("};")
	b.blank()
	b.w(nodeReadRequestBoilerplate)

	for i := 0; i < len(p.Interfaces); i++ {
		iface := p.Interfaces[i]

		b.blank()
		b.f("exports.Dispatch%s = function(rpcreq, svc, onSuccess, onError) {", iface.Name)
		b.w("    var method = rpcreq.method;")
		b.w("    var params = rpcreq.params;")
		for x := 0; x < len(iface.Methods); x++ {
			m := iface.Methods[x]
			if x > 0 {
				b.raw("    else ")
			} else {
				b.raw("    ")
			}

			b.f("if (method === \"%s_%s\") {", iface.Name, m.Name)

			if len(m.Args) == 0 {
				b.f("        svc.%s(onSuccess, onError);", m.Name)
			} else if len(m.Args) == 1 {
				b.f("        svc.%s(params, onSuccess, onError);", m.Name)
			} else {
				plist := ""
				for y := 0; y < len(m.Args); y++ {
					plist += fmt.Sprintf("params[%d], ", y)
				}
				b.f("        svc.%s(%sonSuccess, onError);", m.Name, plist)
			}
			b.w("    }")
		}
		b.w("    else {")
		b.w("        onError(-32601, \"Method not found: \" + method);")
		b.w("    }")
		b.w("};")

		b.blank()
		b.f("exports.%sHttpServer = function(svc, maxPostLen) {", iface.Name)
		b.w("    return _util.createServer(_http, exports.DispatchSampleService, svc, maxPostLen);")
		b.w("};")

		b.blank()
		b.f("exports.%sHttpsServer = function(svc, maxPostLen) {", iface.Name)
		b.w("    return _util.createServer(_https, exports.DispatchSampleService, svc, maxPostLen);")
		b.w("};")

		b.blank()
		b.f("exports.%sClient = function(_url) {", iface.Name)
		GenJsClientFunc(iface, b, "_util")
		b.w("};")
	}

	file := File{JsFilename(p.Name + "-node"), b.b.Bytes()}
	return []File{ file }
}

func GenJsClientFunc(iface Interface, b *StrBuf, utilname string) {
	b.w("        var _me = {};")
    b.w("        var _tmp = _urlmod.parse(_url);")
    b.w("        _url = { 'host': _tmp.hostname, 'port': _tmp.port, 'path': _tmp.pathname, 'protocol': _tmp.protocol };")
	for x := 0; x < len(iface.Methods); x++ {
		m := iface.Methods[x]
		if len(m.Args) == 0 {
			b.f("        _me.%s = function(_onSuccess, _onError) {", m.Name)
			b.w("            var _args = null;")
		} else if len(m.Args) == 1 {
			b.f("        _me.%s = function(%s, _onSuccess, _onError) {", m.Name, m.Args[0].Name)
			b.f("            var _args = %s;", m.Args[0].Name)
		} else {
			args := ""
			for y := 0; y < len(m.Args); y++ {
				if y > 0 {
					args += ", "
				}
				args += m.Args[y].Name
			}
			b.f("        _me.%s = function(%s, _onSuccess, _onError) {", m.Name, args)
			b.f("            var _args = [ %s ];", args)
		}
		b.f("            %s.rpcCall(_url, \"%s_%s\", _args, _onSuccess, _onError);", utilname, iface.Name, m.Name)
		b.w("        };")
	}
	b.w("        return _me;")
}

var jsPostBoilerplate = `    post : function(url, obj, callback) {
        var json = JSON.stringify(obj);
        jQuery.ajax({ type: 'POST', 
                      url: url,
                      dataType: 'json',
                      data: json,
                      success: callback,
                      error: callback});
    },`

var jsBoilerplate = `    S4 : function() {
        return (((1+Math.random())*0x10000)|0).toString(16).substring(1);
    },

    uuid : function() {
        return (this.S4()+this.S4()+"-"+this.S4()+"-"+this.S4()+"-"+this.S4()+"-"+this.S4()+this.S4()+this.S4());
    },

    rpcCall : function(url, method, params, onSuccess, onError) {
        var obj = { "jsonrpc": "2.0", "id": this.uuid(), "method": method };
        if (params) {
            obj.params = params;
        }
        this.post(url, obj, function(rpcResp) {
            if (rpcResp && rpcResp.error) {
                onError(rpcResp.error);
            }
            else if (rpcResp && rpcResp.result) {
                onSuccess(rpcResp.result);
            }
            else {
                onError({ "code" : -33000, 
                          "message" : "Invalid response: " + rpcResp });
            }
        });
    }`

var nodePostBoilerplate = `    post : function(urlInfo, obj, callback) {
        var json = JSON.stringify(obj);
        var settings = { 
            host: urlInfo.host, 
            path: urlInfo.path, 
            method: 'POST'};
        settings.headers = { 
            'Content-Type' : 'application/json',
            'Content-Length' : json.length
        };
        var req;
        if (urlInfo.protocol === 'https:') {
            settings.port = urlInfo.port || 443;
            req = _https.request(settings);
        } else {
            settings.port = urlInfo.port || 80;
            req = _http.request(settings);
        }
        req.write(json);
        req.on('response', function(res) {
            res.body = '';
            res.setEncoding('utf-8');

            // concat chunks
            res.on('data', function(chunk) { res.body += chunk; });
        
            // when the response has finished
            res.on('end', function(){
            
                // fire callback
                callback(JSON.parse(res.body));
            });
        });
        req.end();
    },`

var nodeReadRequestBoilerplate = `exports.ReadServerRequest = function(req, maxlen, onSuccess, onError) {
    var body = '';
    req.on('data', function(chunk) { 
        body += chunk; 
        if (body.length > maxlen) {
            onError(-32000, "Request exceeded max length: " + maxlen);
        }
    });
    req.on('end', function() {
        onSuccess(body);
    });
};`

var nodeBoilerplate = jsBoilerplate + `,

    createServer : function(httpMod, dispatcher, svc, maxPostLength) {
        var sendResp = function(res, jsonresp) {
            var respdata = JSON.stringify(jsonresp);
            res.writeHead(200, { "Content-Length": respdata.length });
            res.end(respdata);
        };

        return httpMod.createServer(function (req, res) {
            var onPostData = function(data) {
                var jsonreq = JSON.parse(data);
                var jsonresp = { "jsonrpc": "2.0", "id" : jsonreq.id };
            
                var onSuccess = function(data) {
                    if ((data === null) || (data === undefined)) {
                        data = true;
                    }
                    jsonresp.result = data;
                    sendResp(res, jsonresp);
                };
            
                var onError = function(code, message) {
                    jsonresp.error = { "code" : code, "message" : message };
                    sendResp(res, jsonresp);
                };
                
                dispatcher(jsonreq, svc, onSuccess, onError);
            };
            
            var onPostErr = function(code, message) {
                var error = { "code" : code, "message" : message };
                sendResp(res, { "jsonrpc": "2.0", "error" : error });
            };
            
            exports.ReadServerRequest(req, maxPostLength, onPostData, onPostErr);
        });
    }`