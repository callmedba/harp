package main

import (
  "arp"
  "fmt"
  "flag"
  "config"
  "net"
)

func parseFlags(conf *string, dev *string) {
  flag.StringVar(conf, "conf", "", "config description.\n"+ 
                 "e.g. file:/path/to/file;\n"+
                 "      etcd:node_name@10.10.10.1:2379,10.10.10.2:2379/prefix")
  flag.StringVar(dev, "dev", "", "network device. e.g. eth0")
  flag.Parse()
}

func main() {
  var confDesc, devName string
  parseFlags(&confDesc)
  conf, err := config.NewConfig(confDesc, devName)
  if err != nil {
    fmt.Println(err)
    return
  }
  dev, err = net.InterfaceByName(devName)
  listener, err := arp.NewListener()
  if err != nil {
    fmt.Println(err)
    return
  }
  defer listener.Close()
  
  conf.SetOnChange(func(ip string, mac string){
    if mac == "" {
      return
    }
    ipAddr:= net.PasreIP(ip)
    hwAddr, err:= net.ParseMac(mac)
    if err != nil {
      return
    }
    self.listener.Send(arp.GratuitousArpOp(ipAddr, hwAddr), arp.AddrFromInterface(dev))
  })
  listener.Listen(func(op arp.ArpOp, addr arp.Addr){
    if op.Op != arp.ARP_REQUEST {
      return
    }
    mac := conf.GetMacByIp(op.RcptIpAddr.String())
    //fmt.Printf("%s -> %s\n", op.RcptIpAddr.String(), mac)
    if mac == "" {
      return
    }
    var reply arp.ArpOp
    reply.Op = arp.ARP_REPLY
    reply.SndrIpAddr = op.RcptIpAddr
    reply.SndrHwAddr, err = net.ParseMAC(mac)
    if err != nil {
      return
    }
    reply.RcptIpAddr = op.SndrIpAddr
    reply.RcptHwAddr = op.SndrHwAddr
    listener.Send(reply, addr)
  })
}



