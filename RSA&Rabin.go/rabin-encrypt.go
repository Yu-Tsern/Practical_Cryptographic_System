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
    //testReadParse()

    // read command line and parse
    agrv_s := os.Args

    // declare var
    var N, m, c big.Int

    // Read and Parse
    text := Readfile(agrv_s[1])
    ParseKey(&text, &N)

    // Pad message to make decryption unique
    if len(agrv_s[2]) < 10 {
        agrv_s[2] = agrv_s[2] + agrv_s[2]
    } else {
        agrv_s[2] = agrv_s[2] + agrv_s[2][:10]
    }
    m.SetString(agrv_s[2], 10)

    // raise m^2 mod N
    c.Mul(&m, &m)
    c.Mod(&c, &N)

    // stdout
    fmt.Print(c.String())

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




func testReadParse(){
    filename := "Rabin_pk.txt"
    text := Readfile(filename)
    fmt.Println(string(text))
    var N big.Int
    ParseKey(&text, &N)
    fmt.Println(N.String())
}

