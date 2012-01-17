package polygenlib

import (
	"fmt"
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
		files = append(files, g.genServiceDispatcher(p, iface))
		files = append(files, g.genServiceHttpServer(p, iface))
		files = append(files, g.genServiceClient(p, iface))
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

func (g JavaGenerator) genServiceDispatcher(p *Package, iface Interface) File {
	cname := iface.Name + "Dispatcher"
	b := StartFile(p)
	b.w("import org.codehaus.jackson.map.ObjectMapper;")
	b.w("import org.codehaus.jackson.JsonNode;")
	b.w("import org.codehaus.jackson.node.ObjectNode;")
	b.blank()
	b.f("public class %s {", cname)
	b.blank()
	b.f("    private %s _service;", iface.Name)
	b.blank()
	b.f("    public %s(%s service) {", cname, iface.Name)
	b.w("        this._service = service;")
	b.w("    }")
	b.blank()
	b.w("    public String exec(String _json) {")
	b.w("        ObjectMapper _m = new ObjectMapper();")
	b.w("        ObjectNode _resp = _m.createObjectNode();")
	b.w("        JsonNode _r = null;")
	b.w("        String _id = null;")
	b.w("        try { _r = _m.readTree(_json); }")
	b.w("        catch (java.io.IOException _e) { ")
	b.w("          return rpcErr(_resp,-32700, \"Parse error: \" +_json,_id);")
	b.w("        }")
	b.w("        if (_r.has(\"id\")) { _id = _r.get(\"id\").asText(); }")
	b.w("        if (_r.has(\"method\")) {")
	b.w("          String _meth = _r.get(\"method\").asText();")
	b.w("          try {")
	for i := 0; i < len(iface.Methods); i++ {
		m := iface.Methods[i]
		
		
		b.raw("            ")
		if i > 0 {
			b.raw("else ")
		}
		b.f("if (_meth.equals(\"%s_%s\")) {", iface.Name, m.Name)
		jtype := JavaType(m.ReturnType)
		params := ""
		if len(m.Args) > 0 {
			b.w("              JsonNode _par = _r.get(\"params\");")
		}
		for x := 0; x < len(m.Args); x++ {
			prefix := "_par"
			if len(m.Args) > 1 {
				prefix  += fmt.Sprintf(".get(%d)", x)
			}
			arg := m.Args[x]
			if x > 0 {
				params += ","
			}
			argjtype := JavaType(arg.Type)
			if arg.Type == argjtype {
				params += fmt.Sprintf("_m.treeToValue(%s, %s.class)", prefix, arg.Type)
			} else if argjtype == "String" {
				params += prefix + ".asText()"
			} else {
				params += prefix + ".as" + argjtype + "()"
			}
		}
		if m.ReturnType == "" {
			b.f("              _service.%s(%s);", m.Name, params)
			b.f("              _resp.put(\"result\", true);")
		} else if m.ReturnType == jtype {
			b.f("              _resp.put(\"result\", _m.valueToTree(_service.%s(%s)));", m.Name, params)
		} else {
			b.f("              _resp.put(\"result\", _service.%s(%s));", m.Name, params)
		} 
		b.w("              _resp.put(\"jsonrpc\", \"2.0\");")
		b.w("              _resp.put(\"id\", _id);")
		b.w("              return _resp.toString();")
		b.w("            }")
	}
    b.w("            else { return rpcErr(_resp, -32601, \"Method not found: \" + _meth, _id); }")
	b.w("          }")
	b.w("          catch (RPCException e) { return rpcErr(_resp, e.getCode(), e.getMessage(), _id); }")
	b.w("          catch (Throwable t) { return rpcErr(_resp, -32005, \"Unknown error: \" + t.getMessage(), _id); }")
	b.w("        }")
	b.w("        else { return rpcErr(_resp, -32600, \"Invalid Request. method missing: \" +_json, _id); }")
	b.w("    }")
	b.w(dispatcherRpcErr)
	b.w("}")
	return File{JavaFilename(cname), b.b.Bytes()}
}

func (g JavaGenerator) genServiceHttpServer(p *Package, iface Interface) File {
	cname := iface.Name + "HttpServer"
	b := StartFile(p)
	b.f("public class %s {", cname)
	b.blank()
	b.f("    private %sDispatcher dispatcher;", iface.Name)
	b.blank()
	b.f("    public %s(%sDispatcher d) { ", cname, iface.Name)
	b.w("        this.dispatcher = d;")
	b.w("    }")
	b.blank()
	b.w(httpServerBoilerplate)
	b.w("}")
	return File{JavaFilename(cname), b.b.Bytes()}
}

func (g JavaGenerator) genServiceClient(p *Package, iface Interface) File {
	cname := iface.Name + "Client"
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

var httpServerBoilerplate = `    private java.util.List<Worker> pool;

    public void serve(int port) throws Exception {
        int count = 0;
        pool = new java.util.ArrayList<Worker>();
        java.net.ServerSocket ss = new java.net.ServerSocket(port);
        while (true) {

            java.net.Socket s = ss.accept();

            Worker w = null;
            synchronized (pool) {
                if (pool.isEmpty()) {
                    Worker ws = new Worker(pool);
                    ws.setSocket(s);
                    count++;
                    Thread t = new Thread(ws, "ServiceHttpWorker-"+count);
                    t.setDaemon(true);
                    t.start();
                } else {
                    w = pool.remove(0);
                    w.setSocket(s);
                }
            }
        }
    }

    class Worker implements Runnable {
        final byte[] EOL = {(byte)'\r', (byte)'\n' };
        final String STATUS_OK = "200 OK";
        final String STATUS_ERR = "500 Internal Server Error";

        java.net.Socket s;
        java.util.List<Worker> pool;
        java.util.List<Byte> bytes;

        Worker(java.util.List<Worker> pool) {
            this.s = null;
            this.pool = pool;
            bytes = new java.util.ArrayList<Byte>(2048);
        }

        synchronized void setSocket(java.net.Socket s) {
            this.s = s;
            notify();
        }

        public synchronized void run() {
            while (true) {
                if (s == null) {
                    try {
                        wait();
                    } catch (InterruptedException e) {
                        continue;
                    }
                }
                try {
                    handleClient();
                } catch (Exception e) {
                    e.printStackTrace();
                }

                s = null;
                synchronized (this.pool) {
                    if (this.pool.size() >= 5) {
                        return;
                    } else {
                        this.pool.add(this);
                    }
                }
            }
        }
        
        void handleClient() throws Exception {
            java.io.InputStream in = new java.io.BufferedInputStream(s.getInputStream());
            java.io.PrintStream out = new java.io.PrintStream(s.getOutputStream());
            s.setSoTimeout(30000);
            s.setTcpNoDelay(true);

            String status = STATUS_OK;
            String ctype = "text/plain";
            String outStr = null;
            bytes.clear();

            try {
                StringBuilder sb = new StringBuilder();
                int r = 0;
                int b = 0;
                int lastb = 0;
                boolean inbody = false;
                int contentLen = 0;
                
                while ((b = in.read()) > -1) {
                    if (inbody) {
                        bytes.add((byte)b);
                        if (bytes.size() == contentLen) {
                            break;
                        }
                    }
                    else {
                        if (b == '\n' && lastb == '\r') {
                            String header = 
                                new String(toArr(bytes, bytes.size()-1), 
                                           "utf-8");
                            if (header.trim().equals("")) {
                                inbody = true;
                            }
                            else if (header.toLowerCase().startsWith("content-length")) {
                                header = header.toLowerCase();
                                int start = header.indexOf(":");
                                contentLen = Integer.parseInt(header.substring(start+1).trim());
                            }
                            bytes.clear();
                        } else {
                            bytes.add((byte)b);
                        }
                    }
                    lastb = b;
                }
                String post = new String(toArr(bytes, bytes.size()), "utf-8");

                System.out.println("Server, got: " + post);
                outStr = dispatcher.exec(post);

                out.print("HTTP/1.0 ");
                out.print(status);
                out.write(EOL);
                out.print("Content-Length: ");
                out.print(outStr.length());
                out.print(EOL);
                out.print("Content-Type: ");
                out.print(ctype);
                out.write(EOL);
                out.write(EOL);
                out.print(outStr);
                out.flush();
                s.close();

                System.out.println("Server, sent: " + outStr);
            } finally {
                in.close();
                out.close();
                s.close();
            }
        }

        byte[] toArr(java.util.List<Byte> list, int size) {
            byte[] arr = new byte[size];
            for (int i = 0; i < size; i++) {
                arr[i] = list.get(i);
            }
            return arr;
        }

    }`

var dispatcherRpcErr = `    private String rpcErr(ObjectNode resp, int code, String msg, String id) {
        resp.put("jsonrpc", "2.0");
        resp.put("id", id);
        ObjectNode err = resp.putObject("error");
        err.put("code", code);
        err.put("message", msg);
        return resp.toString();
    }`