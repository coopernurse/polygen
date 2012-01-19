package foolib;

import java.util.List;
import java.util.Map;
import java.util.HashMap;

public class Example1Client {

    public static void main(String argv[]) throws Exception {
        SampleServiceClient client = new SampleServiceClient("http://localhost:9009");
        
        Person p = new Person();
        p.setEmail("foo@bar.com");
        p.setName("my name");
        System.out.println("Calling Create()");
        Result r = client.Create(p);
        System.out.println("  note=" + r.getNote() + "  code=" + r.getCode());

        System.out.println("Calling Add()");
        System.out.println(" 2+3=" + client.Add(2L, 3L));

        long start = System.currentTimeMillis();
        for (int i = 0; i < 10000; i++) {
            //client.Add(2L, 3L);
        }
        System.err.println("elapsed: " + (System.currentTimeMillis()-start));

        System.out.println("Calling StoreName()");
        client.StoreName("bob");

        System.out.println("Say_Hi(): " + client.Say_Hi());

        System.out.println("Calling getPeople()");
        Map<String,String> keys = new HashMap<String,String>();
        keys.put("key1", "val1");
        keys.put("key2", "val2");
        List<Person> people = client.getPeople(keys);
        for (Person per : people) {
            System.out.println("  person.name=" + per.getName());
        }

        System.out.println("done!");
    }

}