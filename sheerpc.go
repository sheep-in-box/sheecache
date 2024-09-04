package sheecache

type Args struct {
	Group string
	Key   string
}

type Reply struct {
	Value []byte
}
