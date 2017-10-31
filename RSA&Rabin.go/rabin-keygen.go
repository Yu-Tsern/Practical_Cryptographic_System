package main

import(
    "crypto/rand"
    "fmt"
    "os"
    "math/big"
)
func main(){
    //testGenN()
    //testWrite()
    //testBoundedGen()
    //testMilRab()
    //testPwoMod()
    //testGenPrime()

    // read command line and parse
    agrv_s := os.Args

    // declare var
    var N, p, q big.Int

    // generate p, q, N
    GenN(&p, &q, &N)

    // write key to file
    WritePK ( &N, agrv_s[1])
    WriteSK ( &N, &p, &q, agrv_s[2])
}

func GenN(p *big.Int, q *big.Int, N *big.Int){
    GenPrime(p)
    GenPrime(q)
    N.Mul(p, q)
}

func GenPrime(p *big.Int){ // generate 512 bits prime number that is 4k+3
    var lower, upper big.Int
    big1 := big.NewInt(1)
    big4 := big.NewInt(4)
    lower.Lsh(big1, 511)
    upper.Lsh(&lower, 1)
    BoundedGen(&lower, &upper, p)
    p.SetBit(p, 0, 1)
    p.SetBit(p, 1, 1)
    p.SetBit(p, 510, 1)
    for !MilRab(p, 20) || p.BitLen()!=512 {
        if p.BitLen() > 512 {
            p.Sub(p, &lower)
        }
        p.Add(big4, p)
    }
}

func WritePK(N *big.Int, filename string){
    f, err := os.Create(filename) // creating file
    if err != nil {
        fmt.Printf("error creating file: %v", err)
        return
    }
    defer f.Close()
    s := "Public key. (" + N.String() +")"
    _, err = f.WriteString(s) // writing number
    if err != nil {
        fmt.Printf("error writing string: %v", err)
    }
}

func WriteSK(N *big.Int, p *big.Int, q *big.Int, filename string){
    f, err := os.Create(filename) // creating file
    if err != nil {
        fmt.Printf("error creating file: %v", err)
        return
    }
    defer f.Close()
    s := "Private key. (" + N.String() + "," + p.String() + "," + q.String() + ")"
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
    p := big.NewInt(3)
    q := big.NewInt(4)
    WritePK(N, "Rabin_PK.txt")
    WriteSK(N, p, q, "Rabin_SK.txt")
}

func testGenN(){
    var p, q, N big.Int
    GenN(&p, &q, &N)
    fmt.Println(p.String(), q.String(), N.String())
}
