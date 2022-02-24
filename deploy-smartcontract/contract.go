package main

type contract struct {
	Deploy struct {
		VM struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"VM:-"`
		Main1 struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"main:1"`
		Ropsten3 struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"ropsten:3"`
		Rinkeby4 struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"rinkeby:4"`
		Kovan42 struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"kovan:42"`
		Görli5 struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"görli:5"`
		Custom struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			AutoDeployLib bool `json:"autoDeployLib"`
		} `json:"Custom"`
	} `json:"deploy"`
	Data struct {
		Bytecode struct {
			LinkReferences struct {
			} `json:"linkReferences"`
			Object    string `json:"object"`
			Opcodes   string `json:"opcodes"`
			SourceMap string `json:"sourceMap"`
		} `json:"bytecode"`
		DeployedBytecode struct {
			ImmutableReferences struct {
			} `json:"immutableReferences"`
			LinkReferences struct {
			} `json:"linkReferences"`
			Object    string `json:"object"`
			Opcodes   string `json:"opcodes"`
			SourceMap string `json:"sourceMap"`
		} `json:"deployedBytecode"`
		GasEstimates struct {
			Creation struct {
				CodeDepositCost string `json:"codeDepositCost"`
				ExecutionCost   string `json:"executionCost"`
				TotalCost       string `json:"totalCost"`
			} `json:"creation"`
			External struct {
				GetMessage       string `json:"get_message()"`
				SetMessageString string `json:"set_message(string)"`
			} `json:"external"`
		} `json:"gasEstimates"`
		MethodIdentifiers struct {
			GetMessage       string `json:"get_message()"`
			SetMessageString string `json:"set_message(string)"`
		} `json:"methodIdentifiers"`
	} `json:"data"`
	Abi []struct {
		Inputs []struct {
			InternalType string `json:"internalType"`
			Name         string `json:"name"`
			Type         string `json:"type"`
		} `json:"inputs"`
		StateMutability string `json:"stateMutability"`
		Type            string `json:"type"`
		Name            string `json:"name,omitempty"`
		Outputs         []struct {
			InternalType string `json:"internalType"`
			Name         string `json:"name"`
			Type         string `json:"type"`
		} `json:"outputs,omitempty"`
	} `json:"abi"`
}
