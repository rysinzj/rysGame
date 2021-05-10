package AttemperManage

import (
	"LobbyServer/DataBaseDefine"
	"LobbyServer/GlobalFunc"
	"LobbyServer/ServerDefine"
	"bytes"
	"encoding/binary"
)

type AttemperEngineSink struct{

}
func NewAttemperEngineSink() *AttemperEngineSink{
	attemp :=&AttemperEngineSink{ }
	return attemp
}

func(a *AttemperEngineSink)OnEventTimer(wTimerID uint32 ,wBindParam uint32)bool{
	println("定时器通知:",wTimerID)
	return false
}
func(a *AttemperEngineSink)OnEventSocketRead(MainCmdID uint16,SubCmdID uint16,pDataBuffer[]byte,wDataSize uint16,wRoundID int64,wIndex int64,wClientIP uint32)bool{
	switch MainCmdID{
		case ServerDefine.MDM_GR_LOGON:
			return a.OnSocketMainLogon(SubCmdID,pDataBuffer,wDataSize,wRoundID,wIndex,wClientIP)
	}
	defer GlobalFunc.OnRecoverError()
	return false
}
func(a *AttemperEngineSink)OnEventSocketClose(wRoundID int64,wIndex	 int64)bool{
	defer GlobalFunc.OnRecoverError()
	return true
}
func(a *AttemperEngineSink)OnSocketMainLogon(SubCmdID uint16,dataBuffer[]byte,wDataSize uint16,wRoundID int64,wIndex int64,wClientIP uint32)bool{
	switch SubCmdID {
	case ServerDefine.SUB_GP_REGET_USERMONEY:
		return a.OnRqUserMoney(dataBuffer,wDataSize,wRoundID,wIndex,wClientIP)
	case ServerDefine.SUB_GP_DAILYTASK_INFO_KIND:
		return a.OnRqUserTaskInfo(dataBuffer,wDataSize,wRoundID,wIndex,wClientIP)
	case ServerDefine.SUB_GP_DAILYTASK_OP:
		return a.OnRqUserDailyTaskOP(dataBuffer,wDataSize,wRoundID,wIndex,wClientIP)
	case ServerDefine.SUB_GP_USER_BANKOP:
		return a.OnRqUserBankOP(dataBuffer,wDataSize,wRoundID,wIndex,wClientIP)
	case ServerDefine.SUB_GP_TREASURE_RANKING:
		return a.OnRqUserScoreRank(dataBuffer,wDataSize,wRoundID,wIndex,wClientIP)
	}
	return false
}
/*------------------请求用户金币-----------------*/
func(a *AttemperEngineSink)OnRqUserMoney(dataBuffer[]byte,wDataSize uint16,wRoundID int64,wIndex int64,wClientIP uint32)bool{
	if wDataSize<40 {
		return false
	}
	bytesBuffer := bytes.NewBuffer(dataBuffer)
	UserMoney:=&ServerDefine.CMD_GP_UserMoney{}
	err :=binary.Read(bytesBuffer, binary.LittleEndian, UserMoney)
	if err != nil {
		return false
	}

	var DBRReGetMoney DataBaseDefine.DBR_GP_ReGetMoney
	DBRReGetMoney.UserID=UserMoney.UserID
	copy(DBRReGetMoney.PassWord[:],UserMoney.PassWord[:])
	DBRReGetMoney.IngotType=0
	GlobalFunc.PostDataBaseEvent(DataBaseDefine.DBR_GP_REGET_MONEY,wIndex,wRoundID,DBRReGetMoney,uint16(binary.Size(DBRReGetMoney)))
	return true
}
/*------------------请求用户任务-----------------*/
func(a *AttemperEngineSink)OnRqUserTaskInfo(dataBuffer[]byte,wDataSize uint16,wRoundID int64,wIndex int64,wClientIP uint32)bool{
	if wDataSize<8 {
		return false
	}
	bytesBuffer := bytes.NewBuffer(dataBuffer)
	UserTaskInfo:=&ServerDefine.CMD_GP_UserDailyTaskInfo_KindID{}
	err :=binary.Read(bytesBuffer, binary.LittleEndian, UserTaskInfo)
	if err != nil {
		return false
	}
	var DBRUserTaskInfo DataBaseDefine.DBR_GP_UserDailyTaskInfo
	DBRUserTaskInfo.UserID=UserTaskInfo.UserID
	DBRUserTaskInfo.KindID=UserTaskInfo.KindID
	GlobalFunc.PostDataBaseEvent(DataBaseDefine.DBR_GP_RETURN_DAILYTASKINFO,wIndex,wRoundID,DBRUserTaskInfo,uint16(binary.Size(DBRUserTaskInfo)))
	return true
}
/*------------------领取任务奖励-----------------*/
func(a *AttemperEngineSink)OnRqUserDailyTaskOP(dataBuffer[]byte,wDataSize uint16,wRoundID int64,wIndex int64,wClientIP uint32)bool{
	if wDataSize!=44 {
		return false
	}
	bytesBuffer := bytes.NewBuffer(dataBuffer)
	UserDailyTaskInfoOp:=&ServerDefine.CMD_GP_UserDailyTaskInfoOp{}
	err :=binary.Read(bytesBuffer, binary.LittleEndian, UserDailyTaskInfoOp)
	if err != nil {
		return false
	}
	var DBRUserTaskOp DataBaseDefine.DBR_GP_UserDailyTaskInfoOp
	DBRUserTaskOp.TaskID=UserDailyTaskInfoOp.TaskID
	DBRUserTaskOp.UserID=UserDailyTaskInfoOp.UserID
	copy(DBRUserTaskOp.UserPass[:],UserDailyTaskInfoOp.UserPass[:])
	GlobalFunc.PostDataBaseEvent(DataBaseDefine.DBR_GP_UPDATE_DAILYTASKOP,wIndex,wRoundID,DBRUserTaskOp,uint16(binary.Size(DBRUserTaskOp)))
	return true
}
/*----------------------------银行操作-------------------------*/
func(a *AttemperEngineSink)OnRqUserBankOP(buffer[]byte,wDataSize uint16,wRoundID ,wIndex int64,wClientIP uint32)bool{
	if wDataSize!=48 {
		return false
	}
	bytesBuffer := bytes.NewBuffer(buffer)
	UserBankOp:=&ServerDefine.CMD_GP_UserBankOp{}
	err :=binary.Read(bytesBuffer, binary.LittleEndian, UserBankOp)
	if err != nil {
		return false
	}
	var DBRUserBanOp DataBaseDefine.DBR_GP_UserBankOp
	DBRUserBanOp.UserID=UserBankOp.UserID
	DBRUserBanOp.Score=UserBankOp.Score
	DBRUserBanOp.ClientIP=wClientIP
	DBRUserBanOp.Type=UserBankOp.Type
	copy(DBRUserBanOp.BankPass[:],UserBankOp.BankPass[:])
	GlobalFunc.PostDataBaseEvent(DataBaseDefine.DBR_GP_OPERAT_BANK,wIndex,wRoundID,DBRUserBanOp,uint16(binary.Size(DBRUserBanOp)))
	return true
}
func(a *AttemperEngineSink)OnRqUserScoreRank(buffer[]byte,wDataSize uint16,wRoundID ,wIndex int64,wClientIP uint32)bool{
	if wDataSize!=4 {
		return false
	}
	bytesBuffer := bytes.NewBuffer(buffer)
	UserRank:=&ServerDefine.CMD_GP_UserRank{}
	err :=binary.Read(bytesBuffer, binary.LittleEndian, UserRank)
	if err != nil {
		return false
	}
	var DBRUserRank DataBaseDefine.DBR_GP_UserRank
	DBRUserRank.UserID=UserRank.UserID
	GlobalFunc.PostDataBaseEvent(DataBaseDefine.DBR_GP_TREASURE_RANKING,wIndex,wRoundID,UserRank,uint16(binary.Size(UserRank)))
	return true
}
/*------------------------数据库反馈------------------------*/
func(a *AttemperEngineSink)OnEventDataBase(dataBuffer []byte,wDataSize uint16,wRequestID uint16,wRoundID int64,wIndex	 int64)bool{
	switch wRequestID {
	case DataBaseDefine.DBR_GP_RETURN_MONEY:
		return a.OnDBReturnMoney(dataBuffer,wDataSize,wRoundID,wIndex)
	case DataBaseDefine.DBR_GP_RETURN_DAILYTASKINFO:
		return a.OnDBReturnDailyTask(dataBuffer,wDataSize,wRequestID,wRoundID,wIndex)
	case DataBaseDefine.DBR_GP_RETURN_DAILYTASKINFO_FINISH:
		return a.OnDBReturnDailyTask(dataBuffer,wDataSize,wRequestID,wRoundID,wIndex)
	case DataBaseDefine.DBR_GP_UPDATE_DAILYTASKOP:
		return a.OnDBReturUpdateDailyOP(dataBuffer,wDataSize,wRoundID,wIndex)
	case DataBaseDefine.DBR_GP_OPERAT_BANK:
		return a.OnDBReturnOperateBank(dataBuffer,wDataSize,wRoundID,wIndex)
	case DataBaseDefine.DBR_GR_TREASURE_RANKING:
		fallthrough
	case DataBaseDefine.DBR_GR_TREASURE_RANKING_FINISH:
		return a.OnDBReturnTreasureRank(dataBuffer,wDataSize,wRequestID,wRoundID,wIndex)
	}
	defer GlobalFunc.OnRecoverError()
	return false
}
//获取金币
func(a *AttemperEngineSink)OnDBReturnMoney(pDataBuffer[]byte,wDataSize uint16,wRoundID int64,wIndex int64)bool {
	bytesBuffer := bytes.NewBuffer(pDataBuffer)
	ReGetMoney:=&DataBaseDefine.DBR_GP_UserMoney{}
	err :=binary.Read(bytesBuffer, binary.LittleEndian, ReGetMoney)
	if err != nil {
		return false
	}
	ReturnUserMoney:=ServerDefine.CMD_GP_ReturnUserMoney{}
	ReturnUserMoney.RoomCard=ReGetMoney.RoomCard
	ReturnUserMoney.ReviveCard=ReGetMoney.ReviveCard
	ReturnUserMoney.RedNum=ReGetMoney.RedNum
	ReturnUserMoney.Lottery=ReGetMoney.Lottery
	ReturnUserMoney.Ingot=ReGetMoney.Ingot
	ReturnUserMoney.GradeScore=ReGetMoney.GradeScore
	ReturnUserMoney.BankScore=ReGetMoney.BankScore
	ReturnUserMoney.Score=ReGetMoney.Score
	GlobalFunc.OnSendData(wRoundID,wIndex,ServerDefine.MDM_GR_LOGON,ServerDefine.SUB_GP_REGET_USERMONEY,&ReturnUserMoney,uint16(binary.Size(ReturnUserMoney)))
	GlobalFunc.OnSendCloseSocket(wRoundID,wIndex)
	return true
}
//任务返回
func(a *AttemperEngineSink)OnDBReturnDailyTask(dataBuffer []byte,wDataSize uint16,wRequestID uint16,wRoundID int64,wIndex	 int64)bool{
	if wRequestID==DataBaseDefine.DBR_GP_RETURN_DAILYTASKINFO_FINISH{
		GlobalFunc.OnSendData(wRoundID,wIndex,ServerDefine.MDM_GR_LOGON,ServerDefine.SUB_GP_DAILYTASK_INFO_FINISH,dataBuffer,wDataSize)
		GlobalFunc.OnSendCloseSocket(wRoundID,wIndex)
	}else{
		GlobalFunc.OnSendData(wRoundID,wIndex,ServerDefine.MDM_GR_LOGON,ServerDefine.SUB_GP_DAILYTASK_INFO,dataBuffer,wDataSize)
	}
	return true
}
//领取返回
func(a *AttemperEngineSink)OnDBReturUpdateDailyOP(dataBuffer []byte,wDataSize uint16,wRoundID int64,wIndex	 int64)bool{
	bytesBuffer := bytes.NewBuffer(dataBuffer)
	Result:=&DataBaseDefine.DBR_GP_UserDailyTaskInfoOpResult{}
	err :=binary.Read(bytesBuffer, binary.LittleEndian, Result)
	if err != nil {
		return false
	}
	var UserDailyTaskInfoOpResult ServerDefine.CMD_GP_UserDailyTaskInfoOpResult
	UserDailyTaskInfoOpResult.GiftType=Result.GiftType
	UserDailyTaskInfoOpResult.Result=Result.Result
	UserDailyTaskInfoOpResult.TaskValue=Result.TaskValue
	UserDailyTaskInfoOpResult.TaskID=Result.TaskID
	copy(UserDailyTaskInfoOpResult.ErrorDescribe[:],Result.ErrorDescribe[:])
	GlobalFunc.OnSendData(wRoundID,wIndex,ServerDefine.MDM_GR_LOGON,ServerDefine.SUB_GP_DAILYTASK_OP,dataBuffer,wDataSize)
	GlobalFunc.OnSendCloseSocket(wRoundID,wIndex)
	return true
}
//银行操作返回
func(a *AttemperEngineSink)OnDBReturnOperateBank(dataBuffer []byte,wDataSize uint16,wRoundID int64,wIndex	 int64)bool {
	bytesBuffer := bytes.NewBuffer(dataBuffer)
	Result := &DataBaseDefine.DBR_GP_OperateResult{}
	err := binary.Read(bytesBuffer, binary.LittleEndian, Result)
	if err != nil {
		return false
	}
	var OperateResult ServerDefine.CMD_GP_OperateResult
	OperateResult.ErrorCode=Result.ErrorCode
	OperateResult.Bank=Result.Bank
	OperateResult.Score=Result.Score
	copy(OperateResult.ErrorDescribe[:],Result.ErrorDescribe[:])
	GlobalFunc.OnSendData(wRoundID,wIndex,ServerDefine.MDM_GR_LOGON,ServerDefine.SUB_GP_USER_BANKOP,dataBuffer,wDataSize)
	GlobalFunc.OnSendCloseSocket(wRoundID,wIndex)
	return true
}
//排行榜
func(a *AttemperEngineSink)OnDBReturnTreasureRank(dataBuffer []byte,wDataSize uint16,wRequestID uint16,wRoundID int64,wIndex	 int64)bool{
	if wRequestID==DataBaseDefine.DBR_GR_TREASURE_RANKING_FINISH{
	GlobalFunc.OnSendData(wRoundID,wIndex,ServerDefine.MDM_GR_LOGON,ServerDefine.SUB_GR_TREASURE_RANKING_FINISH,dataBuffer,wDataSize)
	GlobalFunc.OnSendCloseSocket(wRoundID,wIndex)
	}else{
	GlobalFunc.OnSendData(wRoundID,wIndex,ServerDefine.MDM_GR_LOGON,ServerDefine.SUB_GR_TREASURE_RANKING,dataBuffer,wDataSize)
	}
	return true
}