package bettermolang

import "fmt"

const testString string = `
func rotate(x, y, a){
    ca = cos(a);
    sa = sin(a);
    x = x*ca - y*sa;
    y = x*sa + y*ca;
}

func square(x){
    return x*x;
}

var str = 'stuff n () { things'
vec var p = (1,2)
arr tmp a = [1,2,2,7,5]
var px = a[3];
tmp py = true + 1;
py = square(py); #comment
px ++;
rotate(px,py,45.0);
/*
"(//filter this -> ))"
'(stuff)'
*/
`

const testString2 string = `
### This is a comment
+- / * ### This is also a comment
[{([{()}])}] ### This is also a comment
'This is a string!!!!!!!!! {ajuidhaiudhaiuhdiuahwiud} []() ******'
### Multi line string 
'THIS IS A MULTI LINE STRING
THIS IS A MULTI LINE STRING
THIS IS A MULTI LINE STRING'

' An escaped string!\' '
/* 
This is a multi line comment
Please fucking work 21312321 {ahduawhdu'dhwaiyudghuaiyghbdu '}
*/  
--++*
12.1 + 99

func functionName() {
    return 1 + 2
}

var b = 20.1
var a = b + functionName()
`

const testString3 string = `
/*
Test Comment Please Work You Fuckhead
*/
`

const testString4 string = `
'123\\\'321'
`

func RunTest() {
	scanner := NewScanner()
	tokens := scanner.ScanTokens(testString)
	for _, token := range tokens {
		println(fmt.Sprintf("String: %s, ID: %d", token.Value, token.TokenType))
	}
}
