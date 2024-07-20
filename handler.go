package naspad

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc
