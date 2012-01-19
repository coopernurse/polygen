package foolib;

import java.util.List;
import java.util.ArrayList;
import java.util.Map;
import java.util.HashMap;

public class Example1Server implements SampleService {

    public static void main(String argv[]) throws Exception {
        Example1Server es = new Example1Server();
        SampleServiceHttpServer server = 
            new SampleServiceHttpServer(new SampleServiceDispatcher(es));

        System.out.println("Starting server on port 9009");
        server.serve(9009);
    }

    public Result Create(Person p) throws RPCException {
        System.out.println("creating person with name: " + p.getName());
        Result r = new Result();
        r.setNote("this is a note");
        r.setCode(393L);
        return r;
    }

    public Long Add(Long a, Long b) throws RPCException {
        return a+b;
    }

    public void StoreName(String name) throws RPCException {
        System.out.println("StoreName() name=" + name);
    }

    public String Say_Hi() throws RPCException {
        return "howdy ho";
    }

    public List<Person> getPeople(Map<String,String> params) throws RPCException {

        List<Person> list = new ArrayList<Person>();

        for (String key : params.keySet()) {
            System.out.println("getPeople key=" + key);
            Person p = new Person();
            p.setName("name with key: " + key + " val: " + params.get(key));
            list.add(p);
        }

        return list;
    }

}