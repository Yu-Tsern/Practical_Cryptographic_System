package main

import(
    "crypto/rand"
    "fmt"
    "os"
    "math/big"
)
func main(){
    //testGenN()
    //testExEu()
    //testGenExp()
    //testFindInv()
    //testWrite()
    //testBoundedGen()
    //testMilRab()
    //testPwoMod()
    //testGenPrime()

    // read command line and parse
    agrv_s := os.Args

    // declare var
    var N, e, d, p, q, phi_N, p_1, q_1 big.Int
    big1 := big.NewInt(1)

    // generate p, q, N
    GenN(&p, &q, &N)
    p_1.Sub(&p, big1)
    q_1.Sub(&q, big1)
    phi_N.Mul(&p_1, &q_1)

    // generate exponent
    GenExp(&phi_N, &e)
    FindInv(&e, &d, &phi_N)

    // write key to file
    WritePK ( &N, &e, agrv_s[1])
    WriteSK ( &N, &d, &p, &q, agrv_s[2])

}

func FindInv(e *big.Int, d *big.Int, phi_N *big.Int) {
    a := big.NewInt(0)
    ExtendedEu(a, d, phi_N, e)
    d.Mod(d, phi_N)
}

func GenExp(phi_N *big.Int, e *big.Int) { // generate exponent
    big1 := big.NewInt(1)
    big2 := big.NewInt(2)
    BoundedGen(big1, phi_N, e)
    e.SetBit(e, 0, 1)
    var a, b big.Int
    for big1.Cmp(ExtendedEu(&a, &b, phi_N, e))!=0 || big2.Cmp(e)>-1 {
        e.Add(e, big2)
        if e.Cmp(phi_N)>-1 {
            e.Sub(e, phi_N)
        }
    }
}

func GenN(p *big.Int, q *big.Int, N *big.Int){
    GenPrime(p)
    GenPrime(q)
    N.Mul(p, q)
}

func GenPrime(p *big.Int){ // generate 1024 bit prime number 
    var lower, upper big.Int
    big1 := big.NewInt(1)
    big2 := big.NewInt(2)
    lower.Lsh(big1, 511)
    upper.Lsh(&lower, 1)
    BoundedGen(&lower, &upper, p)
    p.SetBit(p, 0, 1)
    p.SetBit(p, 510, 1)
    for !MilRab(p, 20) || p.BitLen()!=512{
        if p.BitLen() > 512 {
            p.Sub(p, &lower)
        }
        p.Add(big2, p)
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

func WritePK(N *big.Int, e *big.Int, filename string){
    f, err := os.Create(filename) // creating file
    if err != nil {
        fmt.Printf("error creating file: %v", err)
        return
    }
    defer f.Close()
    s := "Public key. (" + N.String() + "," + e.String() +")"
    _, err = f.WriteString(s) // writing number
    if err != nil {
        fmt.Printf("error writing string: %v", err)
    }
}

func WriteSK(N *big.Int, d *big.Int, p *big.Int, q *big.Int, filename string){
    f, err := os.Create(filename) // creating file
    if err != nil {
        fmt.Printf("error creating file: %v", err)
        return
    }
    defer f.Close()
    s := "Private key. (" + N.String() + "," + d.String() + "," + p.String() + "," + q.String() + ")"
    _, err = f.WriteString(s) // writing number
    if err != nil {
        fmt.Printf("error writing string: %v", err)
    }
}

func MilRab(n *big.Int, t int) bool { //true if it is probably a prime
    var n_1, a, r, y big.Int
    s := 0
    big1 := big.NewInt(1)
    big2 := big.NewInt(2)
    n_1.Sub(n, big1)
    r.Set(&n_1)
    for r.Bit(0) == 0 {
        r.Rsh(&r, 1)
        s += 1
    }
    for i:=0; i<t; i++ {
        BoundedGen(big2, &n_1, &a)
        PowMod(&y, &a, &r, n)
        if big1.Cmp(&y)!=0 && n_1.Cmp(&y)!=0 {
            for j:=1; j<s && n_1.Cmp(&y)!=0; j++ {
                y.Mul(&y, &y)
                y.Mod(&y, n)
                if big1.Cmp(&y)==0 {
                    //fmt.Println("s, r, a, y: ", s, " ", r.String(), " ", a.String(), " ", y.String())
                    return false
                }
            }
            if n_1.Cmp(&y)!=0 {
                //fmt.Println("s, r, a, y: ", s, " ", r.String(), " ", a.String(), " ", y.String())
                return false
            }
        }
    }
    return true
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





func testGenPrime(){
    var p big.Int
    GenPrime(&p)
    fmt.Println(p.ProbablyPrime(30))
}

func testPwoMod(){
    var y, z, N, e, d, p, q, phi_N, p_1, q_1 big.Int
    big1 := big.NewInt(1)
    GenN(&p, &q, &N)
    p_1.Sub(&p, big1)
    q_1.Sub(&q, big1)
    phi_N.Mul(&p_1, &q_1)
    GenExp(&phi_N, &e)
    FindInv(&e, &d, &phi_N)

    m := big.NewInt(53)
    y.Exp(m, &e, &N)
    PowMod(&z, &y, &d, &N)
    fmt.Println(y.String(), z.String())
}

func testMilRab(){
    tar := big.NewInt(103)
    res := MilRab(tar, 10)
    fmt.Println("My test: ", res)
    tar.Mul(tar, tar)
    fmt.Println("square tar: ", tar)
    var p, q, N big.Int
    GenN(&p, &q, &N)
    fmt.Println("probably prime : ", p.ProbablyPrime(30), q.ProbablyPrime(30), N.ProbablyPrime(30))
    fmt.Println("My Miller Rab  : ", MilRab(&p, 20), MilRab(&q, 20), MilRab(&N, 20))
}

func testBoundedGen(){
    l := big.NewInt(4)
    u := big.NewInt(6)
    res := big.NewInt(0)
    BoundedGen(l, u, res)
    fmt.Println(l.String(), u.String(), res.String())
}

func testWrite(){
    N := big.NewInt(100)
    e := big.NewInt(20)
    d := big.NewInt(2)
    p := big.NewInt(3)
    q := big.NewInt(4)
    WritePK(N, e, "PK.txt")
    WriteSK(N, d, p, q, "SK.txt")
}

func testExEu(){
    x := big.NewInt(56)
    y := big.NewInt(49)
    a := big.NewInt(0)
    b := big.NewInt(0)
    gcd := ExtendedEu(a, b, x, y)
    fmt.Println(gcd.String(), a.String(), b.String())
}

func testGenN(){
    var p, q, N big.Int
    GenN(&p, &q, &N)
    fmt.Println(p.String(), q.String(), N.String())
}

func testGenExp(){
    e := big.NewInt(0)
    phi_N := big.NewInt(58*12)
    GenExp(phi_N, e) // generate exponent
    fmt.Println(phi_N.String(), e.String())
}

func testFindInv(){
    e := big.NewInt(263)
    d := big.NewInt(263)
    phi_N := big.NewInt(58*12)
    FindInv(e, d, phi_N)
    fmt.Println(phi_N.String(), e.String(), d.String())
}
