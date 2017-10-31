package main

import(
/*
    "unicode/utf8"
    "unicode"
    "crypto/cipher"
    "strconv"
    "bufio"
    "flag"
    "encoding/hex"
*/
    "crypto/aes"
    "crypto/rand"
    "crypto/sha256"
    "fmt"
    "os"
    "io/ioutil"
    "log"
    "bytes"
    "os/exec"
)
func main(){


    // read command line and parse
    agrv_s := os.Args
    text, err := ioutil.ReadFile((agrv_s[2]))
    if err != nil {
        panic(err)
    }
    output:=make([]byte, len(text))
    //var output []byte
    for len(text)>16 {
        temp:=lastblockattack(text)
        //fmt.Println(temp)
        for i:=0; i<16; i++ {
            output[len(text)-16+i]=temp[i]
        }
        text=text[:len(text)-16]
    }
    //fmt.Println(output)
    if(!PADDING_CHECK(&output)){
        return
    }
    if(len(output)<33){
        fmt.Print("INVALID MAC")
        return
    }
    output=output[:len(output)-32]
    output=output[16:]
    fmt.Print(string(output))
}

func lastblockattack(text []byte) []byte{
    //fmt.Println("text= ", text)
    var duplic=make([]byte, len(text))
    var output=make([]byte, 16) // only 16 bytes in a block
    var aesout=make([]byte, 16) // only 16 bytes in a block
    var counter byte
    copy(duplic, text)
    filename:="blockattack.txt"
    duplic[len(duplic)-17]++
    WriteText(filename, &duplic)
    duplic[len(duplic)-17]--
    for i:=1; i<17; i++ { // the first 16 byte is IV
        //if(i==3){ break}
        counter=0x01
        //a:=test(filename)
        //fmt.Println(len(a))
        //fmt.Println(test(filename)=="INVALID PADDING\n")
        for test(filename)=="INVALID PADDING" {
            counter++
            inc1byte(filename, len(text)-i-16)
        }
        //fmt.Println("counter ", counter)
        //fmt.Println("byte(i) ", byte(i))
        //fmt.Println("text[len(text)-i-16] ", text[len(text)-i-16])
        //fmt.Println("text[len(text)-i-16]+counter ", text[len(text)-i-16]+counter)
        aesout[16-i]=byte(i)^(text[len(text)-i-16]+counter)
        output[16-i]=aesout[16-i]^(text[len(text)-i-16])
        //fmt.Println(output[16-i])
        for j:=1; j<i+1; j++ {
            duplic[len(duplic)-16-j]=aesout[16-j]^byte(i+1)
        }
        if(i!=16){
            duplic[len(duplic)-i-17]++
        }
        //fmt.Println("duplic= ", duplic)
        //fmt.Println("text= ", text)
        WriteText(filename, &duplic)
    }
    return output
}


func inc1byte(filename string, index int){
    text, err := ioutil.ReadFile((filename))
    if err != nil {
        panic(err)
    }
    text[index]++
    WriteText(filename, &text)
}

