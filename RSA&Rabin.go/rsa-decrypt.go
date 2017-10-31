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
    var N, d, m, c, p, q big.Int

    // Read and Parse
    text := Readfile(agrv_s[1])
    ParseKey(&text, &N, &d, &p, &q)
    c.SetString(agrv_s[2], 10)

    // raise m^e
    PowMod(&m, &c, &d, &N)

    // stdout
    fmt.Print(m.String())

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

func ParseKey(b *[]byte, N *big.Int, d *big.Int, p *big.Int, q *big.Int){
    s := string(*b)
    s_s := strings.Split(s, ",")
    s_s0 := strings.Split(s_s[0], "(")
    s_sn := strings.Split(s_s[len(s_s)-1], ")")
    N.SetString(s_s0[1], 10)
    d.SetString(s_s[1], 10)
    p.SetString(s_s[2], 10)
    q.SetString(s_sn[0], 10)
}



func testReadParse(){
    filename := "sk.txt"
    text := Readfile(filename)
    fmt.Println(string(text))
    var N, d, p, q big.Int
    ParseKey(&text, &N, &d, &p, &q)
    fmt.Println(N.String(), d.String(), p.String(), q.String())
}
