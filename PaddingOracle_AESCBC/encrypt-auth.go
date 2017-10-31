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
    "crypto/aes"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "os"
    "io/ioutil"
)

func main(){


    // read command line and parse
    //reader := bufio.NewReader(os.Stdin)
    //agrv, _ := reader.ReadBytes('\n')
    //agrv_s := bytes.Split(agrv, []byte(" "))
    agrv_s := os.Args
    // parse keys
    dst := make([]byte, hex.DecodedLen(len(agrv_s[3])))
    hex.Decode(dst, []byte(agrv_s[3]))
    k_enc := dst[:len(dst)/2]
    k_mac := dst[len(dst)/2:]
    //  read file
    text, err := ioutil.ReadFile((agrv_s[5]))
    if err != nil {
        panic(err)
    }
    // encryption & decryptionn
    if(agrv_s[1]=="encrypt"){
        T:=HMAC(k_mac, text)
        M_2:=PADDING(text, T)
        IV:=GEN_IV(16)
        //fmt.Println(IV)
        output:=AES_ENC(M_2, IV, k_enc)
        Append(&IV, &output)
        WriteText(agrv_s[7], &IV)

    }else if((agrv_s[1])==("decrypt")){
        IV:=text[:16]
        text=text[16:]
        output:=AES_DEC(text, IV, k_enc)
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
        WriteText( agrv_s[7], &output)
    }
    //fmt.Println("plain text : ", plaintext)
    //fmt.Println(k_enc, k_mac)
    //fmt.Println("M_2 = ", M_2)
    // set IV
    //fmt.Println("IV = ", IV)
    //fmt.Println("encrypted ", output)

    //fmt.Println("decrypted ", output)
    //fmt.Println(ck)
    //fmt.Println("depad ", output)

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

