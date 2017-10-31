package main

import(
    "strings"
    "io/ioutil"
    "log"
    "fmt"
    "os"
    "math/big"
    "os/exec"
    "bytes"
    "crypto/rand"
)
func main(){
    //testExEu()
    //testReadParse()
    //testSqrtN()
    //testExe()

    // read command line and parse
    agrv_s := os.Args

    // declare var
    var N, p, q big.Int

    // Read and Parse
    text := Readfile(agrv_s[1])
    ParseKey(&text, &N)

    // Try until m_1 != m_2
    FindFactor(agrv_s[1], &N, &p)
    q.Quo(&N, &p)
    fmt.Print(p.String(), ",", q.String())
/*
    var t1, t2 big.Int
    t1.Mod(&N, &p)
    t2.Mod(&N, &q)
    fmt.Println("\n", t1.String(), "\n", t2.String())
*/
}

func BoundedGen( l *big.Int, u *big.Int, res *big.Int){ // [ lower, upper )
    var m big.Int
    m.Sub(u, l)
    r, err := rand.Int(rand.Reader, &m)
    if err != nil {
        fmt.Println(err)
    }
    res.Add(l, r)
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

func Readfile(filename string) []byte {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    b, err := ioutil.ReadAll(file)
    return b
}

func ParseKey(b *[]byte, N *big.Int){
    s := string(*b)
    s_s0 := strings.Split(s, "(")
    s_sn := strings.Split(s_s0[1], ")")
    N.SetString(s_sn[0], 10)
}

func ExCrack(pkfile string, c_text string) string {
    name := "./rabin-crack"
    cmd := exec.Command(name, pkfile, c_text)
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
    outStr := string(stdout.Bytes())
    return outStr
}

func FindFactor(pkfile string, N *big.Int, res *big.Int){
    var c, m_1p, m_1n, m_2, a, b big.Int
    big3 := big.NewInt(3)
    for true {
        BoundedGen(big3, N, &m_1p)
        m_1n.Neg(&m_1p)
        m_1n.Mod(&m_1n, N)
        c.Mul(&m_1p, &m_1p)
        c.Mod(&c, N)
        s := ExCrack(pkfile, c.String())
        m_2.SetString(s, 10)
        if m_2.Cmp(&m_1p) != 0 && m_2.Cmp(&m_1n) != 0 {
            m_2.Sub(&m_2, &m_1p)
            m_2.Mod(&m_2, N)
            p := ExtendedEu(&a, &b, &m_2, N)
            res.Set(p)
            return
        }
    }
}

func testExe(){
    pkfile := "Rabin_pk.txt"
    c_text := "9"
    p := ExCrack(pkfile, c_text)
    fmt.Println(p)
    //p.Mul(&p, &p)
    //p.Mod(&p, &N)
}

func testReadParse(){
    filename := "Rabin_sk.txt"
    text := Readfile(filename)
    fmt.Println(string(text))
    var N big.Int
    ParseKey(&text, &N)
    fmt.Println(N.String())
}

func testExEu(){
    x := big.NewInt(56)
    y := big.NewInt(49)
    a := big.NewInt(0)
    b := big.NewInt(0)
    gcd := ExtendedEu(a, b, x, y)
    fmt.Println(gcd.String(), a.String(), b.String())
}

