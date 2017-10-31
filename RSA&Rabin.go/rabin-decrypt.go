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
    //testExEu()
    //testPwoMod()
    //testReadParse()
    //testSqrtN()
    //testRep()

    // read command line and parse
    agrv_s := os.Args

    // declare var
    var N, p, q, m, c big.Int

    // Read and Parse
    text := Readfile(agrv_s[1])
    ParseKey(&text, &N, &p, &q)
    c.SetString(agrv_s[2], 10)

    // decrypt
    suc := SqrtN(&c, &p, &q, &N, &m)
    if suc {
        fmt.Print(m.String())
    } else {
        fmt.Println("Fail")
    }

}


func ExtendedEu(a *big.Int, b *big.Int, x *big.Int, y *big.Int) *big.Int { // assume x > y
    // A = B*Q + R
    var X, Y, Q, R, a_p, b_p, a_n, b_n big.Int
    big0 := big.NewInt(0)
    big1 := big.NewInt(1)
    X.Set(x)
    Y.Set(y)
    Q.QuoRem(x, y, &R)
    a.Set(big1)
    b.Neg(&Q)
    a_p.Set(big0)
    b_p.Set(big1)
    for R.Cmp(big0)!=0 {
        X.Set(&Y)
        Y.Set(&R)
        Q.QuoRem(&X, &Y, &R)
        a_n.Add(&a_p, a_n.Neg(a_n.Mul(&Q, a)))
        a_p.Set(a)
        a.Set(&a_n)
        b_n.Add(&b_p, b_n.Neg(b_n.Mul(&Q, b)))
        b_p.Set(b)
        b.Set(&b_n)
    }
    a.Set(&a_p)
    b.Set(&b_p)
    return &Y
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

func ParseKey(b *[]byte, N *big.Int, p *big.Int, q *big.Int){
    s := string(*b)
    s_s := strings.Split(s, ",")
    s_s0 := strings.Split(s_s[0], "(")
    s_sn := strings.Split(s_s[len(s_s)-1], ")")
    N.SetString(s_s0[1], 10)
    p.SetString(s_s[1], 10)
    q.SetString(s_sn[0], 10)
}

func SqrtN(c *big.Int, p *big.Int, q *big.Int, N *big.Int, m *big.Int) bool {
    var root [4]big.Int
    var a, b, r, s, x, y big.Int
    big1 := big.NewInt(1)
    big4 := big.NewInt(4)
    ExtendedEu(&a, &b, p, q)
    r.Add(p, big1)
    r.Quo(&r, big4)
    PowMod(&r, c, &r, p)
    r.Mul(&r, &b)
    r.Mod(&r, N)
    r.Mul(&r, q)
    r.Mod(&r, N)
    s.Add(q, big1)
    s.Quo(&s, big4)
    PowMod(&s, c, &s, q)
    s.Mul(&s, &a)
    s.Mod(&s, N)
    s.Mul(&s, p)
    s.Mod(&s, N)
    x.Add(&s, &r)
    x.Mod(&x, N)
    y.Sub(&s, &r)
    y.Mod(&y, N)
    root[0]=x
    root[1]=*(x.Neg(&x))
    root[2]=y
    root[3]=*(y.Neg(&y))
    for i:=0; i<len(root); i++ {
        root[i].Mod(&(root[i]), N)
        //fmt.Println("root in Sqrt: ", (root[i]).String())
        if CheckRep(&root[i]) {
            m.Set(&root[i])
            return true
        }
    }
    return false
}

func CheckRep(root *big.Int) bool {
    s := root.String()
    //fmt.Println("root in CheckRep: ", s)
    l := len(s)
    if l < 20 {
        //fmt.Println("l < 20")
        l = l/2
        for i:=0; i<l; i++ {
            if s[i] != s[l+i] {
                return false
            }
        }
        s = s[:l]
        root.SetString(s, 10)
    } else {
        //fmt.Println("l >= 20")
        for i:=0; i<10; i++ {
            if s[i] != s[l+i-10] {
                //fmt.Println("difference: ",i, " ",  s[i], " ", s[l+i-10])
                return false
            }
        }
        s = s[:l-10]
        root.SetString(s, 10)
    }
    return true
}


func testRep(){
    c := big.NewInt(9)
    p := big.NewInt(19)
    c.SetString(c.String() + c.String(), 10)
    fmt.Println("c: ", CheckRep(c))
    fmt.Println("p: ", CheckRep(p))
}

func testSqrtN(){
    c := big.NewInt(9) // m=1, m'=11, c=1001
    p := big.NewInt(19)
    q := big.NewInt(11)
    N := big.NewInt(209)
    m := big.NewInt(209)
    SqrtN(c, p, q, N, m)
}


func testReadParse(){
    filename := "Rabin_sk.txt"
    text := Readfile(filename)
    fmt.Println(string(text))
    var N, p, q big.Int
    ParseKey(&text, &N, &p, &q)
    fmt.Println(N.String())
    fmt.Println(p.String())
    fmt.Println(q.String())
}

func testPwoMod(){
    z := big.NewInt(0)
    y := big.NewInt(31415926)
    m := big.NewInt(271828)
    e := big.NewInt(746538)
    N := big.NewInt(5783648)

    y.Exp(m, e, N)
    PowMod(z, m, e, N)
    fmt.Println(y.String(), z.String())
}

func testExEu(){
    x := big.NewInt(56)
    y := big.NewInt(49)
    a := big.NewInt(0)
    b := big.NewInt(0)
    gcd := ExtendedEu(a, b, x, y)
    fmt.Println(gcd.String(), a.String(), b.String())
}

