package dh

import (
	"crypto/rand"
	"math/big"
)

func isProbablePrimeFermat(n *big.Int, k int) bool {
	if n.Cmp(big.NewInt(2)) < 0 {
		return false
	}
	if n.Cmp(big.NewInt(2)) == 0 {
		return true
	}
	for i := 0; i < k; i++ {
		a, err := rand.Int(rand.Reader, new(big.Int).Sub(n, big.NewInt(2)))
		if err != nil {
			return false
		}
		a = new(big.Int).Add(a, big.NewInt(2))
		if new(big.Int).Exp(a, new(big.Int).Sub(n, big.NewInt(1)), n).Cmp(big.NewInt(1)) != 0 {
			return false
		}
	}
	return true
}

// Функция для теста Миллера-Рабина
func isProbablePrimeMillerRabin(n *big.Int, k int) bool {
	if n.Cmp(big.NewInt(2)) < 0 {
		return false
	}
	if n.Cmp(big.NewInt(2)) == 0 {
		return true
	}
	if n.Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	// Разложение n-1 на 2^s * d
	s := 0
	d := new(big.Int).Sub(n, big.NewInt(1))
	for d.Mod(d, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		d.Rsh(d, 1)
		s++
	}

	for i := 0; i < k; i++ {
		a, err := rand.Int(rand.Reader, new(big.Int).Sub(n, big.NewInt(2)))
		if err != nil {
			return false
		}
		a = new(big.Int).Add(a, big.NewInt(2))
		x := new(big.Int).Exp(a, d, n)
		if x.Cmp(big.NewInt(1)) != 0 && x.Cmp(new(big.Int).Sub(n, big.NewInt(1))) != 0 {
			for r := 1; r < s; r++ {
				x = new(big.Int).Exp(x, big.NewInt(2), n)
				if x.Cmp(big.NewInt(1)) == 0 {
					return false
				}
				if x.Cmp(new(big.Int).Sub(n, big.NewInt(1))) == 0 {
					break
				}
			}
			if x.Cmp(new(big.Int).Sub(n, big.NewInt(1))) != 0 {
				return false
			}
		}
	}
	return true
}

// Функция для генерации простого числа
func generatePrime(bits int) (*big.Int, error) {
	for {
		// Генерация случайного числа с заданным количеством бит
		primeCandidate, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), uint(bits)))
		if err != nil {
			return nil, err
		}
		// Убедимся, что число нечетное
		primeCandidate.SetBit(primeCandidate, 0, 1)

		// Проверка на простоту
		if isProbablePrimeMillerRabin(primeCandidate, 20) && isProbablePrimeFermat(primeCandidate, 10) {
			return primeCandidate, nil
		}

	}
}

// Функция для генерации ключей Диффи-Хеллмана
func GenerateKeys(p, g *big.Int) (*big.Int, *big.Int, error) {
	// Генерация приватного ключа
	privateKey, err := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(2)))
	if err != nil {
		return nil, nil, err
	}
	privateKey.Add(privateKey, big.NewInt(1)) // Убедимся, что ключ больше 0

	// Вычисление публичного ключа
	publicKey := new(big.Int).Exp(g, privateKey, p)

	return privateKey, publicKey, nil
}

func ComputeSharedSecret(publicKeyA, privateKeyB, p *big.Int) *big.Int {
	sharedSecret := new(big.Int).Exp(publicKeyA, privateKeyB, p)
	return sharedSecret
}
