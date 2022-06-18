package ziface

type IRequest interface {
	GetConnection() IConnection

	GetDate() []byte
}
