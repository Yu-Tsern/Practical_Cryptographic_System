package main

import(
    "fmt"
    "io/ioutil"
    "bytes"
    "unicode/utf8"
    "unicode"
    "strconv"
    "os"
//    "time"
)
func main(){

//  input of command line

    agrv := os.Args
    //var agrv[3] string
    //fmt.Scanln(&agrv[0], &agrv[1], &agrv[2])

        var words[4][3] rune
        var words_f[4] int
        var key []byte
        var separatedtext [][]byte
        const lowerbound=11
        const upperbound=21
        text:=ReadCiphertext(&agrv[1])
        Find3Letters (&text, &words_f, &words)
        keylengths:=make( []int, upperbound, upperbound)
        FactorizeDistances( &words_f, lowerbound, upperbound, &words, &text, &keylengths)
        pkl:=make([]int, 2, 2)
        pkl[0] ,pkl[1]=GuessKeylength(&keylengths)
        key_l_final:=ICVariance(&pkl, &text)
        BuildSeparated( &separatedtext, &text, key_l_final)
        RecoverKey( &separatedtext, &key)
        CheckRepeated(&key)
        fmt.Print(len(key))

}

func ReadPlaintext( argument *string) []byte{
    //read the file
    content, err := ioutil.ReadFile(*argument)
    if err != nil {
        panic(err)
    }
    //convert to upper case
    content=bytes.ToUpper(content)
    //strip out characters that are not letters
    content=bytes.Map(LetterOnly, content)
    return content
}

func ReadCiphertext( argument *string) []byte{
    //read the file
    content, err := ioutil.ReadFile(*argument)
    if err != nil {
        panic(err)
    }
    return content
}

func LetterOnly(r rune) rune {
    if !unicode.IsLetter(r) {
        return -1
    }
    return r
}

func KeyExpansion( length int, argument *string) []byte {
    var long_key []byte
    var short_key []byte
    short_key=[]byte(*argument)
    if length/len(short_key)>1{
        long_key=bytes.Repeat(short_key, length/len(short_key))
    }
    if length%len(short_key)!=0{
        long_key=append(long_key, short_key[:length%len(short_key)]...)
    }
    return long_key
}


func Encrypt( key *[]byte, text *[]byte) {
    var r rune
    buf:=make([]byte, 1)
    for i:=0; i<len(*key); i++ {
        r=rune((*key)[i])+rune((*text)[i])-65
        if r>90 {
            r=r-26
        }
        utf8.EncodeRune(buf, r)
        (*text)[i]=buf[0]
    }
}

func Decrypt( key *[]byte, text *[]byte) {
    var r rune
    buf:=make([]byte, 1)
    for i:=0; i<len(*key); i++ {
        r=-rune((*key)[i])+rune((*text)[i])+65
        if r<65 {
            r=r+26
        }
        utf8.EncodeRune(buf, r)
        (*text)[i]=buf[0]
    }
}

func WriteText ( filename string, text *[]byte) {  // writing plain text is the same as writing cipher text
    err:=ioutil.WriteFile( filename, *text, 0644)
    if err!=nil {
    panic(err)
    }
}

