package app

import (
	"fmt"
	"github.com/kolo/xmlrpc"
)

type ProcessInfo struct {
	Name          string `xmlrpc:"name" json:"name"`
	Group         string `xmlrpc:"group" json:"group"`
	Description   string `xmlrpc:"description" json:"description"`
	Start         int    `xmlrpc:"start" json:"start"`
	Stop          int    `xmlrpc:"stop" json:"stop"`
	Now           int    `xmlrpc:"now" json:"now"`
	State         int    `xmlrpc:"state" json:"state"`
	StateName     string `xmlrpc:"statename" json:"statename"`
	SpawnErr      string `xmlrpc:"spawnerr" json:"spawnerr"`
	ExitStatus    int    `xmlrpc:"exitstatus" json:"exitstatus"`
	LogFile       string `xmlrpc:"logfile" json:"logfile"`
	StdoutLogfile string `xmlrpc:"stdout_logfile" json:"stdout_logfile"`
	StderrLogfile string `xmlrpc:"stderr_logfile" json:"stderr_logfile"`
	Pid           int    `xmlrpc:"pid" json:"pid"`
	Server        string `json:"server"`
}

type Std struct {
	Log    string `xmlrpc:"log" json:"log"`
	Server string `json:"server"`
}

type ProcessStatus struct {
	Description string `xmlrpc:"description"`
	Group       string `xmlrpc:"group"`
	Name        string `xmlrpc:"name"`
	Status      int    `xmlrpc:"status"`
}

func (rpc *SupervisorRpc) GetClient() (*xmlrpc.Client, error) {
	return xmlrpc.NewClient(rpc.Url, nil)
}

func (rpc *SupervisorRpc) GetAPIVersion() (string, error) {
	ret := ""
	err := rpc.Client.Call("supervisor.getAPIVersion", nil, &ret)
	return ret, err
}

func (rpc *SupervisorRpc) Shutdown() (bool, error) {
	ret := false
	err := rpc.Client.Call("supervisor.shutdown", nil, &ret)
	return ret, err
}

func (rpc *SupervisorRpc) Restart() (bool, error) {
	ret := false
	err := rpc.Client.Call("supervisor.restart", nil, &ret)
	return ret, err
}

func (rpc *SupervisorRpc) GetAllProcessInfo() ([]ProcessInfo, error) {
	ret := make([]ProcessInfo, 0)
	err := rpc.Client.Call("supervisor.getAllProcessInfo", nil, &ret)
	return ret, err
}

func (rpc *SupervisorRpc) GetProcessInfo(group string, name string) (ProcessInfo, error) {
	ret := ProcessInfo{}
	params := fmt.Sprintf("%s:%s", group, name)
	err := rpc.Client.Call("supervisor.getProcessInfo", params, &ret)
	return ret, err
}

func (rpc *SupervisorRpc) StartProcess(group string, name string) (bool, error) {
	ret := false
	params := fmt.Sprintf("%s:%s", group, name)
	err := rpc.Client.Call("supervisor.startProcess", params, &ret)
	return ret, err
}

func (rpc *SupervisorRpc) StopProcess(group string, name string) (bool, error) {
	ret := false
	params := fmt.Sprintf("%s:%s", group, name)
	err := rpc.Client.Call("supervisor.startProcess", params, &ret)
	return ret, err
}

func (rpc *SupervisorRpc) StartAllProcess(wait bool) ([]interface{}, error) {
	params := []interface{}{wait}
	ret := make([]interface{}, 0)
	err := rpc.Client.Call("supervisor.startAllProcesses", params, &ret)
	return ret, err
}

func (rpc *SupervisorRpc) StartProcessGroup(group string) ([]ProcessStatus, error) {
	ret := make([]ProcessStatus, 0)
	params := []interface{}{group, true}
	err := rpc.Client.Call("supervisor.startProcessGroup", params, &ret)
	return ret, err
}

func (rpc *SupervisorRpc) StopAllProcesses() ([]ProcessStatus, error) {
	ret := make([]ProcessStatus, 0)
	params := []interface{}{true}
	err := rpc.Client.Call("supervisor.stopAllProcesses", params, &ret)
	return ret, err
}

func (rpc *SupervisorRpc) StopProcessGroup(group string) ([]ProcessStatus, error) {
	ret := make([]ProcessStatus, 0)
	params := []interface{}{group, true}
	err := rpc.Client.Call("supervisor.stopProcessGroup", params, &ret)
	return ret, err
}

func (rpc *SupervisorRpc) ReloadConfig() (bool, error) {
	ret := false
	err := rpc.Client.Call("supervisor.reloadConfig", nil, &ret)
	return ret, err
}

func (rpc *SupervisorRpc) GetStdErr(group string, name string) (Std, error) {
	ret := make([]interface{}, 0)
	params := []interface{}{fmt.Sprintf("%s:%s", group, name), 0, 5000}
	_ = rpc.Client.Call("supervisor.tailProcessStderrLog", params, &ret)
	log := ""
	if len(ret) > 0 && ret[0] != nil {
		log = ret[0].(string)
	}
	return Std{
		Log: log,
	}, nil
}

func (rpc *SupervisorRpc) GetStdOut(group string, name string) (Std, error) {
	ret := make([]interface{}, 0)
	params := []interface{}{fmt.Sprintf("%s:%s", group, name), 0, 5000}
	_ = rpc.Client.Call("supervisor.tailProcessStdoutLog", params, &ret)
	log := ""
	if len(ret) > 0 && ret[0] != nil {
		log = ret[0].(string)
	}
	return Std{
		Log: log,
	}, nil
}