func test(filename string) string{
    name := "./decrypt-test"
    cmd := exec.Command(name, "-i", filename)
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

func Append( a *[]byte, b *[]byte){
    for i:=0; i<len(*b); i++ {
        (*a)=append((*a), (*b)[i])
    }
}

func WriteText ( filename string, text *[]byte) {  // writing plain text is the same as writing cipher text
    err:=ioutil.WriteFile( filename, *text, 0644)
    if err!=nil {
        panic(err)
    }
}

func MAC_CHECK(k_mac []byte, M []byte, T_2 []byte) bool {
    t:=HMAC(k_mac, M)
    for i:=0; i<len(t); i++ {
        if t[i]!=T_2[i] {
            fmt.Println("INVALID MAC")
            return false
        }
    }
    return true
}

func PADDING_CHECK( M_2 *[]byte) bool{
    r:=((*M_2)[len(*M_2)-1])
    for i:=int(r); i>len(*M_2)-int(r)-1-1; i--{
        if (*M_2)[i]!=r{
            fmt.Println("INVALID PADDING")
            return false
        }
    }
    (*M_2)=(*M_2)[:len(*M_2)-int(r)]
    return true
}

func AES_DEC( C []byte, IV []byte, k_enc []byte) []byte {
    AES, err := aes.NewCipher(k_enc)
    if err != nil {
        panic(err)
    }
    i_block:=make([]byte, 16)
    o_block:=make([]byte, 16)
    var output []byte
    for len(C)>0 {
        i_block=C[:16]
        AES.Decrypt(o_block, i_block)
        for i:=0; i<len(IV); i++ {
            o_block[i]=IV[i]^o_block[i]
        }
        IV=i_block
        for i:=0; i<16; i++ {
            output=append(output, o_block[i])
        }
        if len(C)>16{
            C=C[16:]
        }else{
            C=C[:0]
        }
    }
    return output
}

func AES_ENC(M_2 []byte, IV []byte, k_enc []byte) []byte {
    AES, err := aes.NewCipher(k_enc)
    if err != nil {
        panic(err)
    }
    i_block:=make([]byte, 16)
    o_block:=make([]byte, 16)
    var output []byte
    for i:=0; i<len(IV); i++ {
        i_block[i]=IV[i]^M_2[i]
    }
    for len(M_2)>0 {
        AES.Encrypt(o_block, i_block)
        for i:=0; i<16; i++ {
            output=append(output, o_block[i])
        }
        if len(M_2)>16{
            M_2=M_2[16:]
            for i:=0; i<16; i++ {
                i_block[i]=o_block[i]^M_2[i]
            }
        }else{
            M_2=M_2[:0]
        }
    }
    return output
}

func PADDING(plaintext []byte, T [32]byte) []byte {
    // compute M_1
    M_1:=append(plaintext,(T[0]))
    for i:=1; i<len(T); i++{
        M_1=append(M_1,(T[i]))
    }

    // compute M_2
    var n byte
    var PS []byte
    n=byte(len(M_1)%16)
    if n!=0 {
        PS=make([]byte, int(16-n))
        for i:=0; i<len(PS); i++ {
            PS[i]=byte((0x10-n))
        }
    }else{
        PS=make([]byte, 16)
        for i:=0; i<16; i++ {
            PS[i]=0x10
        }
    }
    M_2:=append(M_1, PS[0])
    for i:=1; i<len(PS); i++ {
        M_2=append(M_2, PS[i])
    }
    return M_2
}

func HMAC( k_mac []byte, plaintext []byte) [32]byte{
    // append key
    for len(k_mac)<64 {
        k_mac=append(k_mac, 0x00)
    }

    // create ipad and opad 
    ipad:=make([]byte, 64)
    for i:=0; i<64; i++ {
        ipad[i]=0x36
    }
    opad:=make([]byte, 64)
    for i:=0; i<64; i++ {
        opad[i]=0x5C
    }

    // K XOR with ipad
    temp:=make([]byte, 64)
    for i:=0; i<64; i++ {
        temp[i]=ipad[i]^k_mac[i]
    }
    // append text
    for i:=0; i<len(plaintext); i++ {
        temp=append(temp, plaintext[i])
    }

    // first hash
    first_hash := sha256.Sum256(temp)

    // K XOR with opad
    temp2:=make([]byte, 64)
    for i:=0; i<64; i++ {
        temp2[i]=opad[i]^k_mac[i]
    }

    // append first hash
    for i:=0; i<len(first_hash); i++ {
        temp2=append(temp2, first_hash[i])
    }

    // second hash
    second_hash := sha256.Sum256(temp2)
    //fmt.Println(second_hash)
    return second_hash
}

func GEN_IV(c int) []byte {
    b := make([]byte, c)
    rand.Read(b)
    return b
}


/*
    block, err := aes.NewCipher(k_enc)
    if err != nil {
        panic(err)
    }
    // The IV needs to be unique, but not secure. Therefore it's common to
    // include it at the beginning of the ciphertext.
    ciphertext := make([]byte, aes.BlockSize+len(M_2))
    //iv := ciphertext[:aes.BlockSize]
    mode := cipher.NewCBCEncrypter(block, IV)
    mode.CryptBlocks(ciphertext[aes.BlockSize:], M_2)
    // It's important to remember that ciphertexts must be authenticated
    // (i.e. by using crypto/hmac) as well as being encrypted in order to
    // be secure.
    fmt.Printf("%x\n", ciphertext)
*/

    // 
/*
    H := 6a09e667
    ipad = the byte 0x36 repeated B times
    opad = the byte 0x5C repeated B times

    H(K XOR opad, H(K XOR ipad, text))
*/