func Find3Letters ( text *[]byte, max *[4]int, max_index *[4][3]rune) (){
    var n[26][26][26] int
    for i:=0; i<len(*text)-2; i++ {
        n[rune((*text)[i])-65][rune((*text)[i+1])-65][rune((*text)[i+2])-65]++
    }
    for i:=0; i<26; i++ {
        for j:=0; j<26; j++ {
            for k:=0; k<26; k++ {
                if n[i][j][k]>(*max)[3] {
                    if n[i][j][k]>(*max)[2]{
                        if n[i][j][k]>(*max)[1] {
                            if n[i][j][k]>(*max)[0] {
                                (*max)[3]=(*max)[2]
                                (*max)[2]=(*max)[1]
                                (*max)[1]=(*max)[0]
                                (*max)[0]=n[i][j][k]
                                (*max_index)[3][0]=(*max_index)[2][0]
                                (*max_index)[3][1]=(*max_index)[2][1]
                                (*max_index)[3][2]=(*max_index)[2][2]
                                (*max_index)[2][0]=(*max_index)[1][0]
                                (*max_index)[2][1]=(*max_index)[1][1]
                                (*max_index)[2][2]=(*max_index)[1][2]
                                (*max_index)[1][0]=(*max_index)[0][0]
                                (*max_index)[1][1]=(*max_index)[0][1]
                                (*max_index)[1][2]=(*max_index)[0][2]
                                (*max_index)[0][0]=rune(i+65)
                                (*max_index)[0][1]=rune(j+65)
                                (*max_index)[0][2]=rune(k+65)
                            } else {
                                (*max)[3]=(*max)[2]
                                (*max)[2]=(*max)[1]
                                (*max)[1]=n[i][j][k]
                                (*max_index)[3][0]=(*max_index)[2][0]
                                (*max_index)[3][1]=(*max_index)[2][1]
                                (*max_index)[3][2]=(*max_index)[2][2]
                                (*max_index)[2][0]=(*max_index)[1][0]
                                (*max_index)[2][1]=(*max_index)[1][1]
                                (*max_index)[2][2]=(*max_index)[1][2]
                                (*max_index)[1][0]=rune(i+65)
                                (*max_index)[1][1]=rune(j+65)
                                (*max_index)[1][2]=rune(k+65)
                            }
                        } else {
                            (*max)[3]=(*max)[2]
                            (*max)[2]=n[i][j][k]
                            (*max_index)[3][0]=(*max_index)[2][0]
                            (*max_index)[3][1]=(*max_index)[2][1]
                            (*max_index)[3][2]=(*max_index)[2][2]
                            (*max_index)[2][0]=rune(i+65)
                            (*max_index)[2][1]=rune(j+65)
                            (*max_index)[2][2]=rune(k+65)
                        }
                    } else {
                        (*max)[3]=n[i][j][k]
                        (*max_index)[3][0]=rune(i+65)
                        (*max_index)[3][1]=rune(j+65)
                        (*max_index)[3][2]=rune(k+65)
                    }
                }
            }
        }
    }
    //fmt.Println((*max), (*max_index))
}

func FactorizeDistances( words_f *[4]int, lower int, upper int, words *[4][3]rune, text*[]byte, keylength *[]int) {
    var distances [][]int
    var index int
    var cursor int
    distances=make([][]int, 4, 4)
    for k:=0; k<4; k++ {
        index=0
        cursor=0
        for i:=0; i<len((*text))-2; i++ {
            if rune((*text)[i])==(*words)[k][0] && rune((*text)[i+1])==(*words)[k][1] && rune((*text)[i+2])==(*words)[k][2] {
                cursor=i
                break
            }
        }
        distances[k]=make( []int, (*words_f)[k]-1, (*words_f)[k]-1)
        for i:=cursor+1; i<len((*text))-2; i++ {
            if rune((*text)[i])==(*words)[k][0] && rune((*text)[i+1])==(*words)[k][1] && rune((*text)[i+2])==(*words)[k][2] {
                distances[k][index]=i-cursor
                index++
                cursor=i
            }
        }
        for i:=0; i<(*words_f)[k]-1; i++ {
            for j:=lower; j<upper; j++ {
                if distances[k][i]%j==0 {
                    (*keylength)[j]++
                }
            }
        }
    }
}

func GuessKeylength( k *[]int) (int, int) {
    var n_max_value1, n_max_value2 int
    var n_max_index1, n_max_index2 int
    n_max_value1=0
    n_max_value2=0
    n_max_index1=0
    n_max_index2=0
    for i:=0; i<len(*k); i++ {
        if (*k)[i]>n_max_value2 {
            if (*k)[i]>n_max_value1 {
                n_max_value2=n_max_value1
                n_max_value1=(*k)[i]
                n_max_index2=n_max_index1
                n_max_index1=i
            } else {
                n_max_value2=(*k)[i]
                n_max_index2=i
            }
        }
    }
    if n_max_index1>=n_max_index2 {
        return n_max_index1, n_max_index2
    } else {
        if n_max_index2 % n_max_index1==0 {
            return n_max_index2, n_max_index1
        } else {
            return n_max_index1, n_max_index2
        }
    }
}

func ReadKeylength( argument *string) int {
    keylength, err:=strconv.Atoi(*argument)
        if err != nil {
            panic(err)
        }
    return keylength
}

func BuildSeparated(s_text *[][]byte, text *[]byte, keylength int, ){
    var s_textlength, num_longer int
    s_textlength=len(*text)/keylength
    num_longer=len(*text)%keylength
    *s_text=make( [][]byte, keylength, keylength)
    for i:=0; i<num_longer; i++ {
        (*s_text)[i]=make( []byte, s_textlength+1, s_textlength+1)
    }
    for i:=num_longer; i<keylength; i++ {
        (*s_text)[i]=make( []byte, s_textlength, s_textlength)
    }
    for i:=0; i<s_textlength; i++ {
        for j:=0; j<keylength; j++ {
            (*s_text)[j][i]=(*text)[j+keylength*i]
        }
    }
    for i:=0; i<num_longer; i++ {
        (*s_text)[i][s_textlength]=(*text)[i+keylength*s_textlength]
    }
}

