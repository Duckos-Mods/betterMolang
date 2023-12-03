package bettermolang

const testString string = `
func rotate(inout x, inout y, a){
    ca = cos(a);
    sa = sin(a);
    x = x*ca - y*sa;
    y = x*sa + y*ca;
}

func square(x){
    return x*x;
}

var str = "stuff n () { things'
vec var p = (1,2)
arr tmp a = [1,2,2,7,5]
var px = a[3];
tmp py = true + 1;
py = square(py); //comment
px ++;
rotate(px,py,45.0);

/*
"(//filter this -> ))"
'(stuff)'
*/

//comment
`

const testString2 string = `
### This is a comment
+- / * ### This is also a comment
[{([{()}])}] ### This is also a comment
`

func RunTest() {
	scanner := NewScanner()
	tokens := scanner.ScanTokens(testString2)
	for _, token := range tokens {
		println(token.Value)
	}
}
