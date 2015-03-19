package kmgRpc

type RpcRequest struct {
	ApiName string
	InData  [][]byte
}
type RpcResponse struct {
	Error   string
	OutData [][]byte
}