func Sort3 (i1 int, i2 int, i3 int) (int, int, int ){
    if i1<i2 {
        if i3<=i1 {
            return i3, i1, i2
        } else if i2<=i3{
            return i1, i2, i3
        } else {
            return i1, i3, i2
        }
    } else {// i2<=i1
        if i3<=i2 {
            return i3, i2, i1
        } else if i1<=i3 {
            return i2, i1, i3
        } else {
            return i2, i3, i1
        }
    }
}

func GetShift(r[3] rune, s[3] rune) int {//r is unknow, while s is the result from Table()
    var sorted_r[3] int
    var sorted_s[3] int
    sorted_r[0], sorted_r[1], sorted_r[2]=Sort3(int(r[0]), int(r[1]), int(r[2]))
    sorted_s[0], sorted_s[1], sorted_s[2]=Sort3(int(s[0]), int(s[1]), int(s[2]))
    if sorted_s[1]-sorted_s[0]==sorted_r[1]-sorted_r[0] {
        return (sorted_r[0]-sorted_s[0])
    } else if sorted_s[1]-sorted_s[0]==sorted_r[2]-sorted_r[1] {
        return (sorted_r[1]-sorted_s[0])
    } else if sorted_s[1]-sorted_s[0]==26-sorted_r[2]+sorted_r[0] {
        return (sorted_r[2]-sorted_s[0])
    }
    return 27
}

func Table(r[3] rune) [3]rune {
    var sorted[3] int
    var output[3] rune
    var temp[3] int
    sorted[0], sorted[1], sorted[2]=Sort3(int(r[0]), int(r[1]), int(r[2]))
    temp[0], temp[1],temp[2]=Sort3(sorted[2]-sorted[1], sorted[1]-sorted[0], 26-(sorted[2]-sorted[0]) )
    if (temp[0]==4 && temp[1]==4 && temp[2]==18 ){
        output[0]=65
        output[1]=69
        output[2]=73
    } else if (temp[0]==4 && temp[1]==10 && temp[2]==12 ){
        output[0]=65
        output[1]=69
        output[2]=79
    } else if (temp[0]==4 && temp[1]==7 && temp[2]==15 ){
        output[0]=65
        output[1]=69
        output[2]=84
    } else if (temp[0]==6 && temp[1]==8 && temp[2]==12 ){
        output[0]=65
        output[1]=73
        output[2]=79
    } else if (temp[0]==7 && temp[1]==8 && temp[2]==11 ){
        output[0]=65
        output[1]=73
        output[2]=84
    } else if (temp[0]==5 && temp[1]==7 && temp[2]==14 ){
        output[0]=65
        output[1]=79
        output[2]=84
    } else if (temp[0]==4 && temp[1]==6 && temp[2]==16 ){
        output[0]=69
        output[1]=73
        output[2]=79
    } else if (temp[0]==4 && temp[1]==11 && temp[2]==11 ){
        output[0]=69
        output[1]=73
        output[2]=84
    } else if (temp[0]==5 && temp[1]==10 && temp[2]==11 ){
        output[0]=69
        output[1]=79
        output[2]=84
    } else if (temp[0]==5 && temp[1]==6 && temp[2]==15 ){
        output[0]=73
        output[1]=79
        output[2]=84
    }
    return output
}



func RecoverKey( separatedtext *[][]byte, key *[]byte){
    keylength:=len(*separatedtext)
    var r[3] rune
    (*key)=make([]byte, keylength, keylength)
    buf:=make([]byte, 1)
    // Find possible ET and determine
    for i:=0; i<keylength; i++ {
        r[0], r[1], r[2]=FindMax3( &((*separatedtext)[i]))
        //fmt.Println("r[0], r[1], r[2] are: ", r[0], r[1], r[2], i)
        shift:=GetShift(r, Table(r))
        //fmt.Println(Table(r))
        //fmt.Println(shift)
        if shift<27 {
            if shift<0 {
                shift=shift+26
            }
            utf8.EncodeRune(buf, rune(shift+65))
            (*key)[i]=buf[0]
        } else {
            //fmt.Println("Not sure the character at ", i)
            if r[0]-15==r[1] || r[0]-15==r[1]-26 {
                fmt.Println("the amount of T more than E", i)
                if r[1]-69<0{
                    r[1]=r[1]+26
                }
                utf8.EncodeRune(buf, r[1]-69+65)
                (*key)[i]=buf[0]
            } else {
                if r[0]-69<0{
                    r[0]=r[0]+26
                }
                utf8.EncodeRune(buf, r[0]-69+65)
                (*key)[i]=buf[0]
            }
        }
    }

}


