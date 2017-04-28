package demo

// file generated by
// github.com.mh-vbon/channeler
// do not edit

type MyTomate struct {
	embed Tomate
	ops   chan func()
}

func NewMyTomate() *MyTomate {
	ret := &MyTomate{}
	ret.loop()
	return ret
}
func (t *MyTomate) loop() {
	for {
		select {
		case op := <-t.ops:
			op()
		}
	}
}
func (t *MyTomate) Hello() {
	t.ops <- func() {
		t.embed.Hello()
	}
}
func (t *MyTomate) Good() {
	t.ops <- func() {
		t.embed.Good()
	}
}
func (t *MyTomate) Name(it string) string {
	var retVar0 string
	t.ops <- func() {
		retVar0 = t.embed.Name(it)
	}
	return retVar0
}
