package main

import(
    "strings"
    "io/ioutil"
    "log"
    "fmt"
    "os"
    "math/big"
)
func main(){
    //testPwoMod()
    //testReadParse()

    // read command line and parse
    agrv_s := os.Args

    // declare var
    var N, e, m, c big.Int

    // Read and Parse
    text := Readfile(agrv_s[1])
    ParseKey(&text, &N, &e)
    m.SetString(agrv_s[2], 10)

    // raise m^e
    PowMod(&c, &m, &e, &N)

    // stdout
    fmt.Print(c.String())
}


func PowMod(z *big.Int, x *big.Int, y *big.Int, n *big.Int){ // z = x^y mod n
    var p, r big.Int
    big1 := big.NewInt(1)
    p.Set(y)
    r.Set(x)
    z.Set(big1)
    for p.BitLen() > 0 {
        if p.Bit(0) != 0 {
            z.Mul(z, &r)
            z.Mod(z, n)
        }
        p.Rsh(&p, 1)
        r.Mul(&r, &r)
        r.Mod(&r, n)
    }
}

func Readfile(filename string) []byte {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    b, err := ioutil.ReadAll(file)
    return b
}

func ParseKey(b *[]byte, N *big.Int, e *big.Int){
    s := string(*b)
    s_s := strings.Split(s, ",")
    s_s0 := strings.Split(s_s[0], "(")
    s_sn := strings.Split(s_s[len(s_s)-1], ")")
    N.SetString(s_s0[1], 10)
    e.SetString(s_sn[0], 10)
}



func testReadParse(){
    filename := "pk.txt"
    text := Readfile(filename)
    fmt.Println(string(text))
    var N, e big.Int
    ParseKey(&text, &N, &e)
    fmt.Println(N.String(), e.String())
}
