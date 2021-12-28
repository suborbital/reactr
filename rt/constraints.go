//go:build go1.18

package rt

// Input is the constraint for types used as payloads for jobs
type Input interface{
	string | []byte | interface{}
}

// Output is the constraint for types used as returns from jobs
type Output interface{
	string | []byte | interface{}
}
