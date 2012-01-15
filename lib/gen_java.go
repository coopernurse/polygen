package polygenlib

import (
	"strings"
)

type JavaGenerator struct { }

func (g JavaGenerator) GenFiles(p *Package) []File {
	files := make([]File, 0)
	files = append(files, g.genRPCException(p))
	files = append(files, g.genRPCError(p))

	for i := 0; i < len(p.Structs); i++ {
		files = append(files, g.genStructClass(p, p.Structs[i]))
	}

	for i := 0; i < len(p.Interfaces); i++ {
		iface := p.Interfaces[i]
		files = append(files, g.genServiceInterface(p, iface))
		files = append(files, g.genServiceRPCServer(p, iface))
		files = append(files, g.genServiceRPCClient(p, iface))
		files = append(files, g.genServiceTypes(p, iface))
	}

	return files
}

func JavaFilename(s string) string {
	return s + ".java"
}

func JavaType(s string) string {
	switch s {
	case "int":
		return "Long"
	case "float":
		return "Double"
	case "bool":
		return "Boolean"
	case "string":
		return "String"
	case "":
		return "void"
	}

	return s
}

func JavaName(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func VarName(s string) string {
	return strings.ToLower(s)
}

func ServiceResponseType(t string) string {
	return JavaType(t) + "TypeRespObj"
}

func ParamsAsList(m Method) string {
	b := NewStrBuf("//")
	b.raw("java.util.Arrays.asList(")
	for x := 0; x < len(m.Args); x++ {
		if x > 0 { 
			b.raw(",")
		}
		b.fraw("%s", VarName(m.Args[x].Name))
	}
	b.raw(")")
	return b.b.String()
}

func MethodSig(m Method) string {
	ret := JavaType(m.ReturnType)
	b := NewStrBuf("//")
	b.fraw("public %s %s(", ret, m.Name)
	for x := 0; x < len(m.Args); x++ {
		if x > 0 { 
			b.raw(", ")
		}
		b.fraw("%s %s", JavaType(m.Args[x].Type), VarName(m.Args[x].Name))
	}
	b.raw(") throws RPCException")
	return b.b.String()
}

func StartFile(p *Package) *StrBuf {
	b := NewStrBuf("//")
	b.prelude()
	b.f("package %s;", p.Name)
	b.blank()
	return b
}

func (g JavaGenerator) genStructClass(p *Package, s Struct) File {
	b := StartFile(p)
	b.f("public class %s {", s.Name)
	for i := 0; i < len(s.Props); i++ {
		vname := VarName(s.Props[i].Name)
		b.f("    private %s %s;", JavaType(s.Props[i].Type), vname)
	}
	b.blank()
	for i := 0; i < len(s.Props); i++ {
		t := JavaType(s.Props[i].Type)
		upper := JavaName(s.Props[i].Name)
		vname := VarName(s.Props[i].Name)
		b.f("    public %s get%s() { return this.%s; }", t, upper, vname)
		b.f("    public void set%s(%s val) { this.%s = val; }", upper, t, vname)
	}
	b.blank()
	b.w("}")
	return File{JavaFilename(s.Name), b.b.Bytes()}
}

func (g JavaGenerator) genServiceInterface(p *Package, iface Interface) File {
	cname := iface.Name
	b := StartFile(p)
	b.f("public interface %s {", cname)
	b.blank()
	for i := 0; i < len(iface.Methods); i++ {
		b.f("    %s;", MethodSig(iface.Methods[i]))
	}
	b.blank()
	b.w("}")
	return File{JavaFilename(cname), b.b.Bytes()}
}

func (g JavaGenerator) genServiceRPCServer(p *Package, iface Interface) File {
	cname := iface.Name + "RPCServer"
	b := StartFile(p)
	b.f("public class %s {", cname)
	b.blank()
	b.w("}")
	return File{JavaFilename(cname), b.b.Bytes()}
}

func (g JavaGenerator) genServiceRPCClient(p *Package, iface Interface) File {
	cname := iface.Name + "RPCClient"
	tclass := iface.Name + "Types"
	b := StartFile(p)
	b.w("import org.codehaus.jackson.map.ObjectMapper;")
	b.blank()
	b.f("public class %s implements %s {", cname, iface.Name)
	b.blank()
	b.w(rpcClientBoilerplate)
	b.f("    public %s(String url) {", cname)
	b.w("        this._prv = new PolygenHttpProvider(url); }")
	b.f("    public %s(PolygenProvider provider) {", cname)
	b.w("        this._prv = provider; }")
	b.blank()
	for i := 0; i < len(iface.Methods); i++ {
		m := iface.Methods[i]
		mname := iface.Name + "_" + m.Name
		b.f("    %s {", MethodSig(m))
		if len(m.Args) == 0 {
			b.f("        %s.BaseReqObj _rq = ", tclass)
            b.f("          new %s.BaseReqObj(\"%s\");", tclass, mname)
		} else if len(m.Args) == 1 {
			b.f("        %s.BaseParamsReqObj _rq = ", tclass)
            b.f("          new %s.BaseParamsReqObj(\"%s\", %s);", 
				tclass, mname, VarName(m.Args[0].Name))
		} else {
			b.f("        %s.BaseParamsReqObj _rq = ", tclass)
            b.f("          new %s.BaseParamsReqObj(\"%s\", %s);", 
				tclass, mname, ParamsAsList(m))
		}
		b.w("        ObjectMapper _m = new ObjectMapper();")
		b.w("        try {")
		b.w("            String _j = _prv.execRPC(_m.writeValueAsString(_rq));")
		if m.ReturnType == "" {
			b.f("          %s.BaseRespObj _resp = ", tclass)
            b.f("            new %s.BaseRespObj(_m, _j);", tclass)
            b.w("          if (_resp.getError() != null) ")
			b.w("            throw new RPCException(_resp.getError());")
		} else {
			rtype := tclass + "." + ServiceResponseType(m.ReturnType)
			b.f("            return new %s(_m, _j).getResult();", rtype)
		}
		b.w("        } catch (java.io.IOException _e) {")
		b.w("            throw new RPCException(-32001, _e.getMessage());")
		b.w("        }")
		b.f("    }")
		b.blank()
	}
	b.w("}")
	return File{JavaFilename(cname), b.b.Bytes()}
}

func (g JavaGenerator) genServiceTypes(p *Package, iface Interface) File {
	cname := iface.Name + "Types"
	retTypes := make(map[string] bool)
	b := StartFile(p)
	b.w("import org.codehaus.jackson.map.ObjectMapper;")
	b.w("import org.codehaus.jackson.JsonNode;")
	b.blank()
	b.f("public class %s {", cname)
	b.blank()
	b.w(typesBoilerplate)
	b.blank()
	for i := 0; i < len(iface.Methods); i++ {
		rtype := iface.Methods[i].ReturnType
		if rtype != "" {
			if _, ok := retTypes[rtype]; !ok {
				retTypes[rtype] = true
				jtype := JavaType(rtype)
			
				b.f("    public static class %sTypeRespObj extends BaseRespObj {", jtype)
				b.f("        %s result;", jtype)
				b.f("        public %sTypeRespObj(ObjectMapper m, String j) throws java.io.IOException {", jtype)
				b.w("            super(m, j);")
				b.w("            if (root.has(\"result\"))")
				if jtype == rtype {
					// custom object, not a built in java type
					b.f("                result = m.treeToValue(root.get(\"result\"), %s.class);", jtype)
				} else if jtype == "String" {
					b.f("                result = root.get(\"result\").asText();")
				} else {
					b.f("                result = root.get(\"result\").as%s();", jtype)
				}
				b.w("        }")
                b.f("        public %s getResult() throws RPCException {", jtype)
                b.w("            if (error != null) throw new RPCException(error);")
                b.w("            else return result;")
                b.w("        }")
                b.w("    }")
                b.blank()
			}
		}
	}
	b.w("}")
	return File{JavaFilename(cname), b.b.Bytes()}
}

func (g JavaGenerator) genRPCException(p *Package) File {
	cname := "RPCException"
	b := StartFile(p)
	b.f("public class %s extends Exception {", cname)
	b.w("    private int code;");
	b.w("    public RPCException(int code, String msg) {")
	b.w("        super(msg); this.code = code;")
	b.w("    }")
	b.w("    public RPCException(RPCError err) {");
	b.w("        this(err.getCode(), err.getMessage());")
	b.w("    }")
	b.w("    public int getCode() { return this.code; }")
	b.w("}")
	return File{JavaFilename(cname), b.b.Bytes()}
}

func (g JavaGenerator) genRPCError(p *Package) File {
	cname := "RPCError"
	b := StartFile(p)
	b.f("public class %s {", cname)
	b.w("    private int code;")
	b.w("    private String message;")
	b.w("    public int getCode() { return this.code; }")
	b.w("    public String getMessage() { return this.message; }")
	b.w("    public void setCode(int c) { this.code = c; }")
	b.w("    public void setMessage(String m) { this.message = m; }")
	b.w("}")
	return File{JavaFilename(cname), b.b.Bytes()}
}

var rpcClientBoilerplate = `    public interface PolygenProvider {
        public String execRPC(String json) throws RPCException;
    }

    class PolygenHttpProvider implements PolygenProvider {
        String _endpointUrl;

        PolygenHttpProvider(String endpointUrl) { this._endpointUrl = endpointUrl; }

        public String execRPC(String json) throws RPCException {
            try {
            java.net.URL url = new java.net.URL(_endpointUrl);
            java.net.URLConnection conn = url.openConnection();
            conn.setDoOutput(true);
            java.io.OutputStreamWriter wr = 
                new java.io.OutputStreamWriter(conn.getOutputStream());
            wr.write(json);
            wr.flush();

            // Get the response
            java.io.BufferedReader rd = 
                new java.io.BufferedReader(new java.io.InputStreamReader(conn.getInputStream()));
            StringBuilder sb = new StringBuilder();
            String line;
            while ((line = rd.readLine()) != null) {
                sb.append(line);
            }
            wr.close();
            rd.close();

            return sb.toString();
            }
            catch (java.io.IOException e) {
                throw new RPCException(-32000, e.getMessage());
            }
        }
    }

    private PolygenProvider _prv;`

var typesBoilerplate = `    public static class BaseReqObj {
        String id;
        String method;

        public BaseReqObj(String method) {
            this.id = java.util.UUID.randomUUID().toString();
            this.method = method;
        }
        public String getId() { return id; }
        public String getJsonrpc() { return "2.0"; }
        public String getMethod() { return method; }
    }

    public static class BaseParamsReqObj extends BaseReqObj {
        Object params;
        public BaseParamsReqObj(String method, Object params) {
            super(method);
            this.params = params;
        }
        public Object getParams() { return params; }
    }

    public static class BaseRespObj {
        String json;
        JsonNode root;
        String id;
        RPCError error;

        public BaseRespObj(ObjectMapper mapper, String json)
            throws java.io.IOException {
            this.json = json;
            root = mapper.readTree(json);
            if (root.has("id")) id = root.get("id").asText();
            if (root.has("error")) error = mapper.treeToValue(root.get("error"), RPCError.class);
        }

        public String getId() { return this.id; }
        public RPCError getError() { return this.error; }
        public String toString() { return json; }
    }`