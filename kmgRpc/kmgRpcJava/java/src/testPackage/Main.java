package testPackage;

import testPackage.RpcDemo;
import java.lang.System;
import java.time.Instant;
import java.util.Date;

public class Main {
    public static void main(String[] args)
    {
        try {
            RpcDemo.ConfigDefaultClient("http://127.0.0.1:34895", "abc psk");
            String out = RpcDemo.GetDefaultClient().PostScoreInt("1", 1);
            if (!out.equals("1")){
                throw new Exception("1");
            }
            Date inT = RpcDemo.KmgTime.ParseGolangDate("2001-01-01T01:01:01+08:00");
            Date outT = RpcDemo.GetDefaultClient().DemoTime2(inT);
            if (!RpcDemo.KmgTime.ParseGolangDate("2001-01-01T02:01:01.001+08:00").equals(outT)){
                throw new Exception("2 "+RpcDemo.KmgTime.FormatGolangDate(outT));
            }
            inT = RpcDemo.KmgTime.ParseGolangDate("2001-01-01T01:01:01.1+08:00");
            outT = RpcDemo.GetDefaultClient().DemoTime(inT);
            if (!RpcDemo.KmgTime.ParseGolangDate("2001-01-01T02:01:01.101+08:00").equals(outT)){
                throw new Exception("3 "+RpcDemo.KmgTime.FormatGolangDate(outT));
            }
            System.out.println("Success");
        }catch (Exception e){
            System.out.println(e.getMessage());
            e.printStackTrace();
        }
    }
}
