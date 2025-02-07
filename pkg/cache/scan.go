package cache

type Scanner interface {
	Scan() bool    // move to next, return false if reach the end
	Key() string   // return the key of current entry
	Value() []byte // return the value of current entry
	Close()        // close and release the scanner
}
