package main

import (
	"LobbyServer/GlobalDefine"
	"LobbyServer/GlobalFunc"
	"LobbyServer/ServerDefine"
	"fmt"
	"golang.org/x/net/websocket"
	"math"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)
/*
*websocket 里面 hybi.go  文件:
* PayloadType 默认值改成 BinaryFrame,// TextFrame,
*rysinzj 2020424 修改掉 要不然Web报错 Could not decode a text frame as UTF-8
 */
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU()) //最大并发
	//select 心条包 超时断开
	service:=":7106"
	/*---------------TCP being---------------------*/
	/*
		tcpAddr,err:=net.ResolveTCPAddr("tcp4",service)
		checkError(err)
		listener,err:=net.ListenTCP("tcp",tcpAddr)
		checkError(err)
		fmt.Println("开始启动监听端口：",service)
		GlobalDefine.DataBaseSink.StartService(100,100)
		go handleAllDataChannel()
		for{
			conn,err:=listener.Accept()
			if err != nil{
				continue
			}
			go handleClientConnect(conn)
		}*/
	/*---------------TCP end---------------------*/
	/*--------------webSocekt begin -------------------*/
	go handleAllDataChannel()
	http.Handle("/Lobby", websocket.Handler(handleWSConnect))
	if err := http.ListenAndServe(service, nil); err != nil {
		fmt.Println("ListenAndServe:", err)
	}
	/*--------------webSocekt end -------------------*/
}
func EchoServer(ws *websocket.Conn) {
	go func() {
		buf := make([]byte, 100)
		for{
			n, err := ws.Read(buf)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("receive: ", string(buf[:n]))
		}
	}()
	var send int
	for {
		sendStr := strconv.Itoa(send)
		_, err := ws.Write([]byte(sendStr))
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("send: ", sendStr)
		time.Sleep(time.Second)
		send++
		if send>=5 {
			break
		}

	}
}
func handleWSConnect(conn *websocket.Conn){
	if(conn!=nil) {
		if ServerDefine.ConnectIndex>math.MaxInt32{
			ServerDefine.ConnectIndex=1
		}
		ServerDefine.ConnectIndex++;
		dwClientAddr,err:=GlobalFunc.IPString2Long(strings.Split(conn.Request().RemoteAddr, ":")[0])
		if(err!=nil){
			fmt.Println(err);
			conn.Close()
		}else{
			GlobalDefine.SocketEngine.OnSocketAcceptEvent(conn,uint32(dwClientAddr),ServerDefine.ConnectIndex,true)
		}
	}
	for {
		time.Sleep(time.Minute)   //防止退出程序
	}
}
func handleClientConnect(conn net.Conn) {
	if(conn!=nil){
		if ServerDefine.ConnectIndex>math.MaxInt32{
			ServerDefine.ConnectIndex=1
		}
		ServerDefine.ConnectIndex++;
		dwClientAddr,err:=GlobalFunc.IPString2Long(strings.Split(conn.LocalAddr().String(), ":")[0])
		if(err!=nil){
			conn.Close()
		}else{
			GlobalDefine.SocketEngine.OnSocketAcceptEvent(conn,uint32(dwClientAddr),ServerDefine.ConnectIndex,false)
		}
	}
}
/*func handleClientClose(){
	for{
		data,isClose := <-ServerDefine.GlobalCloseChan
		if !isClose {
			fmt.Println("channel closed!")
			break
		}
		GlobalDefine.SocketEngine.OnSocketCloseEvent(data.RoundID,data.Index)
		GlobalDefine.AttemperEngineSink.OnEventSocketClose(data.RoundID,data.Index)
	}

}
func handleAttemperEngineSocketRead() {
	for{
		data,isClose := <-ServerDefine.RcSocketDataChan
		if !isClose {
			fmt.Println("receive data channel closed!")
			break
		}
		if !GlobalDefine.AttemperEngineSink.OnEventSocketRead(data.MainCmdID,data.SubCmdID,data.SocketData,data.SocketSize,data.RoundID,data.Index,data.ClientIP) {
			 GlobalFunc.OnSendCloseSocket(data.RoundID,data.Index)
		}
	}
}
func handleSendToClient(){
	for{
		data,isClose := <-ServerDefine.AllDataChan
		if !isClose {
			fmt.Println("senddata channel closed!")
			break
		}
		if !GlobalDefine.SocketEngine.OnSocketSendDataEvent(data) {
			GlobalFunc.OnSendCloseSocket(data.RoundID,data.Index)
		}
	}
}
func handleTimeEngineData(){
	for{
		data,isClose := <-ServerDefine.TimeDataChan
		if !isClose {
			fmt.Println("timer channel closed!")
			break
		}
		if data.TimeNotic==0{
			GlobalDefine.TimerEngine.SetTimer(data)
		}else if data.TimeNotic==1{
			GlobalDefine.AttemperEngineSink.OnEventTimer(data.TimeID,data.TimeParam)
		}else{
			GlobalDefine.TimerEngine.KillTimer(data.TimeID)
		}
	}
}
func handleDataBaseSocketRead(){
	for{
		data,isClose := <-ServerDefine.RcDataBaseChan
		if !isClose {
			fmt.Println("receive database channel closed!")
			break
		}
		if data.EventType==0{
			GlobalDefine.DataBaseSink.OnDataBaseRequest(data.RequestID,data.SocketData,data.SocketSize,data.RoundID,data.Index)
		}else if data.EventType==1 {
			GlobalDefine.AttemperEngineSink.OnEventDataBase(data.SocketData,data.SocketSize,data.RequestID,data.RoundID,data.Index)
		}else{
			fmt.Println("非法数据库消息")
		}
	}
}*/
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		GlobalFunc.OnWriteServerLog(fmt.Sprintf("Fatal error: %s", err.Error()))
		os.Exit(1)
	}
}
func handleAllDataChannel(){
	for{
		allData,isClose := <-ServerDefine.AllDataChan
		if !isClose {
			GlobalFunc.OnWriteServerLog("receive  channel closed!")
			break
		}
		if allData.DataType==ServerDefine.ALL_DATA_TYPE_EXIT{
			data, ok := allData.DataBody.(ServerDefine.ClientIdentity)
			if ok{
				GlobalDefine.SocketEngine.OnSocketCloseEvent(data.RoundID,data.Index)
				GlobalDefine.AttemperEngineSink.OnEventSocketClose(data.RoundID,data.Index)
			}
		}else if allData.DataType== ServerDefine.ALL_DATA_TYPE_REC{
			data, ok := allData.DataBody.(ServerDefine.CMD_SocketData)
			if ok{
				if !GlobalDefine.AttemperEngineSink.OnEventSocketRead(data.MainCmdID,data.SubCmdID,data.SocketData,data.SocketSize,data.RoundID,data.Index,data.ClientIP) {
					GlobalFunc.OnSendCloseSocket(data.RoundID,data.Index)
				}
			}
		}else if allData.DataType==ServerDefine.ALL_DATA_TYPE_SEND{
			data, ok := allData.DataBody.(ServerDefine.CMD_SocketData)
			if ok{
				if !GlobalDefine.SocketEngine.OnSocketSendDataEvent(data) {
					GlobalFunc.OnSendCloseSocket(data.RoundID,data.Index)
				}
			}
		}else if allData.DataType==ServerDefine.ALL_DATA_TYPE_TIME{
			data, ok := allData.DataBody.(ServerDefine.TagTimeParam)
			if ok{
				if data.TimeNotic==0{
					GlobalDefine.TimerEngine.SetTimer(data)
				}else if data.TimeNotic==1{
					GlobalDefine.AttemperEngineSink.OnEventTimer(data.TimeID,data.TimeParam)
				}else{
					GlobalDefine.TimerEngine.KillTimer(data.TimeID)
				}
			}
		}else if allData.DataType==ServerDefine.ALL_DATA_TYPE_DATABASE{
			data, ok := allData.DataBody.(ServerDefine.NTY_DataBaseEvent)
			if ok{
				if data.EventType==0{
					if GlobalDefine.DataBaseSink.OnDataBaseRequest(data.RequestID,data.SocketData,data.SocketSize,data.RoundID,data.Index)==false{
						GlobalFunc.OnSendCloseSocket(data.RoundID,data.Index)
					}
				}else if data.EventType==1 {
					GlobalDefine.AttemperEngineSink.OnEventDataBase(data.SocketData,data.SocketSize,data.RequestID,data.RoundID,data.Index)
				}else{
					GlobalFunc.OnWriteServerLog("非法数据库消息")
				}
			}
		}
		/*select{
		 	case data,isClose := <-ServerDefine.GlobalCloseChan:
				if !isClose {
					fmt.Println("channel closed!")
					break
				}
				GlobalDefine.SocketEngine.OnSocketCloseEvent(data.RoundID,data.Index)
				GlobalDefine.AttemperEngineSink.OnEventSocketClose(data.RoundID,data.Index)
			case data,isClose := <-ServerDefine.RcSocketDataChan:
				if !isClose {
					fmt.Println("receive data channel closed!")
					break
				}
				if !GlobalDefine.AttemperEngineSink.OnEventSocketRead(data.MainCmdID,data.SubCmdID,data.SocketData,data.SocketSize,data.RoundID,data.Index,data.ClientIP) {
					GlobalFunc.OnSendCloseSocket(data.RoundID,data.Index)
				}
			case data,isClose := <-ServerDefine.SdSocketDataChan:
				if !isClose {
					fmt.Println("senddata channel closed!")
					break
				}
				if !GlobalDefine.SocketEngine.OnSocketSendDataEvent(data) {
					GlobalFunc.OnSendCloseSocket(data.RoundID,data.Index)
				}
			case data,isClose := <-ServerDefine.TimeDataChan:
				if !isClose {
					fmt.Println("timer channel closed!")
					break
				}
				if data.TimeNotic==0{
					GlobalDefine.TimerEngine.SetTimer(data)
				}else if data.TimeNotic==1{
					GlobalDefine.AttemperEngineSink.OnEventTimer(data.TimeID,data.TimeParam)
				}else{
					GlobalDefine.TimerEngine.KillTimer(data.TimeID)
				}
			case data,isClose := <-ServerDefine.RcDataBaseChan:
				if !isClose {
					fmt.Println("receive database channel closed!")
					break
				}
				if data.EventType==0{
					GlobalDefine.DataBaseSink.OnDataBaseRequest(data.RequestID,data.SocketData,data.SocketSize,data.RoundID,data.Index)
				}else if data.EventType==1 {
					GlobalDefine.AttemperEngineSink.OnEventDataBase(data.SocketData,data.SocketSize,data.RequestID,data.RoundID,data.Index)
				}else{
					fmt.Println("非法数据库消息")
				}
		 	default:
		 		fmt.Println("Default 处理数据")
				time.Sleep(50 * time.Millisecond)
		 }*/
	}
}