func IndexOfCoincidence (text *[]byte ) float32 {
    var letterfrequency[26] int
    var length int
    var ic float32
    ic=0
    length=len(*text)
    for i:=0; i<length; i++ {
        letterfrequency[rune((*text)[i]-65)]++
    }
    for i:=0; i<26; i++ {
        ic=ic+float32(letterfrequency[i]*(letterfrequency[i]-1))
    }
    ic=ic/float32(length*(length-1))
    return ic
}

func Variance( data *[]float32) float32 {
    var variance float32
    length:=len(*data)
    for i:=0; i<length; i++ {
        variance=variance+((*data)[i]-0.066538462)*((*data)[i]-0.066538462)
    }
    return variance/float32(length)
}

func ICVariance(factors *[]int, text *[]byte) int{
//    var f_length int
    var variance []float32
    f_length:=len(*factors)
    variance=make([]float32, f_length, f_length)
    //separatedtext=make([][][]byte,f_length, f_length) 
    for i:=0; i<f_length; i++ {
        var IC []float32
        IC=make([]float32, (*factors)[i], (*factors)[i] )
        //fmt.Println("factor[i] = ", (*factors)[i])
        var separatedtext [][]byte
        BuildSeparated( &(separatedtext), text, (*factors)[i])
        for j:=0; j<(*factors)[i]; j++ {
            IC[j]=IndexOfCoincidence(&(separatedtext[j]))
        }
        //fmt.Println(IC)
        variance[i]=Variance(&IC)
    }
    //fmt.Println(variance)
    var index int
    var min float32
    index=0
    min=100
    for i:=0; i<f_length; i++ {
        if variance[i]<min {
            min=variance[i]
            index=i
        }
    }
    return (*factors)[index]
}

func FindMax3( text *[]byte) ( rune, rune, rune){
    var n[26] int
    // Counting characters
    for i:=0; i<len(*text); i++ {
        n[rune((*text)[i])-65]++
    }
    // extract the most frequent characters
    var n_max_value1, n_max_value2, n_max_value3 int
    var n_max_index1, n_max_index2, n_max_index3 int
    n_max_value1=0
    n_max_value2=0
    n_max_index1=0
    n_max_index2=0
    n_max_value3=0
    n_max_index3=0
    for i:=0; i<26; i++ {
        if n[i]>n_max_value3 {
            if n[i]>n_max_value2 {
                if n[i]>n_max_value1 {
                    n_max_value3=n_max_value2
                    n_max_value2=n_max_value1
                    n_max_value1=n[i]
                    n_max_index3=n_max_index2
                    n_max_index2=n_max_index1
                    n_max_index1=i
                } else {
                    n_max_value3=n_max_value2
                    n_max_value2=n[i]
                    n_max_index3=n_max_index2
                    n_max_index2=i
                }
            } else {
                n_max_value3=n[i]
                n_max_index3=i
            }
        }
    }
    return rune(n_max_index1+65), rune(n_max_index2+65), rune(n_max_index3+65)
}
func CheckRepeated(key *[]byte){
    var keylength, ceiling, index, check, minvalid int
    var factors []int
    keylength=len(*key)
    ceiling=keylength/2+1
    index=0
    for i:=1; i<ceiling; i++ {
        if keylength%i==0 {
            index=index+1
        }
    }
    factors=make([]int, index, index)
    index=0
    for i:=1; i<ceiling; i++ {
        if keylength%i==0 {
            (factors)[index]=i
            index=index+1
        }
    }
    for i:=index-1; i>0; i-- {
        check=0
        for j:=0; j<keylength-factors[i]; j++ {
            if (*key)[j]!=(*key)[j+factors[i]]{
                check=-1
                break
            }
        }
        if check!=-1 {
            minvalid=factors[i]
        }
    }
    if minvalid!=0 {
        (*key)=(*key)[:minvalid]
    }
    //fmt.Println(minvalid)
}
