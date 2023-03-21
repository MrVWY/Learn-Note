package redis_implementation_mq

type Conf struct {
	A *A
}

type A struct {
	addr string
}

func (A *A) GetAddr() string {
	return A.addr
}

func Config() *Conf {
	return &Conf{}
}
