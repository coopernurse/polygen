package polygenlib

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
	b.w(jsBoilerplate)

	for i := 0; i < len(p.Interfaces); i++ {
		iface := p.Interfaces[i]
		b.blank()
		b.f("    %s : function(_uri) {", iface.Name)
		b.w("        var _me = {};")
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
			b.f("            %s.rpcCall(_uri, \"%s_%s\", _args, _onSuccess, _onError);", p.Name, iface.Name, m.Name)
			b.w("        };")
		}
		b.w("        return _me;")
		if (i+1) < len(p.Interfaces) {
			b.w("    },")
		} else {
			b.w("    }")
		}
	}
	b.w("};")
	file := File{JsFilename(p.Name), b.b.Bytes()}
	return []File{ file }
}


var jsBoilerplate = `    S4 : function() {
        return (((1+Math.random())*0x10000)|0).toString(16).substring(1);
    },

    uuid : function() {
        return (this.S4()+this.S4()+"-"+this.S4()+"-"+this.S4()+"-"+this.S4()+"-"+this.S4()+this.S4()+this.S4());
    },

    post : function(uri, obj, callback) {
        var json = JSON.stringify(obj);
        jQuery.ajax({ type: 'POST', 
                      url: uri,
                      dataType: 'json',
                      data: json,
                      success: callback,
                      error: callback});
    },

    rpcCall : function(uri, method, params, onSuccess, onError) {
        var obj = { "jsonrpc": "2.0", "id": this.uuid(), "method": method };
        if (params) {
            obj.params = params;
        }
        this.post(uri, obj, function(rpcResp) {
            if (rpcResp.error) {
                onError(rpcResp.error);
            }
            else {
                onSuccess(rpcResp.result);
            }
        });
    },`