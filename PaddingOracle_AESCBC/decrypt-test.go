package main

import(
/*
    "unicode/utf8"
    "unicode"
    "crypto/cipher"
    "log"
    "strconv"
    "bytes"
    "bufio"
    "flag"
*/
    "crypto/sha256"
    "crypto/aes"
    "encoding/hex"
    "fmt"
    "os"
    "io/ioutil"
)

func main(){


    // read command line and parse
    agrv_s := os.Args
    // Generate keys
    K := "4141414141414141414141414141414141414141414141414141414141414141"
    dst := make([]byte, hex.DecodedLen(len(K)))
    hex.Decode(dst, []byte(K))
    k_enc := dst[:len(dst)/2]
    k_mac := dst[len(dst)/2:]

    //  read file
    text, err := ioutil.ReadFile((agrv_s[2]))
    if err != nil {
        panic(err)
    }
    // decryptionn

    IV:=text[:16]
    text=text[16:]
    //fmt.Println(text)
    //text[len(text)-17]+=253
    output:=AES_DEC(text, IV, k_enc)

    //fmt.Println(output)
    if(!PADDING_CHECK(&output)){
        return
    }
    if(len(output)<33){
        fmt.Print("INVALID MAC")
        return
    }
    T_2:=output[len(output)-32:]
    output=output[:len(output)-32]
    if(!MAC_CHECK(k_mac, output, T_2)){
        return
    }
    fmt.Print("SUCCESS")
    //WriteText( agrv_s[7], &output)

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
            fmt.Print("INVALID MAC")
            return false
        }
    }
    return true
}

func PADDING_CHECK( M_2 *[]byte) bool{
    r:=((*M_2)[len(*M_2)-1])
    if(r==0 || r>16){
        fmt.Print("INVALID PADDING")
        return false
    }
    for i:=len(*M_2)-1; i>len(*M_2)-int(r)-1; i-- {
        if (*M_2)[i]!=r{
            fmt.Print("INVALID PADDING")
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


