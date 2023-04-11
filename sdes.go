package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	var decifra bool

	if len(os.Args) != 3 {
		if len(os.Args) == 4 {
			if os.Args[3] != "d" {
                    fmt.Println("Comando inválido, para decifrar utilize 'd'")
				os.Exit(1)
			} else {
				decifra = true
			}
		} else {
			fmt.Println("Número de argumentos inválido")
			os.Exit(1)
		}
	}
	entrada, chave := formataArgumentos(os.Args[1], os.Args[2])

	fmt.Println(DES(entrada, chave, decifra))
}

//func cbcCifra(IV, entrada, chave []uint8) (saida []uint8) {

//aux := xor(entrada, IV)
//DES()
//return saida
//}

func formataArgumentos(entrada, chave string) (entradaBytes, chaveBytes []uint8) {
	var aux int
	for i := 1; i <= len(entrada); i++ {
		aux, _ = strconv.Atoi(entrada[i-1 : i])
		entradaBytes = append(entradaBytes, uint8(aux))
	}

	for i := 1; i <= len(chave); i++ {
		aux, _ = strconv.Atoi(chave[i-1 : i])
		chaveBytes = append(chaveBytes, uint8(aux))
	}

	return entradaBytes, chaveBytes
}

func permutacao(lista []uint8, posicoes []uint8) (listaPermutada []uint8) {
	// O for é feito com o len de posições, não de lista, por causa dos P8. No caso dos P8, entra uma lista de 10 posições e sai com 8, se fosse feito com len(lista), iria dar index out of range.
	for i := 0; i < len(posicoes); i++ {
		listaPermutada = append(listaPermutada, lista[posicoes[i]-1])
	}
	return listaPermutada
}

func geraChaves(chave []uint8) (k1, k2 []uint8) {
	P10 := []uint8{3, 5, 2, 7, 4, 10, 1, 9, 8, 6}
	LS1 := []uint8{2, 3, 4, 5, 1}
	LS2 := []uint8{3, 4, 5, 1, 2}
	P8 := []uint8{6, 3, 7, 4, 8, 5, 10, 9}

	chave = permutacao(chave, P10)
	meiaChave1, meiaChave2 := chave[:5], chave[5:10]
	meiaChave1 = permutacao(meiaChave1, LS1)
	meiaChave2 = permutacao(meiaChave2, LS1)
	chave = append(meiaChave1, meiaChave2...)
	k1 = permutacao(chave, P8)
	meiaChave1 = permutacao(meiaChave1, LS2)
	meiaChave2 = permutacao(meiaChave2, LS2)
	chave = append(meiaChave1, meiaChave2...)
	k2 = permutacao(chave, P8)

	return k1, k2
}

func sw(ladoEsquerdo, ladoDireito []uint8) ([]uint8, []uint8) {
	return ladoDireito, ladoEsquerdo
}

func xor(entrada, chave []uint8) (resultadoXor []uint8) {
	for i := 0; i < len(entrada); i++ {
		resultadoXor = append(resultadoXor, entrada[i]^chave[i])
	}
	return resultadoXor
}

func binParaInt(valorBin []uint8) (valorInt uint8) {
	for i := 0; i < len(valorBin); i++ {
		valorInt += valorBin[len(valorBin)-i-1] * uint8(math.Pow(float64(2), float64(i)))
	}
	return valorInt
}

func sBox(listaBits []uint8, matrizS [][]uint8) (valorSaida []uint8) {
	linha := binParaInt([]uint8{listaBits[0], listaBits[3]})
	coluna := binParaInt([]uint8{listaBits[1], listaBits[2]})

	var valorTemp uint8
	valorTemp = matrizS[linha][coluna]
	stringAux := fmt.Sprintf("%b", valorTemp)
	p, _ := strconv.ParseInt(stringAux, 10, 8)
	valorSaida = append(valorSaida, uint8(p)/10)
	valorSaida = append(valorSaida, uint8(p)%10)

	return valorSaida
}

func funcaoF(entradaEsquerda, entradaDireita, chave []uint8) (saidaF []uint8) {
	E_P := []uint8{4, 1, 2, 3, 2, 3, 4, 1}
	P4 := []uint8{2, 4, 3, 1}
	S0 := [][]uint8{
		{1, 0, 3, 2},
		{3, 2, 1, 0},
		{0, 2, 1, 3},
		{3, 1, 3, 2},
	}
	S1 := [][]uint8{
		{0, 1, 2, 3},
		{2, 0, 1, 3},
		{3, 0, 1, 0},
		{2, 1, 0, 3},
	}

	valorTemporario := permutacao(entradaDireita, E_P)
	valorTemporario = xor(valorTemporario, chave)
	ladoSBox0, ladoSBox1 := valorTemporario[0:4], valorTemporario[4:8]
	ladoSBox0 = sBox(ladoSBox0, S0)
	ladoSBox1 = sBox(ladoSBox1, S1)
	saidaParcial := append(ladoSBox0, ladoSBox1...)
	saidaParcial = permutacao(saidaParcial, P4)
	saidaParcial = xor(entradaEsquerda, saidaParcial)

	saidaF = append(saidaParcial, entradaDireita...)

	return saidaF
}

func DES(textoDeEntrada, chave []uint8, decifra bool) (textoDeSaida []uint8) {
	var IP = []uint8{2, 6, 3, 1, 4, 8, 5, 7}
	var IP_1 = []uint8{4, 1, 3, 5, 7, 2, 8, 6}
	var chave_k1, chave_k2 []uint8

	if decifra {
		chave_k2, chave_k1 = geraChaves(chave)
	} else {
		chave_k1, chave_k2 = geraChaves(chave)
	}
	textoDeSaida = permutacao(textoDeEntrada, IP)
	ladoEsquerdo, ladoDireito := textoDeSaida[:4], textoDeSaida[4:8]
	saidaF := funcaoF(ladoEsquerdo, ladoDireito, chave_k1)
	ladoEsquerdo, ladoDireito = saidaF[:4], saidaF[4:8]
	ladoEsquerdo, ladoDireito = sw(ladoEsquerdo, ladoDireito)
	textoDeSaida = funcaoF(ladoEsquerdo, ladoDireito, chave_k2)
	textoDeSaida = permutacao(textoDeSaida, IP_1)

	return textoDeSaida
}
